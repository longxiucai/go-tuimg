//go:build !windows

package termeverything

import (
	"syscall"
	"unsafe"
)

func getWinsize(fd uintptr) (WinSize, error) {
	var ws WinSize
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&ws)))
	if errno != 0 {
		return ws, errno
	}
	return ws, nil
}

func GetWinsize(fd uintptr) (WinSize, error) {
	return getWinsize(fd)
}
