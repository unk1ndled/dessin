package primitives

import (
	"math"

	"github.com/unk1ndled/draw/window"
)

// Bresenham's Line Algorithm
// explanation : https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
func DrawLine(x1, y1, x2, y2 int, color *window.Color) {
	dx := int(math.Abs(float64(x2 - x1)))
	sx := -1
	if x1 < x2 {
		sx = 1
	}
	dy := -int(math.Abs(float64(y2 - y1)))
	sy := -1
	if y1 < y2 {
		sy = 1
	}
	err := dx + dy

	for {
		// window.SetPixel(x1+1, y1-1, color)
		window.SetPixel(x1+1, y1, color)
		// window.SetPixel(x1+1, y1+1, color)
		window.SetPixel(x1, y1-1, color)
		window.SetPixel(x1, y1, color)
		window.SetPixel(x1, y1+1, color)
		// window.SetPixel(x1-1, y1-1, color)
		window.SetPixel(x1-1, y1, color)
		// window.SetPixel(x1-1, y1+1, color)

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

// Digital differential analyzer
func DrawLineOld(x1, y1, x2, y2 int, color *window.Color) {
	dx := x2 - x1
	dy := y2 - y1

	var step float64
	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		step = math.Abs(float64(dx))
	} else {
		step = math.Abs(float64(dy))
	}
	xInc := float64(dx) / step
	yInc := float64(dy) / step

	x, y := x1, y1
	for i := 0; i < int(step); i++ {
		window.SetPixel(x, y, color)
		x += int(xInc)
		y += int(yInc)
	}

}
