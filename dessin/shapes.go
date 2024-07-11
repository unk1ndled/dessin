package dessin

import "github.com/unk1ndled/draw/window"

type ShapeType byte

const (
	LINESHAPE ShapeType = iota
)

type Base struct {
	xStart int32
	yStart int32
	xEnd   int32
	yEnd   int32

	stroke int32
	color  *window.Color
}

type Shape interface {
	Draw(pixels *[]byte)
}

type Line struct {
	Base
}

func NewLine(x1, y1, x2, y2, stroke int32, c *window.Color) *Line {
	bs := Base{x1, y1, x2, y2, stroke, c}
	return &Line{Base: bs}
}

func (l Line) Draw() {
	RenderLine(l.xStart, l.yStart, l.xEnd, l.yEnd, l.stroke, l.color, nil)
}
