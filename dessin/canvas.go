package dessin

import (
	"math"

	"github.com/unk1ndled/draw/util"
	"github.com/unk1ndled/draw/window"
)

type Tool byte
type Operation byte

const (
	MAX_BUFFER_LENGTH = 15

	PEN Tool = iota
	ERASER
	FILL
	LINE

	UNDO Operation = iota
	DRAW
	CLICK
)

type Canvas struct {
	*Component
	buffer *util.Deque[[]byte]

	prevCLickX, prevCLickY int32

	currentTool Tool
	drawColor   *window.Color
	currOp      Operation

	lineWidth int32
}

func NewCanvas(x, y, w, h int32) *Canvas {

	cvs := &Canvas{
		Component:   &Component{X: x, Y: y, Width: w, Height: h},
		buffer:      util.NewDeque[[]byte](),
		drawColor:   colors[0],
		currentTool: PEN, lineWidth: 20,
	}

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
		cvs.currOp = DRAW
		return true
	} else if cvs.isClicked() {
		cvs.currOp = CLICK
		return true
	}
	return false
}

// is called if Update is true
func (cvs *Canvas) Render() {
	switch cvs.currOp {
	case UNDO:
		cvs.undo()
	case CLICK:
		switch cvs.currentTool {
		case FILL:
			clicked := window.GetPixelColor(Mouse.X, Mouse.Y)
			cvs.Fill(Mouse.X, Mouse.Y, cvs.drawColor, &clicked)
		}
	case DRAW:
		switch cvs.currentTool {
		case PEN:
			RenderLine(Mouse.PrevX, Mouse.PrevY, Mouse.X, Mouse.Y, cvs.lineWidth, cvs.drawColor, cvs.DrawWidth)
		case ERASER:
			cvs.Erase(Mouse.PrevX, Mouse.PrevY, Mouse.X, Mouse.Y)
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
}

func (cvs *Canvas) modifyDrawWidth(factor int32) {
	cvs.lineWidth += factor
	cvs.lineWidth = int32(math.Max(0, float64(cvs.lineWidth)))
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
func (cvs *Canvas) DrawWidth(x, y, width int32, color *window.Color) {
	halfWidth := width / 2
	isEven := width%2 == 0
	offset := int32(0)

	if !isEven {
		offset = 1
	}

	xStart := util.Max32(x-halfWidth, cvs.X)
	xEnd := util.Min32(x+halfWidth+offset, cvs.X+cvs.Width)
	yStart := util.Max32(y-halfWidth, cvs.Y)
	yEnd := util.Min32(y+halfWidth+offset, cvs.Y+cvs.Height)

	for i := xStart; i <= xEnd; i++ {
		for j := yStart; j <= yEnd; j++ {
			window.SetPixel(i, j, color)
		}
	}
}

func (cvs *Canvas) Erase(x0, y0, x1, y1 int32) {
	RenderLine(x0, y0, x1, y1, cvs.lineWidth, &basecolor, cvs.DrawWidth)
}
