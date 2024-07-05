package dessin

import "github.com/unk1ndled/draw/window"

type OnCLick func()

type Button struct {
	Component
	OnCLick
	color window.Color
}

func NewButton(x0, y0, width, height int32, btnclr window.Color, fn OnCLick) *Button {
	btn := &Button{
		Component: Component{x0, y0, width, height},
		OnCLick:   fn,

		color: btnclr,
	}
	btn.Init()
	return btn
}

func (btn *Button) Init() {

	for i := btn.X; i < btn.Width+btn.X; i++ {
		for j := btn.Y; j < btn.Height+btn.Y; j++ {
			window.SetPixel(i, j, &btn.color)
		}
	}
}
func (btn *Button) ResetVisuals() {

}
