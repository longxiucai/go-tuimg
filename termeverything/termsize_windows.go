//go:build windows

package termeverything

import "golang.org/x/term"

func getWinsize(fd uintptr) (WinSize, error) {
	w, h, err := term.GetSize(int(fd))
	if err != nil {
		return WinSize{}, err
	}
	return WinSize{Row: uint16(h), Col: uint16(w)}, nil
}

func GetWinsize(fd uintptr) (WinSize, error) {
	return getWinsize(fd)
}
