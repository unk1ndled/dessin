package mouse

import "github.com/veandco/go-sdl2/sdl"

type MouseState struct {
	LeftButton      bool
	RightButton     bool
	PrevLeftButton  bool
	PrevRightButton bool
	PrevX, PrevY    int
	X, Y            int
}

func GetMouseState() *MouseState {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()
	var result MouseState
	result.X = int(mouseX)
	result.Y = int(mouseY)
	result.LeftButton = !(leftButton == 0)
	result.RightButton = !(rightButton == 0)
	return &result
}
func (mouseState *MouseState) Update() {
	mouseState.PrevX = mouseState.X
	mouseState.PrevY = mouseState.Y
	mouseState.PrevLeftButton = mouseState.LeftButton
	mouseState.PrevRightButton = mouseState.RightButton

	X, Y, mousebuttonState := sdl.GetMouseState()
	mouseState.X = int(X)
	mouseState.Y = int(Y)
	mouseState.LeftButton = !((mousebuttonState & sdl.ButtonLMask()) == 0)
	mouseState.RightButton = !((mousebuttonState & sdl.ButtonRMask()) == 0)
}
