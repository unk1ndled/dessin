package dessin

import "github.com/unk1ndled/draw/window"

type TopBar struct {
	buttons []*Button
}

func NewTopBar(x, y, btnheight int32, fns []func()) *TopBar {
	tp := &TopBar{
		buttons: make([]*Button, len(fns)),
	}
	gap := int32(10)
	for i := 0; i < len(fns); i++ {
		tp.buttons[i] = NewButton(x+int32(i)*(gap+btnheight), y,
			btnheight, btnheight,
			&window.Color{R: 40, G: 40, B: 40}, fns[i])
	}
	return tp
}

func (tb *TopBar) Update() bool {
	update := false
	for _, btn := range tb.buttons {
		update = btn.Update() || update
	}
	return update
}
