package dessin

import "github.com/unk1ndled/draw/window"

type OnCLick func()

type Button struct {
	Component
	OnCLick
}

func NewButton(x0, y0, width, height int32, color window.Color, fn OnCLick) *Button {
	return &Button{
		Component: Component{x0, y0, width, height},
		OnCLick:   fn,
	}
}

func (btn *Button) Init() {

}
func (btn *Button) ResetVisuals() {

}
func (btn *Button) Clicked() {
	if btn.wasClicked()
}
