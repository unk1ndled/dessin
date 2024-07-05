package dessin

import (
	"github.com/unk1ndled/draw/mouse"
	"github.com/unk1ndled/draw/primitives"
	"github.com/unk1ndled/draw/util"
	"github.com/unk1ndled/draw/window"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	Mouse  *mouse.MouseState
	pixels []byte

	colors = []window.Color{
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
	pt.canvas = NewCanvas(int32(outline), int32(outline), window.ScreenWidth-int32(outline), window.ScreenHeight-int32(outline))

}

func (pt *Paint) Update() (bool, bool) {
	Mouse.Update()
	return pt.canvas.Update(), checkKeyPress()
}
func (pt *Paint) Render() {
	pt.canvas.Draw()
}

func checkKeyPress() bool {
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
	PEN = iota
)

type Canvas struct {

	// pos[0] = x and pos[1] = y
	pos    [2]int32
	width  int32
	height int32

	buffer *util.Deque[[]byte]

	currentTool Tool
}

func NewCanvas(x, y, w, h int32) *Canvas {
	return &Canvas{pos: [2]int32{x, y}, width: w, height: h, buffer: util.NewDeque[[]byte](), currentTool: PEN}
}

func (cvs *Canvas) isHovered() bool {
	if Mouse.X >= int(cvs.pos[0]) && Mouse.Y > int(cvs.pos[1]) && Mouse.X < int(cvs.width) && Mouse.Y < int(cvs.height) {
		return true
	}
	return false
}

func (cvs *Canvas) isClicked() bool {
	return cvs.isHovered() && Mouse.LeftButton
}

func (cvs *Canvas) Update() bool {

	if CTRL && !Z_PRESS && Z_PREV_PRESS {
		Z_PREV_PRESS = false
		cvs.undo()
		return true
	} else if cvs.isHovered() && Mouse.LeftButton && !Mouse.PrevLeftButton {
		cvs.updateBuffer()
	}

	return cvs.isClicked()
}

func (cvs *Canvas) Draw() {

	if cvs.isClicked() {

		switch cvs.currentTool {
		case PEN:
			primitives.DrawLine(Mouse.PrevX, Mouse.PrevY, Mouse.X, Mouse.Y, &colors[1])
		}
	}

}

func (cvs *Canvas) updateBuffer() {
	if cvs.buffer.Size() == 10 {
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
