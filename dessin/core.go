package dessin

import (
	"log"

	"github.com/unk1ndled/draw/mouse"
	"github.com/unk1ndled/draw/primitives"
	"github.com/unk1ndled/draw/util"
	"github.com/unk1ndled/draw/window"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	Mouse  *mouse.MouseState
	pixels []byte

	basecolor = window.Color{R: 0, G: 0, B: 0}
	colors    = []window.Color{
		{R: 250, G: 0, B: 0},
		{R: 0, G: 255, B: 0},
		{R: 0, G: 0, B: 255}}

	Z_PREV_PRESS, Z_PRESS = false, false
	CTRL                  = false
)

// implements sdl.Runnable
// Paint acts as the wrapper struct for the program
// thus its responsible for the inner components
type Paint struct {
	canvas *Canvas
}

func NewPaint() *Paint {
	return &Paint{}
}
func (pt *Paint) Init(pxls []byte) {
	pixels = pxls
	Mouse = mouse.GetMouseState()
	outline := 10
	pt.canvas = NewCanvas(int32(outline), int32(outline), window.ScreenWidth-2*int32(outline), window.ScreenHeight-2*int32(outline))

}

func (pt *Paint) Update() (bool, bool) {
	Mouse.Update()
	return pt.canvas.Update(), pt.checkKeyPress()
}
func (pt *Paint) Render() {
	pt.canvas.Draw()
}

func (pt *Paint) checkKeyPress() bool {
	for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		switch t := e.(type) {
		case *sdl.QuitEvent:
			return true
		case *sdl.KeyboardEvent:
			switch t.Type {
			case sdl.KEYDOWN:
				switch t.Keysym.Scancode {
				case sdl.SCANCODE_LCTRL, sdl.SCANCODE_RCTRL:
					CTRL = true
				case sdl.SCANCODE_Z:
					Z_PREV_PRESS = Z_PRESS
					Z_PRESS = true
				case sdl.SCANCODE_F:
					pt.canvas.setTool(FILL)
				case sdl.SCANCODE_D:
					pt.canvas.setTool(PEN)

				}

			case sdl.KEYUP:
				switch t.Keysym.Scancode {
				case sdl.SCANCODE_LCTRL, sdl.SCANCODE_RCTRL:
					CTRL = false
				case sdl.SCANCODE_Z:
					Z_PREV_PRESS = Z_PRESS
					Z_PRESS = false
				}
			}
		}
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Tool byte

const (
	MAX_BUFFER_LENGTH = 15

	PEN = iota
	FILL
)

type Canvas struct {
	X      int32
	Y      int32
	Width  int32
	Height int32

	buffer *util.Deque[[]byte]

	currentTool Tool
}

func NewCanvas(x, y, w, h int32) *Canvas {
	return &Canvas{X: x, Y: y, Width: w, Height: h, buffer: util.NewDeque[[]byte](), currentTool: PEN}
}

func (cvs *Canvas) doesCover(x, y int) bool {
	if x >= int(cvs.X) && y > int(cvs.Y) && x < int(cvs.X+cvs.Width) && y < int(cvs.Y+cvs.Height) {
		return true
	}
	return false
}

func (cvs *Canvas) isHovered() bool {
	return cvs.doesCover(Mouse.X, Mouse.Y)
}

func (cvs *Canvas) isClicked() bool {
	return cvs.isHovered() && Mouse.LeftButton
}

func (cvs *Canvas) wasClicked() bool {
	res := Mouse.PrevLeftButton && cvs.doesCover(Mouse.PrevX, Mouse.PrevY)
	if res {
		log.Println("was clicked")
	}
	return res
}

func (cvs *Canvas) Update() bool {

	if CTRL && !Z_PRESS && Z_PREV_PRESS {
		Z_PREV_PRESS = false
		cvs.undo()
		return true
	} else if cvs.isClicked() && !Mouse.PrevLeftButton {
		cvs.updateBuffer()
	}

	return cvs.isClicked() || cvs.wasClicked()
}

func (cvs *Canvas) Draw() {

	if cvs.isClicked() {

		switch cvs.currentTool {
		case PEN:
			x1, y1, x2, y2 := cvs.ValidateX(Mouse.PrevX), cvs.ValidateY(Mouse.PrevY), cvs.ValidateX(Mouse.X), cvs.ValidateY(Mouse.Y)
			primitives.DrawLine(x1, y1, x2, y2, &colors[1])
		}
	} else if cvs.wasClicked() && !Mouse.LeftButton {
		log.Println("will fill")

		switch cvs.currentTool {
		case FILL:
			clicked := window.GetPixelColor(Mouse.X, Mouse.Y)
			cvs.Fill(Mouse.X, Mouse.Y, &colors[2], &clicked)
		}
	}

}

func (cvs *Canvas) updateBuffer() {
	if cvs.buffer.Size() == MAX_BUFFER_LENGTH {
		cvs.buffer.PopFront()
	}

	// Save a copy of the current pixel state
	cvs.buffer.PushBack(append([]byte(nil), pixels...))
}

func (cvs *Canvas) undo() {
	if cvs.buffer.Size() > 0 {
		lastState, _ := cvs.buffer.PopBack()
		copy(pixels, lastState) // Restore pixel data
	}
}

func (cvs *Canvas) setTool(t Tool) {
	cvs.currentTool = t
}

func (cvs *Canvas) ValidateX(x int) int {
	if x < int(cvs.X) {
		return int(cvs.X)
	} else if x > int(cvs.X)+int(cvs.Width) {
		return int(cvs.X + cvs.Width)
	} else {
		return x
	}
}
func (cvs *Canvas) ValidateY(y int) int {
	if y < int(cvs.Y) {
		return int(cvs.Y)
	} else if y > int(cvs.Y)+int(cvs.Height) {
		return int(cvs.Y + cvs.Height)
	} else {
		return y
	}
}

func (cvs *Canvas) Fill(x0, y0 int, fillColor, clickedColor *window.Color) {
	if x0 > int(cvs.X) && x0 < int(cvs.X)+int(cvs.Width) && y0 > int(cvs.Y) && y0 < int(cvs.Y+cvs.Height) {
		curPixelCLr := window.GetPixelColor(x0, y0)
		if curPixelCLr.R == clickedColor.R && curPixelCLr.G == clickedColor.G && curPixelCLr.B == clickedColor.B {
			window.SetPixel(x0, y0, fillColor)
			cvs.Fill(x0+1, y0, fillColor, clickedColor)
			cvs.Fill(x0-1, y0, fillColor, clickedColor)
			cvs.Fill(x0, y0-1, fillColor, clickedColor)
			cvs.Fill(x0, y0+1, fillColor, clickedColor)
		}
	}
}
