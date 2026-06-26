package viewer

import (
	"fmt"
	"os"
	"time"

	"github.com/longxiucai/go-tuimg/termeverything"
	"golang.org/x/term"
)

func (v *ImageViewer) HandleInput() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	defer func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		v.ExitFullScreen()
	}()

	fmt.Print("\033[?1049h")
	fmt.Print("\033[?25l")
	fmt.Print("\033[?1003h")
	fmt.Print("\033[?1006h")
	fmt.Print("\033[2J")
	fmt.Print("\033[H")

	ts := termeverything.MakeTermSize()
	v.LastWidth = ts.WidthCells
	v.LastHeight = ts.HeightCells

	v.Render()

	defer func() {
		fmt.Print("\033[?1006l")
		fmt.Print("\033[?1003l")
		fmt.Print("\033[?25h")
		fmt.Print("\033[?1049l")
		drainStdin()
	}()

	var isDragging bool
	var lastMouseX, lastMouseY int
	inputChan := make(chan []byte, 1)

	go func() {
		buf := make([]byte, 4096)
		for v.Running {
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				continue
			}
			data := make([]byte, n)
			copy(data, buf[:n])
			inputChan <- data
		}
	}()

	sigWinch := watchWinch()

	moveStep := 3
	for v.Running {
		select {
		case chunk := <-inputChan:
			v.handleInputChunk(chunk, &isDragging, &lastMouseX, &lastMouseY, moveStep)
		case <-v.NeedRedraw:
			v.Render()
		case <-sigWinch:
			ts := termeverything.MakeTermSize()
			if ts.WidthCells != v.LastWidth || ts.HeightCells != v.LastHeight {
				v.OffsetX = 0
				v.OffsetY = 0
				v.LastWidth = ts.WidthCells
				v.LastHeight = ts.HeightCells
			}
			v.Render()
		case <-time.After(50 * time.Millisecond):
		}
	}
}

func (v *ImageViewer) handleInputChunk(chunk []byte, isDragging *bool, lastMouseX, lastMouseY *int, moveStep int) {
	if len(chunk) >= 4 && chunk[0] == '\x1b' && chunk[1] == '[' && chunk[2] == '<' {
		events := termeverything.ParseSGRMouseSequences(chunk)
		for _, ev := range events {
			switch e := ev.(type) {
			case *termeverything.PointerWheel:
				if e.Up {
					v.ScaleFactor = minVal(v.MaxScale, v.ScaleFactor+v.ScaleStep)
				} else {
					v.ScaleFactor = maxVal(v.MinScale, v.ScaleFactor-v.ScaleStep)
				}
				v.Render()
			case *termeverything.PointerMove:
				if *isDragging {
					v.OffsetX += e.Col - *lastMouseX
					v.OffsetY += e.Row - *lastMouseY
					v.Render()
				}
				*lastMouseX = e.Col
				*lastMouseY = e.Row
			case *termeverything.PointerButtonPress:
				if e.Button == termeverything.BTN_LEFT {
					*isDragging = true
				}
			case *termeverything.PointerButtonRelease:
				if e.Button == termeverything.BTN_LEFT {
					*isDragging = false
				}
			}
		}
		return
	}

	if len(chunk) >= 6 && chunk[0] == '\x1b' && chunk[1] == '[' && chunk[2] == 'M' {
		events := termeverything.ParseSGRMouseSequences(chunk)
		for _, ev := range events {
			switch e := ev.(type) {
			case *termeverything.PointerWheel:
				if e.Up {
					v.ScaleFactor = minVal(v.MaxScale, v.ScaleFactor+v.ScaleStep)
				} else {
					v.ScaleFactor = maxVal(v.MinScale, v.ScaleFactor-v.ScaleStep)
				}
				v.Render()
			case *termeverything.PointerMove:
				if *isDragging {
					v.OffsetX += e.Col - *lastMouseX
					v.OffsetY += e.Row - *lastMouseY
					v.Render()
				}
				*lastMouseX = e.Col
				*lastMouseY = e.Row
			case *termeverything.PointerButtonPress:
				if e.Button == termeverything.BTN_LEFT {
					*isDragging = true
				}
			case *termeverything.PointerButtonRelease:
				if e.Button == termeverything.BTN_LEFT {
					*isDragging = false
				}
			}
		}
		return
	}

	codes := termeverything.ConvertKeycodeToXbdCode(chunk)
	for _, code := range codes {
		if c, ok := code.(*termeverything.KeyCode); ok {
			switch c.KeyCode {
			case termeverything.KEY_Q:
				v.Running = false
				return
			case termeverything.KEY_EQUAL, termeverything.KEY_KPPLUS:
				v.ScaleFactor = minVal(v.MaxScale, v.ScaleFactor+v.ScaleStep)
				v.Render()
			case termeverything.KEY_MINUS, termeverything.KEY_KPMINUS:
				v.ScaleFactor = maxVal(v.MinScale, v.ScaleFactor-v.ScaleStep)
				v.Render()
			case termeverything.KEY_0, termeverything.KEY_KP0:
				v.ScaleFactor = 1.0
				v.OffsetX = 0
				v.OffsetY = 0
				v.Render()
			case termeverything.KEY_H, termeverything.KEY_LEFT:
				v.OffsetX += moveStep
				v.Render()
			case termeverything.KEY_L, termeverything.KEY_RIGHT:
				v.OffsetX -= moveStep
				v.Render()
			case termeverything.KEY_J, termeverything.KEY_DOWN:
				v.OffsetY -= moveStep
				v.Render()
			case termeverything.KEY_K, termeverything.KEY_UP:
				v.OffsetY += moveStep
				v.Render()
			}
		}
	}
}
