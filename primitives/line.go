package primitives

import "github.com/unk1ndled/draw/window"

func DrawLine(x1, y1, x2, y2 int, color *window.Color) {

	slope := 0
	if x1 != x2 {
		slope = (y2 - y1) / (x2 - x1)
	}

	for x := x1; x < x2; x++ {
		y := y1 + slope*(x-x1)
		window.SetPixel(x, y, color)
	}
}
