package dessin

type ShapeType byte

const (
	Line ShapeType = iota
)

type Shape struct {
	xStart int32
	yStart int32
	xEnd   int32
	yEnd   int32
	stype  ShapeType
}

func NewShape(x1, y1, x2, y2 int32, st ShapeType) *Shape {
	return &Shape{x1, y1, x2, y2, st}
}


