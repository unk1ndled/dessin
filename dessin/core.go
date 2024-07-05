package dessin

import (
	"github.com/unk1ndled/draw/mouse"
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
	xoffset, yoffest := 10, 10
	topbar := 40
	pt.canvas = NewCanvas(int32(xoffset),
		int32(yoffest+topbar),
		window.ScreenWidth-2*int32(xoffset),
		window.ScreenHeight-int32(2*(yoffest)+(topbar)))
	pt.setBackground()

}

func (pt *Paint) setBackground() {
	bg := &window.Color{R: 15, G: 15, B: 15}
	for i := 0; i < int(window.ScreenWidth); i++ {
		for j := 0; j < int(window.ScreenHeight); j++ {
			if !pt.canvas.contains(int32(i), int32(j)) {
				window.SetPixel(int32(i), int32(j), bg)
			}
		}
	}
}

func (pt *Paint) Update() (bool, bool) {
	Mouse.Update()
	quit := pt.checkKeyPress()
	update := pt.canvas.Update()
	return update, quit
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
				case sdl.SCANCODE_E:
					pt.canvas.setTool(ERASER)
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
