package rec

import (
	"encoding/csv"
	"io"
	"os"
	"path"
	"strconv"
	"time"
)

// Play replays a terminal session on stdout
func Play(out *os.File, dir string) error {
	data, err := os.Open(path.Join(dir, "term.data"))
	if err != nil {
		return err
	}
	timing, err := os.Open(path.Join(dir, "term.timing"))
	if err != nil {
		return err
	}

	times := csv.NewReader(timing)
	times.FieldsPerRecord = -1 // variable number of fields
	buf := [100 * 1024]byte{}
	start := time.Now()

	for {
		// Read field from timing info
		field, err := times.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		// Resize event
		if field[1] == "r" {
			continue
		}

		t, _ := strconv.Atoi(field[0])
		n, _ := strconv.Atoi(field[1])
		dt := time.Duration(t)
		ch := time.After(dt - time.Since(start))
		sl := buf[:n]
		io.ReadFull(data, sl)
		<-ch
		os.Stdout.Write(sl)
	}

	return nil
}
