package terminal

// #include <sys/ioctl.h>
import "C"

import (
	"errors"
	"os"
	unix "syscall" // will be changed to sys/unix on 1.4
	"unsafe"
)

// Resize a given tty to have a size of rows x columns using an ioctl
// call
func Resize(tty *os.File, cols, rows int) error {
	var ws C.struct_winsize
	ws.ws_row = C.ushort(rows)
	ws.ws_col = C.ushort(cols)

	_, _, e := unix.Syscall(
		unix.SYS_IOCTL,
		tty.Fd(),
		unix.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&ws)),
	)

	if e != 0 {
		return errors.New("Failed to resize terminal")
	}
	return nil
}
