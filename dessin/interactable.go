package dessin

import "github.com/unk1ndled/draw/window"

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

func (cmpt *Component) isClicked() bool {
	return cmpt.wasPressed() && !Mouse.LeftButton

}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
type OnCLick func()

const (
	BtnBorder = 3
)

type Button struct {
	Component
	OnCLick

	baseColor *window.Color
	light     *window.Color
	verylight *window.Color
	dark      *window.Color
}

func NewButton(x0, y0, width, height int32, btnclr *window.Color, fn OnCLick) *Button {
	btn := &Button{
		Component: Component{x0, y0, width, height},
		OnCLick:   fn,

		baseColor: btnclr,
		verylight: btnclr.Lighter(50),
		light:     btnclr.Lighter(15),
		dark:      btnclr.Darker(20),
	}
	btn.Init()
	return btn
}

func (btn *Button) Init() {
	btn.ResetVisuals()
}

func (btn *Button) Update() {

	if btn.isPressed() {
		btn.clickVisuals()

	} else if btn.wasPressed() {
		btn.OnCLick()
		btn.ResetVisuals()
	}
}

func (btn *Button) clickVisuals() {
	//left
	DrawRect(btn.X, btn.Y, btn.X+BtnBorder, btn.Y+btn.Height, btn.baseColor)
	//right
	DrawRect(btn.X+btn.Width-BtnBorder, btn.Y, btn.X+btn.Width, btn.Y+btn.Height, btn.baseColor)

	//top
	DrawRect(btn.X+BtnBorder, btn.Y, btn.X+btn.Width-BtnBorder, btn.Y+BtnBorder, btn.baseColor)
	//bot
	DrawRect(btn.X+BtnBorder, btn.Y+btn.Height-BtnBorder, btn.X+btn.Width-BtnBorder, btn.Y+btn.Height, btn.baseColor)

	//center
	DrawRect(btn.X+BtnBorder, btn.Y+BtnBorder, btn.X+btn.Width-BtnBorder, btn.Y+btn.Height-BtnBorder, btn.dark)

}

func (btn *Button) ResetVisuals() {
	//left
	DrawRect(btn.X, btn.Y, btn.X+BtnBorder, btn.Y+btn.Height, btn.verylight)
	//right
	DrawRect(btn.X+btn.Width-BtnBorder, btn.Y, btn.X+btn.Width, btn.Y+btn.Height, btn.light)

	//top
	DrawRect(btn.X+BtnBorder, btn.Y, btn.X+btn.Width-BtnBorder, btn.Y+BtnBorder, btn.verylight)
	//bot
	DrawRect(btn.X+BtnBorder, btn.Y+btn.Height-BtnBorder, btn.X+btn.Width-BtnBorder, btn.Y+btn.Height, btn.light)

	//center
	DrawRect(btn.X+BtnBorder, btn.Y+BtnBorder, btn.X+btn.Width-BtnBorder, btn.Y+btn.Height-BtnBorder, btn.baseColor)

}
