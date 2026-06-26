package viewer

import (
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/longxiucai/go-tuimg/termeverything"
)

type ImageViewer struct {
	ImagePath     string
	OriginalImage image.Image
	ScaleFactor   float64
	MinScale      float64
	MaxScale      float64
	ScaleStep     float64
	OffsetX       int
	OffsetY       int
	LastWidth     int
	LastHeight    int
	Running       bool
	GIF           *gif.GIF
	GIFFrames     []image.Image
	GIFDelays     []int
	GIFIndex      int
	IsGIF         bool
	NeedRedraw    chan bool
}

func New() *ImageViewer {
	return &ImageViewer{
		ScaleFactor: 1.0,
		MinScale:    0.2,
		MaxScale:    5.0,
		ScaleStep:   0.2,
		Running:     true,
		NeedRedraw:  make(chan bool, 1),
	}
}

func (v *ImageViewer) LoadImage(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "无法打开文件: %v\n", err)
		return false
	}
	defer file.Close()

	if gifImg, err := gif.DecodeAll(file); err == nil && len(gifImg.Image) > 1 {
		v.GIF = gifImg
		v.GIFFrames = make([]image.Image, len(gifImg.Image))
		v.GIFDelays = make([]int, len(gifImg.Image))
		for i, frame := range gifImg.Image {
			v.GIFFrames[i] = frame
			delay := gifImg.Delay[i]
			if delay <= 0 {
				delay = 10
			}
			v.GIFDelays[i] = delay * 10
		}
		v.OriginalImage = v.GIFFrames[0]
		v.IsGIF = true
	} else {
		file.Seek(0, 0)
		img, _, err := image.Decode(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "无法解码图片: %v\n", err)
			return false
		}
		v.OriginalImage = img
		v.IsGIF = false
	}

	v.ImagePath = path
	return true
}

func (v *ImageViewer) NextFrame() {
	if v.IsGIF && len(v.GIFFrames) > 1 {
		v.GIFIndex = (v.GIFIndex + 1) % len(v.GIFFrames)
		v.OriginalImage = v.GIFFrames[v.GIFIndex]
	}
}

func (v *ImageViewer) GetFrameDelay() int {
	if v.IsGIF && len(v.GIFDelays) > 0 {
		return v.GIFDelays[v.GIFIndex]
	}
	return 0
}

func (v *ImageViewer) ImageToColorANSI(img image.Image, targetWidth, targetHeight int) string {
	if img == nil {
		return ""
	}

	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	imgAspect := float64(origWidth) / float64(origHeight)
	charAspect := 0.5

	canvasHeight := targetHeight * 2
	newWidth := int(float64(targetWidth) * v.ScaleFactor)
	newHeight := int(float64(canvasHeight) * v.ScaleFactor)

	canvasAspect := float64(newWidth) * charAspect * 2.0 / float64(newHeight)

	if canvasAspect > imgAspect {
		newWidth = int(float64(newHeight) * imgAspect / (charAspect * 2.0))
	} else {
		newHeight = int(float64(newWidth) * charAspect * 2.0 / imgAspect)
	}

	if newWidth < 1 {
		newWidth = 1
	}
	if newHeight < 1 {
		newHeight = 1
	}

	centerOffsetX := (targetWidth - newWidth) / 2
	centerOffsetY := (targetHeight - newHeight/2) / 2

	offsetX := centerOffsetX + v.OffsetX
	offsetY := centerOffsetY + v.OffsetY

	// -----------------------------
	// 预计算 X 映射
	// -----------------------------
	xMap := make([]int, newWidth)
	for x := 0; x < newWidth; x++ {
		xMap[x] = bounds.Min.X + x*origWidth/newWidth
	}

	// -----------------------------
	// 预计算 Y 映射
	// -----------------------------
	yMap := make([]int, newHeight)
	for y := 0; y < newHeight; y++ {
		yMap[y] = bounds.Min.Y + y*origHeight/newHeight
	}

	// Builder，预估容量
	var b strings.Builder
	b.Grow(targetWidth * targetHeight * 45)

	for y := 0; y < targetHeight; y++ {

		topY := (y - offsetY) * 2
		botY := topY + 1

		for x := 0; x < targetWidth; x++ {

			imgX := x - offsetX

			if imgX >= 0 &&
				imgX < newWidth &&
				topY >= 0 &&
				botY < newHeight {

				srcX := xMap[imgX]

				tr, tg, tb, _ := img.At(srcX, yMap[topY]).RGBA()
				br, bg, bb, _ := img.At(srcX, yMap[botY]).RGBA()

				tr >>= 8
				tg >>= 8
				tb >>= 8

				br >>= 8
				bg >>= 8
				bb >>= 8

				if tr == br && tg == bg && tb == bb {

					writeBG(&b,
						uint8(tr),
						uint8(tg),
						uint8(tb),
					)
					b.WriteByte(' ')

				} else {

					writeFG(&b,
						uint8(tr),
						uint8(tg),
						uint8(tb),
					)

					writeBG(&b,
						uint8(br),
						uint8(bg),
						uint8(bb),
					)

					b.WriteRune('▀')
				}

			} else {

				b.WriteString("\033[48;2;0;0;0m ")
			}
		}

		b.WriteString("\033[0m\r\n")
	}

	return b.String()
}

func writeFG(b *strings.Builder, r, g, bl uint8) {
	b.WriteString("\033[38;2;")
	appendUint8(b, r)
	b.WriteByte(';')
	appendUint8(b, g)
	b.WriteByte(';')
	appendUint8(b, bl)
	b.WriteByte('m')
}

func writeBG(b *strings.Builder, r, g, bl uint8) {
	b.WriteString("\033[48;2;")
	appendUint8(b, r)
	b.WriteByte(';')
	appendUint8(b, g)
	b.WriteByte(';')
	appendUint8(b, bl)
	b.WriteByte('m')
}

func appendUint8(b *strings.Builder, v uint8) {
	b.WriteString(strconv.Itoa(int(v)))
}

func (v *ImageViewer) ImageToColorANSIold(img image.Image, targetWidth, targetHeight int) string {
	if img == nil {
		return ""
	}

	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()
	imgAspect := float64(origWidth) / float64(origHeight)
	charAspect := 0.5

	canvasHeight := targetHeight * 2
	newWidth := int(float64(targetWidth) * v.ScaleFactor)
	newHeight := int(float64(canvasHeight) * v.ScaleFactor)
	canvasAspect := float64(newWidth) * charAspect * 2.0 / float64(newHeight)

	if canvasAspect > imgAspect {
		newWidth = int(float64(newHeight) * imgAspect / (charAspect * 2.0))
	} else {
		newHeight = int(float64(newWidth) * charAspect * 2.0 / imgAspect)
	}

	if newWidth < 1 {
		newWidth = 1
	}
	if newHeight < 1 {
		newHeight = 1
	}

	centerOffsetX := (targetWidth - newWidth) / 2
	centerOffsetY := (targetHeight - newHeight/2) / 2
	offsetX := centerOffsetX + v.OffsetX
	offsetY := centerOffsetY + v.OffsetY

	var result []byte
	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			topY := (y - offsetY) * 2
			botY := topY + 1
			imgX := x - offsetX

			if imgX >= 0 && imgX < newWidth && topY >= 0 && botY < newHeight {
				srcTopX := bounds.Min.X + imgX*origWidth/newWidth
				srcTopY := bounds.Min.Y + topY*origHeight/newHeight
				srcBotX := bounds.Min.X + imgX*origWidth/newWidth
				srcBotY := bounds.Min.Y + botY*origHeight/newHeight

				tr, tg, tb, _ := img.At(srcTopX, srcTopY).RGBA()
				br, bg, bb, _ := img.At(srcBotX, srcBotY).RGBA()

				tr >>= 8
				tg >>= 8
				tb >>= 8
				br >>= 8
				bg >>= 8
				bb >>= 8

				if tr == br && tg == bg && tb == bb {
					result = append(result, fmt.Appendf(nil, "\033[48;2;%d;%d;%dm ", tr, tg, tb)...)
				} else {
					result = append(result, fmt.Appendf(nil, "\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm\u2580", tr, tg, tb, br, bg, bb)...)
				}
			} else {
				result = append(result, []byte("\033[48;2;0;0;0m ")...)
			}
		}
		result = append(result, []byte("\033[0m\r\n")...)
	}
	return string(result)
}

func (v *ImageViewer) Render() {
	if v.OriginalImage == nil {
		return
	}
	ts := termeverything.MakeTermSize()
	ansiOutput := v.ImageToColorANSI(v.OriginalImage, ts.WidthCells, ts.HeightCells)
	fmt.Print("\033[H")
	fmt.Print(ansiOutput)
	fmt.Print("\033[0J")
}

func (v *ImageViewer) ExitFullScreen() {
	fmt.Print("\033[?25h")
	fmt.Print("\033[?1049l")
}

func (v *ImageViewer) Run(path string) {
	if !v.LoadImage(path) {
		return
	}
	if v.IsGIF {
		go v.animateGIF()
	}
	v.HandleInput()
}

func (v *ImageViewer) animateGIF() {
	for v.Running {
		delay := v.GetFrameDelay()
		if delay <= 0 {
			delay = 100
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
		if !v.Running {
			return
		}
		v.NextFrame()
		select {
		case v.NeedRedraw <- true:
		default:
		}
	}
}

func minVal(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxVal(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
