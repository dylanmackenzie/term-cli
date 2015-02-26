package cli

import (
	"flag"
	"fmt"
	"os"

	"dylanmackenzie.com/term-cli/logger"
	"dylanmackenzie.com/term-cli/recorder"
	"dylanmackenzie.com/term-cli/terminal"
)

const usage = `
Usage: term-cast COMMAND [OPTIONS]

A recorder for terminal sessions

Commands:
	play 		Replay a saved terminal session
	record 		Record a new terminal session
	upload 		Upload a saved terminal session

Run 'term-cast COMMMAND --help' for more information on a command
`

func Run() {
	sub := ""

	if len(os.Args) > 1 {
		sub = os.Args[1]
	} else {
		fmt.Print(usage)
		os.Exit(1)
	}

	// Remove subcommand from list of arguments
	os.Args = append(os.Args[:1], os.Args[2:]...)

	switch sub {
	case "record":
		record()
	case "play":
		play()
	case "upload":
		upload()
	case "stream":
		stream()
	default:
		fmt.Print(usage)
		os.Exit(1)
	}
}

func record() {
	f := flag.NewFlagSet("record", flag.ExitOnError)

	f.Parse(os.Args)
	dir := f.Arg(1)
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		logger.Fatal("%s is not a directory", dir)
	}

	v, err := rec.NewVideo(dir)
	if err != nil {
		panic(err)
	}
	defer v.Close()

	logger.Log("Recording session in %s\n", v.Path)
	terminal.Record(v)
	logger.Log("Recording finished")
}

func play() {
	f := flag.NewFlagSet("play", flag.ExitOnError)
	//speed := f.Float64("speed", 1, "Sets the playback speed")

	f.Parse(os.Args)

	dir := f.Arg(1)
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		logger.Fatal("%s is not a directory", dir)
	}

	logger.Log("Replaying session\n")
	err := rec.Play(os.Stdout, dir)
	if err != nil {
		logger.Fatal("%s", err)
	}
	logger.Log("Replay finished\n")
}

func upload() {
	// f := flag.NewFlagSet("upload", flag.ExitOnError)
}

func stream() {

}
