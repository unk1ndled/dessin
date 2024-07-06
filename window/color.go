package window

import "math"

type Color struct {
	R, G, B byte
}

func (clr *Color) Darker(val byte) *Color {

	dec := byte(math.Min(float64(val), float64(clr.GetMin())))
	return &Color{clr.R - dec, clr.G - dec, clr.B - dec}
}

func (clr *Color) Lighter(val byte) *Color {
	inc := byte(math.Min(float64(val), float64(255-clr.GetMAx())))
	return &Color{clr.R + inc, clr.G + inc, clr.B + inc}

}

func (clr *Color) GetMin() byte {
	return byte(math.Min(math.Min(float64(clr.R), float64(clr.G)), float64(clr.B)))
}

func (clr *Color) GetMAx() byte {
	return byte(math.Max(math.Max(float64(clr.R), float64(clr.G)), float64(clr.B)))
}
func (clr *Color) Equals(other *Color) bool {
	return clr.R == other.R && clr.G == other.G && clr.B == other.B
}
