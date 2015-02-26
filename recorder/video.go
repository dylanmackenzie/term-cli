package rec

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

// A Recorder stores the output of a terminal session
type Recorder interface {
	io.WriteCloser
	Resize(cols, rows int) error
}

// A Video recorder stores timing information and the raw data from a
// terminal session.  Videos are meant to be recorded in full and then
// played back at a later time.
type Video struct {
	data, timing *os.File
	Path         string
	Start        time.Time
}

// NewVideo returns a new Video instance
func NewVideo(root string) (*Video, error) {
	const layout = "Recording from 2006-01-02 15:04:05" // used to name files
	const flags = os.O_CREATE | os.O_EXCL | os.O_RDWR
	var err error

	v := new(Video)
	v.Start = time.Now()

	// Make recording folder
	name := v.Start.Format(layout)
	dir := path.Join(root, name)
	v.Path = dir
	if err := os.Mkdir(dir, 0775); err != nil {
		return nil, err
	}

	// Open data file
	dpath := path.Join(dir, "term.data")
	if v.data, err = os.OpenFile(dpath, flags, 0644); err != nil {
		return nil, err
	}

	// Open timing file
	tpath := path.Join(dir, "term.timing")
	if v.timing, err = os.OpenFile(tpath, flags, 0644); err != nil {
		return nil, err
	}

	return v, nil
}

// Resize records a resize event in the timing file
func (v *Video) Resize(cols, rows int) error {
	_, err := fmt.Fprintf(v.timing, "%d,r,%d,%d\n", time.Since(v.Start), cols, rows)
	return err
}

// Write stores the terminal data and records how many bytes were
// written in the timing file
func (v *Video) Write(p []byte) (int, error) {
	// Record timing data
	_, err := fmt.Fprintf(v.timing, "%d,%d\n", time.Since(v.Start), len(p))
	if err != nil {
		return 0, err
	}

	return v.data.Write(p)
}

func (v *Video) Close() error {
	v.data.Close()
	v.timing.Close()

	return nil
}
