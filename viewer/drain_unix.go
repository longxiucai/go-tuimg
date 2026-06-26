//go:build !windows

package viewer

import (
	"os"
	"syscall"
)

func drainStdin() {
	fd := int(os.Stdin.Fd())
	syscall.SetNonblock(fd, true)
	drain := make([]byte, 4096)
	for {
		if _, err := syscall.Read(fd, drain); err != nil {
			break
		}
	}
	syscall.SetNonblock(fd, false)
}
