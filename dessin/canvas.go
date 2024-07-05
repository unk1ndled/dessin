package dessin

import (
	"log"
	"math"

	"github.com/unk1ndled/draw/util"
	"github.com/unk1ndled/draw/window"
)

type Tool byte
type Operation byte

const (
	MAX_BUFFER_LENGTH = 15

	PEN = iota
	ERASER
	FILL

	UNDO = iota
	DRAW
	CLICK
)

type Canvas struct {
	*Component
	buffer *util.Deque[[]byte]

	currentTool Tool
	currOp      Operation

	lineWidth int32
}

func NewCanvas(x, y, w, h int32) *Canvas {
	return &Canvas{
		Component:   &Component{X: x, Y: y, Width: w, Height: h},
		buffer:      util.NewDeque[[]byte](),
		currentTool: PEN, lineWidth: 20}
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
	} else if cvs.Click() {
		cvs.currOp = CLICK
		return true
	}

	return false
}

func (cvs *Canvas) Draw() {

	switch cvs.currOp {
	case UNDO:
		cvs.undo()
		log.Println("undind")

	case CLICK:
		switch cvs.currentTool {
		case FILL:
			clicked := window.GetPixelColor(Mouse.X, Mouse.Y)
			cvs.Fill(Mouse.X, Mouse.Y, &colors[2], &clicked)
			log.Println("filled")

		}
	case DRAW:
		switch cvs.currentTool {
		case PEN:
			cvs.DrawLine(Mouse.PrevX, Mouse.PrevY, Mouse.X, Mouse.Y, cvs.lineWidth, &colors[1])
		case ERASER:
			cvs.Erase(Mouse.PrevX, Mouse.PrevY, Mouse.X, Mouse.Y)
		}
		log.Println("drew")

	}
}

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

func (cvs *Canvas) setTool(t Tool) {
	cvs.currentTool = t
}

// func (cvs *Canvas) ValidateX(x int) int {
// 	if x < int(cvs.X) {
// 		return int(cvs.X)
// 	} else if x > int(cvs.X)+int(cvs.Width) {
// 		return int(cvs.X + cvs.Width)
// 	} else {
// 		return x
// 	}
// }
// func (cvs *Canvas) ValidateY(y int) int {
// 	if y < int(cvs.Y) {
// 		return int(cvs.Y)
// 	} else if y > int(cvs.Y)+int(cvs.Height) {
// 		return int(cvs.Y + cvs.Height)
// 	} else {
// 		return y
// 	}
// }

func (cvs *Canvas) Fill(x0, y0 int32, fillColor, clickedColor *window.Color) {
	if x0 >= (cvs.X) && x0 <= (cvs.X)+(cvs.Width) && y0 >= (cvs.Y) && y0 <= (cvs.Y+cvs.Height) {
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

// Bresenham's Line Algorithm
// explanation : https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
func (cvs *Canvas) DrawLine(x1, y1, x2, y2, width int32, color *window.Color) {
	dx := int32(math.Abs(float64(x2 - x1)))
	sx := int32(-1)
	if x1 < x2 {
		sx = 1
	}
	dy := -int32(math.Abs(float64(y2 - y1)))
	sy := int32(-1)
	if y1 < y2 {
		sy = 1
	}
	err := dx + dy

	for {

		go cvs.DrawWidth(x1, y1, width, color)

		if x1 == x2 && y1 == y2 {
			break
		}
		err2 := 2 * err
		if err2 >= dy {
			if x1 == x2 {
				break
			}
			err += dy
			x1 += sx
		}
		if err2 <= dx {
			if y1 == y2 {
				break
			}
			err += dx
			y1 += sy
		}
	}
}

func (cvs *Canvas) DrawWidth(x1, y1, width int32, color *window.Color) {

	xstart, xend := int32(math.Max(float64(x1-width), float64(cvs.X))), int32(math.Min(float64(x1+width), float64(cvs.X+cvs.Width)))
	ystart, yend := int32(math.Max(float64(y1-width), float64(cvs.Y))), int32(math.Min(float64(y1+width), float64(cvs.Y+cvs.Height)))

	for x := xstart; x <= xend; x++ {
		for y := ystart; y <= yend; y++ {
			window.SetPixel(x, y, color)
		}
	}
}

func (cvs *Canvas) Erase(x0, y0, x1, y1 int32) {
	cvs.DrawLine(x0, y0, x1, y1, cvs.lineWidth, &basecolor)
}

// Digital differential analyzer
// func DrawLineOld(x1, y1, x2, y2 int, color *window.Color) {
// 	dx := x2 - x1
// 	dy := y2 - y1

// 	var step float64
// 	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
// 		step = math.Abs(float64(dx))
// 	} else {
// 		step = math.Abs(float64(dy))
// 	}
// 	xInc := float64(dx) / step
// 	yInc := float64(dy) / step

// 	x, y := x1, y1
// 	for i := 0; i < int(step); i++ {
// 		window.SetPixel(x, y, color)
// 		x += int(xInc)
// 		y += int(yInc)
// 	}

// }
