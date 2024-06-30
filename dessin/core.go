package dessin

import (
	"github.com/unk1ndled/draw/mouse"
	"github.com/unk1ndled/draw/primitives"
	"github.com/unk1ndled/draw/window"
)

var (
	Mouse  *mouse.MouseState
	pixels []byte
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
func (pt *Paint) Update() bool {
	Mouse.Update()
	return true
}
func (pt *Paint) Render() {
	pt.canvas.Draw()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	colors = []window.Color{
		{R: 250, G: 0, B: 0},
		{R: 0, G: 255, B: 0},
		{R: 0, G: 0, B: 255}}
)

type Tool byte

const (
	PEN = iota
)

type Canvas struct {

	// pos[0] = x and pos[1] = y
	pos    [2]int32
	width  int32
	height int32

	currentTool Tool
}

func NewCanvas(x, y, w, h int32) *Canvas {
	return &Canvas{pos: [2]int32{x, y}, width: w, height: h, currentTool: PEN}
}

func (cvs *Canvas) isHovered() bool {
	if Mouse.X >= int(cvs.pos[0]) && Mouse.Y > int(cvs.pos[1]) && Mouse.X < int(cvs.width) && Mouse.Y < int(cvs.height) {
		return true
	}
	return false
}

func (cvs *Canvas) Draw() {

	if cvs.isHovered() && Mouse.LeftButton {

		switch cvs.currentTool {
		case PEN:
			primitives.DrawLine(Mouse.PrevX, Mouse.PrevY, Mouse.X, Mouse.Y, &colors[0])
		}
	}

}
