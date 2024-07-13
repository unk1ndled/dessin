package dessin

import (
	"github.com/unk1ndled/draw/util"
	"github.com/unk1ndled/draw/window"
)

type widthFunc func(x, y, width int32, clr *window.Color)

// Bresenham's Line Algorithm
// explanation : https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
func RenderLine(x1, y1, x2, y2, width int32, color *window.Color, fn widthFunc) {
	dx := (util.Abs32(x2 - x1))
	sx := int32(-1)
	if x1 < x2 {
		sx = 1
	}
	dy := -(util.Abs32(y2 - y1))
	sy := int32(-1)
	if y1 < y2 {
		sy = 1
	}
	err := dx + dy

	for {

		go fn(x1, y1, width, color)

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

func RenderRect(x1, y1, x2, y2 int32, clr *window.Color) {
	for i := x1; i <= x2; i++ {
		for j := y1; j <= y2; j++ {
			window.SetPixel(i, j, clr)
		}
	}
}

func RenderCircle(centerX, centerY, radius, stroke int32, color *window.Color, fn widthFunc) {
	x := radius - 1
	y := int32(0)
	dx := int32(2)
	dy := int32(2)
	err := dx - (radius << 1)

	for x >= y {
		window.SetPixel(centerX+x, centerY+y, color)
		go fn(centerX+x, centerY+y, stroke, color)

		window.SetPixel(centerX+y, centerY+x, color)
		go fn(centerX+y, centerY+x, stroke, color)

		window.SetPixel(centerX-y, centerY+x, color)
		go fn(centerX-y, centerY+x, stroke, color)

		window.SetPixel(centerX-x, centerY+y, color)
		go fn(centerX-x, centerY+y, stroke, color)

		window.SetPixel(centerX-x, centerY-y, color)
		go fn(centerX-x, centerY-y, stroke, color)

		window.SetPixel(centerX-y, centerY-x, color)
		go fn(centerX-y, centerY-x, stroke, color)

		window.SetPixel(centerX+y, centerY-x, color)
		go fn(centerX+y, centerY-x, stroke, color)

		window.SetPixel(centerX+x, centerY-y, color)
		go fn(centerX+x, centerY-y, stroke, color)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (int32(radius) << 1)
		}
	}
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
