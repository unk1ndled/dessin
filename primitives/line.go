package primitives

import (
	"math"

	"github.com/unk1ndled/draw/window"
)

func DrawLine(x1, y1, x2, y2 int, color *window.Color) {
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
