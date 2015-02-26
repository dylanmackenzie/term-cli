package terminal

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"dylanmackenzie.com/term-cli/logger"
	"dylanmackenzie.com/term-cli/recorder"

	"github.com/kr/pty"
	tty "golang.org/x/crypto/ssh/terminal"
)

// Start initiates a new terminal session and returns the file
// descriptor of the tty.
func Start(cmd string) (*os.File, error) {
	c := exec.Command(cmd)

	return pty.Start(c)
}

// Record starts a new terminal session using the current stdin and
// stdout and sends all of its output to dupe.
func Record(dupe rec.Recorder) {
	done := make(chan struct{})
	sh := os.Getenv("SHELL")
	if sh == "" {
		sh = "/bin/bash"
	}

	// Put stdin in raw mode so each keystroke is sent to our new tty
	stdin := int(os.Stdin.Fd())
	oldState, err := tty.MakeRaw(stdin)
	if err != nil {
		logger.Fatal("Could not put stdin in raw mode\n")
	}
	defer tty.Restore(stdin, oldState)

	session, err := Start(sh)
	if err != nil {
		logger.Fatal("Failed to start new tty\n")
	}
	defer session.Close()

	resize := func() {
		cols, rows, err := tty.GetSize(stdin)
		if err != nil {
			logger.Fatal("Could not get terminal size: %s", err)
		}
		Resize(session, cols, rows)
		dupe.Resize(cols, rows)
	}

	resize()

	// Resize session when window is resized
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGWINCH)
	go func() {
		for _ = range signals {
			resize()
		}
	}()
	defer signal.Stop(signals)

	// Copy stdin to our terminal session
	go func() {
		io.Copy(session, os.Stdin)
	}()

	// Copy stdout from our terminal session
	go func() {
		mw := io.MultiWriter(dupe, os.Stdout)
		io.Copy(mw, session)
		done <- struct{}{}
	}()

	// Wait for new tty to exit
	<-done
}
