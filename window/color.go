package window

type Color struct {
	R, G, B byte
}

func (clr *Color) Equals(other *Color) bool {
	return clr.R == other.R && clr.G == other.G && clr.B == other.B
}
