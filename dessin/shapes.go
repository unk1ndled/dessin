package dessin

import (
	"github.com/unk1ndled/draw/util"
	"github.com/unk1ndled/draw/window"
)

type ShapeType byte

const (
	LINESHAPE ShapeType = iota
	RECTSHAPE
	CIRCLESHAPE
	TRIANGLESHAPE
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
	Draw()
}

func NewShape(base Base, stype ShapeType) Shape {
	switch stype {
	case LINESHAPE:
		return NewLine(base)
	case RECTSHAPE:
		return NewRect(base)
	case CIRCLESHAPE:
		return NewCircle(base)
	}

	return nil
}

type Line struct {
	Base
}

func NewLine(base Base) *Line {
	bs := base
	return &Line{Base: bs}
}

func (l *Line) Draw() {
	RenderLine(l.xStart, l.yStart, l.xEnd, l.yEnd, l.stroke, l.color, DrawWidth)
}

type Rect struct {
	Base
}

func NewRect(base Base) *Rect {
	bs := base
	return &Rect{Base: bs}
}

func (l *Rect) Draw() {
	RenderLine(l.xStart, l.yStart, l.xEnd, l.yStart, l.stroke, l.color, DrawWidth)
	RenderLine(l.xStart, l.yStart, l.xStart, l.yEnd, l.stroke, l.color, DrawWidth)
	RenderLine(l.xEnd, l.yEnd, l.xStart, l.yEnd, l.stroke, l.color, DrawWidth)
	RenderLine(l.xEnd, l.yEnd, l.xEnd, l.yStart, l.stroke, l.color, DrawWidth)
}

type Circle struct {
	Base
}

func NewCircle(base Base) *Circle {
	bs := base
	return &Circle{Base: bs}
}

func (l *Circle) Draw() {
	cx := (l.xStart + l.xEnd) / 2
	cy := (l.yStart + l.yEnd) / 2
	radius := util.Min32(util.Abs32(l.yEnd-l.yStart)/2, util.Abs32(l.xEnd-l.xStart)/2)
	RenderCircle(cx, cy, radius, l.stroke, l.color, DrawWidth)
}
