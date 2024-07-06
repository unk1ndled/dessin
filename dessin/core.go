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

	Z_PREV_PRESS, Z_PRESS = false, false
	CTRL                  = false
)

// implements sdl.Runnable
// Paint acts as the wrapper struct for the program
// thus its responsible for the inner components
type Paint struct {
	topbar  *Bar
	leftbar *Bar

	canvas *Canvas
}

func NewPaint() *Paint {
	return &Paint{}
}
func (pt *Paint) Init(pxls []byte) {

	Padding := int32(6)

	BarPadding := Padding / 2
	TopBarHeight := int32(45)
	TopBarX := Padding
	TopBarY := Padding
	TopBarWidth := window.ScreenWidth - 2*Padding

	ButtonWidth := TopBarHeight + int32(TopBarHeight/3)
	ButtonHeight := TopBarHeight - 2*BarPadding
	ButtonGap := Padding - 1

	LeftBarX := Padding
	LeftBarY := TopBarHeight + 2*Padding
	LeftBarWidth := ButtonWidth + 2*BarPadding
	LeftBarHeight := window.ScreenHeight - (2*Padding + TopBarHeight)

	pixels = pxls
	Mouse = mouse.GetMouseState()
	pt.setBackground()

	pt.canvas = NewCanvas(
		(2*Padding)+LeftBarWidth,
		(2*Padding)+TopBarHeight,
		window.ScreenWidth-((2*Padding)+LeftBarWidth+10),
		window.ScreenHeight-((2*Padding)+TopBarHeight+10)) // Corrected from LeftBarHeight to TopBarHeight

	fns := []func(){
		func() { pt.canvas.setTool(FILL) },
		func() { pt.canvas.setTool(PEN) },
		func() { pt.canvas.setTool(ERASER) },
		func() { pt.canvas.setColor(0) },
		func() { pt.canvas.setColor(1) },
		func() { pt.canvas.setColor(2) },
		func() { pt.canvas.modifyDrawWidth(-1) },
		func() { pt.canvas.modifyDrawWidth(1) },
	}
	btnc := &window.Color{R: 40, G: 40, B: 40}
	pt.topbar = NewBar(
		TopBarX, TopBarY, TopBarWidth, TopBarHeight,
		ButtonWidth, ButtonHeight, ButtonGap, BarPadding, HORIZONTAL,
		[]*BtnConfig{
			{Color: btnc, Fn: fns[1]},
			{Color: btnc, Fn: fns[2]},
			{Color: btnc, Fn: fns[6]},
			{Color: btnc, Fn: fns[7]},
		},
	)
	pt.leftbar = NewBar(
		LeftBarX, LeftBarY, LeftBarWidth, LeftBarHeight,
		ButtonWidth, ButtonHeight, ButtonGap, BarPadding, VERTICAL,
		[]*BtnConfig{
			{Color: btnc, Fn: fns[0]},
		},
	)

}

func (pt *Paint) setBackground() {
	bg := &window.Color{R: 25, G: 25, B: 25}
	DrawRect(0, 0, window.ScreenWidth, window.ScreenHeight, bg)

}

func (pt *Paint) Update() (bool, bool) {
	Mouse.Update()
	pt.topbar.Update()
	pt.leftbar.Update()

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
