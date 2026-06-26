//go:build windows

package viewer

import "os"

func watchWinch() <-chan os.Signal {
	return make(chan os.Signal, 1)
}
