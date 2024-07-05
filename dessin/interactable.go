package dessin

type Component struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

// Contains pixel at (x , y) ?
func (cmpt *Component) contains(x, y int) bool {
	if x >= int(cmpt.X) && y > int(cmpt.Y) && x < int(cmpt.X+cmpt.Width) && y < int(cmpt.Y+cmpt.Height) {
		return true
	}
	return false
}

func (cmpt *Component) isHovered() bool {
	return cmpt.contains(Mouse.X, Mouse.Y)
}

func (cmpt *Component) isClicked() bool {
	return cmpt.isHovered() && Mouse.LeftButton
}

// was clicked last frame ?
func (cmpt *Component) wasClicked() bool {
	return Mouse.PrevLeftButton && cmpt.contains(Mouse.PrevX, Mouse.PrevY)
}

