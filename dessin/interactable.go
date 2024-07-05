package dessin

type Component struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

// Contains pixel at (x , y) ?
func (cmpt *Component) contains(x, y int32) bool {
	if x >= (cmpt.X) && y > (cmpt.Y) && x < (cmpt.X+cmpt.Width) && y < (cmpt.Y+cmpt.Height) {
		return true
	}
	return false
}

func (cmpt *Component) isHovered() bool {
	return cmpt.contains(Mouse.X, Mouse.Y)
}

func (cmpt *Component) isPressed() bool {
	return cmpt.isHovered() && Mouse.LeftButton
}

// was Pressed last frame ?
func (cmpt *Component) wasPressed() bool {
	return Mouse.PrevLeftButton && cmpt.contains(Mouse.PrevX, Mouse.PrevY)
}

func (cmpt *Component) Click() bool {
	return cmpt.wasPressed() && !Mouse.LeftButton

}
