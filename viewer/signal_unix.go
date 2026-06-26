//go:build !windows

package viewer

import (
	"os"
	"os/signal"
	"syscall"
)

func watchWinch() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	return ch
}
