package dessin

import (
	"math"

	"github.com/unk1ndled/draw/util"
	"github.com/unk1ndled/draw/window"
)

type Tool byte
type Operation byte
type Mode byte

var (
	CanvasX int32
	CanvasY int32
	CanvasW int32
	CanvasH int32

	test []byte
)

const (
	MAX_BUFFER_LENGTH = 15

	NONE     byte = iota
	DRAWMODE Mode = iota
	SHAPEMODE

	PEN Tool = iota
	ERASER
	FILL
	LINE

	UNDO Operation = iota
	DRAG
	CLICK
)

type Canvas struct {
	*Component
	Mode
	stype                  ShapeType
	prevCLickX, prevCLickY int32

	buffer *util.Deque[[]byte]

	currentTool Tool
	drawColor   *window.Color
	currOp      Operation

	lineWidth   int32
	strokeWidth int32
}

func NewCanvas(x, y, w, h int32) *Canvas {

	CanvasX = x
	CanvasY = y
	CanvasW = w
	CanvasH = h

	cvs := &Canvas{
		Component:   &Component{X: x, Y: y, Width: w, Height: h},
		buffer:      util.NewDeque[[]byte](),
		drawColor:   colors[0],
		currentTool: PEN, lineWidth: 20, strokeWidth: 8,
	}

	test = make([]byte, len(pixels))

	RenderRect(x, y, x+w, y+h, &window.Color{R: 0, G: 0, B: 0})

	return cvs
}

func (cvs *Canvas) Update() bool {

	if CTRL && !Z_PRESS && Z_PREV_PRESS {
		Z_PREV_PRESS = false
		cvs.currOp = UNDO
		return true
	} else if cvs.isPressed() {
		if !Mouse.PrevLeftButton {
			cvs.updateBuffer()
		}
		cvs.currOp = DRAG
		return true
	} else if cvs.isClicked() {
		cvs.currOp = CLICK
		return true
	}

	return false
}

// is called if Update is true
func (cvs *Canvas) Render() {

	switch cvs.Mode {
	case DRAWMODE:
		switch cvs.currOp {
		case UNDO:
			cvs.undo()
			cvs.currOp = Operation(NONE)
		case CLICK:
			switch cvs.currentTool {
			case FILL:
				clicked := window.GetPixelColor(Mouse.X, Mouse.Y)
				cvs.Fill(Mouse.X, Mouse.Y, cvs.drawColor, &clicked)
			}
		case DRAG:
			switch cvs.currentTool {
			case PEN:
				RenderLine(Mouse.PrevX, Mouse.PrevY, Mouse.X, Mouse.Y, cvs.lineWidth, cvs.drawColor, DrawWidth)
			case ERASER:
				cvs.Erase(Mouse.PrevX, Mouse.PrevY, Mouse.X, Mouse.Y)
			}
		}
	case SHAPEMODE:
		switch cvs.currOp {
		case UNDO:
			cvs.undo()
			cvs.currOp = Operation(NONE)
		case DRAG:
			if cvs.initialPress() {
				cvs.prevCLickX = Mouse.X
				cvs.prevCLickY = Mouse.Y
				copy(test, pixels)
			} else {

				if !(Mouse.PrevX == Mouse.X && Mouse.PrevY == Mouse.Y) {
					copy(pixels, test)
					base := Base{Mouse.X, Mouse.Y, cvs.prevCLickX, cvs.prevCLickY, cvs.strokeWidth, cvs.drawColor}
					shape := NewShape(base, cvs.stype)
					if shape != nil {
						shape.Draw()
					}
				}

			}
		}
	}

}

// this function is called before a pixels state mutation occurs to save the prev state
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

// sets current tool
func (cvs *Canvas) setTool(t Tool) {
	cvs.currentTool = t
	cvs.Mode = DRAWMODE
}

func (cvs *Canvas) setShape(t ShapeType) {
	cvs.stype = t
	cvs.Mode = SHAPEMODE
}

func (cvs *Canvas) modifyDrawWidth(factor int32) {
	cvs.lineWidth += factor
	cvs.lineWidth = int32(math.Max(0, float64(cvs.lineWidth)))
	cvs.strokeWidth = int32(math.Min(8, float64(cvs.lineWidth)))
}

// sets draw color
func (cvs *Canvas) setColor(index int) {
	cvs.drawColor = colors[index]
}

// flood fill algo
func (cvs *Canvas) Fill(x0, y0 int32, fillColor, clickedColor *window.Color) {
	if cvs.contains(x0, y0) {
		curPixelCLr := window.GetPixelColor(x0, y0)
		if curPixelCLr.Equals(fillColor) {
			return
		}
		if curPixelCLr.Equals(clickedColor) {
			window.SetPixel(x0, y0, fillColor)
			cvs.Fill(x0+1, y0, fillColor, clickedColor)
			cvs.Fill(x0-1, y0, fillColor, clickedColor)
			cvs.Fill(x0, y0-1, fillColor, clickedColor)
			cvs.Fill(x0, y0+1, fillColor, clickedColor)
		}
	}
}

// controls line drawing width
func DrawWidth(x, y, width int32, color *window.Color) {
	halfWidth := width / 2
	isEven := width%2 == 0
	offset := int32(0)

	if !isEven {
		offset = 1
	}

	xStart := util.Max32(x-halfWidth, CanvasX)
	xEnd := util.Min32(x+halfWidth+offset, CanvasX+CanvasW)
	yStart := util.Max32(y-halfWidth, CanvasY)
	yEnd := util.Min32(y+halfWidth+offset, CanvasY+CanvasH)

	for i := xStart; i <= xEnd; i++ {
		for j := yStart; j <= yEnd; j++ {
			window.SetPixel(i, j, color)
		}
	}
}

func (cvs *Canvas) Erase(x0, y0, x1, y1 int32) {
	RenderLine(x0, y0, x1, y1, cvs.lineWidth, &basecolor, DrawWidth)
}
