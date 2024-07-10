package dessin

import (
	"github.com/unk1ndled/draw/icons"
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

	colors = []*window.Color{

		{R: 58, G: 135, B: 47},
		{R: 102, G: 204, B: 102},
		{R: 204, G: 230, B: 102},
		{R: 248, G: 229, B: 65},
		{230, 150, 80},
		///
		{230, 50, 80},
		{135, 53, 85},
		{166, 85, 95},
		{178, 102, 79},
		{242, 174, 153},
		//
		{102, 178, 255},
		{102, 102, 255},
		{153, 102, 255},
		{131, 77, 196},
		{125, 45, 160},
		///

		{255, 153, 204},
		{238, 143, 238},
		{255, 255, 255},
		{127, 127, 127},
		{10, 10, 10},
	}
)

// implements sdl.Runnable
// Paint acts as the wrapper struct for the program
// thus its responsible for the inner components
type Paint struct {
	topbars []*Bar
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

	ButtonWidth := TopBarHeight + int32(TopBarHeight/4)
	ButtonHeight := TopBarHeight - 2*BarPadding
	ButtonGap := Padding - 1

	LeftBarX := Padding
	LeftBarY := TopBarHeight + 2*Padding
	LeftBarWidth := ButtonWidth + BarPadding
	LeftBarHeight := window.ScreenHeight - (2*Padding + TopBarHeight)

	pixels = pxls
	Mouse = mouse.GetMouseState()
	pt.setBackground()

	pt.canvas = NewCanvas(
		(2*Padding)+LeftBarWidth,
		(2*Padding)+TopBarHeight,
		window.ScreenWidth-((2*Padding)+LeftBarWidth+10),
		window.ScreenHeight-((2*Padding)+TopBarHeight+10)) // Corrected from LeftBarHeight to TopBarHeight

	options := []func(){
		func() { pt.canvas.setTool(PEN) },
		func() { pt.canvas.setTool(ERASER) },
		func() { pt.canvas.modifyDrawWidth(-1) },
		func() { pt.canvas.modifyDrawWidth(1) },
	}
	btnc := &window.Color{R: 40, G: 40, B: 40}
	pt.topbars = make([]*Bar, 3)
	pt.topbars[0] = NewBar(
		TopBarX, TopBarY, TopBarWidth, TopBarHeight,
		ButtonWidth, ButtonHeight, ButtonGap, BarPadding, HORIZONTAL,
		[]*BtnConfig{
			{Color: btnc, Fn: options[0], ButtonIcon: icons.PEN},
			{Color: btnc, Fn: options[1], ButtonIcon: icons.ERASER},
			{Color: btnc, Fn: options[2], ButtonIcon: icons.DECREASE},
			{Color: btnc, Fn: options[3], ButtonIcon: icons.INCREASE},
		},
	)

	clrconfig := *(new([]*BtnConfig))
	for i := 0; i < 20; i++ {
		//closure stuff with using i inside the functions
		index := i
		clrconfig = append(clrconfig, &BtnConfig{Color: colors[i], Fn: func() { pt.canvas.setColor(index) }})
	}

	pt.topbars[1] = NewBar(
		pt.topbars[0].X+pt.topbars[0].Width+Padding, TopBarY, TopBarWidth, TopBarHeight,
		ButtonWidth, ButtonHeight, ButtonGap, BarPadding, GRID,
		clrconfig,
	)

	pt.topbars[2] = NewBar(
		pt.topbars[1].X+pt.topbars[1].Width+Padding, TopBarY, TopBarWidth, TopBarHeight,
		ButtonWidth, ButtonHeight, ButtonGap, BarPadding, HORIZONTAL,
		[]*BtnConfig{
			{Color: btnc, Fn: func() {
				window.OpenPNG(pt.canvas.X, pt.canvas.Y, pt.canvas.Width, pt.canvas.Height)
			}, ButtonIcon: icons.OPENFILE},
			{Color: btnc, Fn: func() {
				window.SaveTextureAsImage(pt.canvas.X, pt.canvas.Y, pt.canvas.Width, pt.canvas.Height, "drawing")
			}, ButtonIcon: icons.SAVE}},
	)
	tools := []func(){
		func() { pt.canvas.setTool(FILL) }}

	pt.leftbar = NewBar(
		LeftBarX, LeftBarY, LeftBarWidth, LeftBarHeight,
		ButtonWidth, ButtonHeight, ButtonGap, BarPadding, VERTICAL,
		[]*BtnConfig{
			{Color: btnc, Fn: tools[0], ButtonIcon: icons.FILL},
		},
	)

}

func (pt *Paint) setBackground() {
	bg := &window.Color{R: 25, G: 25, B: 25}
	RenderRect(0, 0, window.ScreenWidth, window.ScreenHeight, bg)

}

func (pt *Paint) Update() (bool, bool) {
	Mouse.Update()
	update := pt.topbars[0].Update()
	update = update || pt.topbars[1].Update()
	update = update || pt.topbars[2].Update()
	update = update || pt.leftbar.Update()
	update = update || pt.canvas.Update()
	quit := pt.checkEvent()
	return update, quit
}
func (pt *Paint) Render() {
	pt.canvas.Render()
}

func (pt *Paint) checkEvent() bool {
	for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		switch t := e.(type) {
		case *sdl.QuitEvent:
			return true
		case *sdl.KeyboardEvent:
			pt.HandleKeyBoard(t)
		}
	}
	return false
}

func (pt *Paint) HandleKeyBoard(t *sdl.KeyboardEvent) {
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
