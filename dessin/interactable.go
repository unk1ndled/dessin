package dessin

import (
	"github.com/unk1ndled/draw/icons"
	"github.com/unk1ndled/draw/window"
)

type Component struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

// Contains pixel at (x , y) ?
func (cmpt *Component) contains(x, y int32) bool {
	if x >= (cmpt.X) && y >= (cmpt.Y) && x <= (cmpt.X+cmpt.Width) && y <= (cmpt.Y+cmpt.Height) {
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

//////////////////////////////////////////////////////////////////////////////////////////////////////////

type BarType byte

const (
	VERTICAL BarType = iota
	HORIZONTAL
	GRID
)

type BtnConfig struct {
	Color      *window.Color
	Fn         func()
	ButtonIcon icons.IconType
}

type Bar struct {
	Component
	bartype BarType
	buttons []*Button
}

func NewBar(x, y, w, h, btnW, btnH, gap, padding int32, btype BarType, btns []*BtnConfig) *Bar {
	tp := &Bar{
		Component: Component{x, y, w, h},
		buttons:   make([]*Button, len(btns)),
		bartype:   btype,
	}
	btnx, btny := x+padding, y+padding
	if btype == GRID {
		btnH = int32(btnH / 2)
		btnW = int32(btnW / 2)
		j := int32(0)
		k := 0
		for j = 0; j < 2; j++ {
			btny = btny - padding/2 + j*(btnH+gap)
			i := int32(0)
			for i = 0; i < int32(len(btns)/2); i++ {
				tp.buttons[k] = NewButton(
					btnx+(i)*(gap+btnW),
					btny,
					btnW, btnH,
					btns[k].Color, btns[k].Fn)
				k++
			}
		}
		tp.Width = ((gap + btnW) * int32(len(btns)/2)) + padding
		return tp
	}
	i := int32(0)
	xdir, ydir := int32(1), int32(0)
	if btype == VERTICAL {
		xdir = 0
		ydir = 1
	}

	for i = 0; i < int32(len(btns)); i++ {
		var ic *icons.Icon
		if btns[i].ButtonIcon != icons.NONE {
			ic = icons.NewIcon(btns[i].ButtonIcon)
		}
		tp.buttons[i] = NewButton(
			btnx+xdir*(i)*(gap+btnW),
			btny+ydir*(i)*(gap+btnH),
			btnW, btnH,
			btns[i].Color, btns[i].Fn, ic)
	}

	tp.Width = ((gap + btnW) * int32(len(btns))) + padding
	return tp
}

func (bar *Bar) Update() bool {
	update := false
	for _, btn := range bar.buttons {
		update = btn.Update() || update
	}
	// if bar.inner != nil {
	// 	bar.inner.Update()
	// }
	return update
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////

type OnCLick func()

const (
	BtnBorder = 3
)

type Button struct {
	*Component
	OnCLick

	mouseEnter bool
	icon       *icons.Icon

	baseColor *window.Color
	light     *window.Color
	verylight *window.Color
	dark      *window.Color
}

func NewButton(x0, y0, width, height int32, btnclr *window.Color, fn OnCLick, icons ...*icons.Icon) *Button {
	btn := &Button{
		Component: &Component{x0, y0, width, height},
		OnCLick:   fn,
		baseColor: btnclr,
		verylight: btnclr.Lighter(50),
		light:     btnclr.Lighter(15),
		dark:      btnclr.Darker(20),
	}
	if len(icons) != 0 {
		btn.icon = icons[0]
	}
	btn.Init()
	return btn
}

func (btn *Button) Init() {
	btn.ResetVisuals()
}

func (btn *Button) Update() bool {
	if btn.isHovered() {
		if !Mouse.LeftButton {
			if !btn.mouseEnter {
				btn.mouseEnter = true
				// log.Println("entered")
			} else if btn.wasPressed() {
				btn.OnCLick()
				btn.ResetVisuals()
				// log.Println("clicked")
				return true
			}
		} else
		// leftclick and hover
		if btn.mouseEnter {
			if !btn.wasPressed() {
				btn.clickVisuals()
				// log.Println("painted")
				return true
			}
			// log.Println("pressing")
		}
	} else {
		if btn.mouseEnter {
			btn.ResetVisuals()
			btn.mouseEnter = false
		}
	}
	return false
}

func (btn *Button) clickVisuals() {
	btn.Render(btn.dark, btn.baseColor, btn.light)
}

func (btn *Button) ResetVisuals() {
	btn.Render(btn.baseColor, btn.light, btn.verylight)
}

func (btn *Button) Render(dark, base, light *window.Color) {
	//left
	RenderRect(btn.X, btn.Y, btn.X+BtnBorder, btn.Y+btn.Height, light)
	//right
	RenderRect(btn.X+btn.Width-BtnBorder, btn.Y, btn.X+btn.Width, btn.Y+btn.Height, base)

	//top
	RenderRect(btn.X+BtnBorder, btn.Y, btn.X+btn.Width-BtnBorder, btn.Y+BtnBorder, light)
	//bot
	RenderRect(btn.X+BtnBorder, btn.Y+btn.Height-BtnBorder, btn.X+btn.Width-BtnBorder, btn.Y+btn.Height, base)

	//center
	RenderRect(btn.X+BtnBorder, btn.Y+BtnBorder, btn.X+btn.Width-BtnBorder, btn.Y+btn.Height-BtnBorder, dark)

	if btn.icon != nil {
		btn.DrawIcon(light)
	}
}

func (btn *Button) DrawIcon(clr *window.Color) {
	if btn.icon != nil {
		scale := int32((btn.Height / icons.BITMAPWIDTH))
		iconw := scale * icons.BITMAPWIDTH
		iconX, iconY := 1+btn.X+(btn.Width-iconw)/2, 1+btn.Y+(btn.Height-iconw)/2

		for i := int32(1); i < int32(len(*btn.icon.Data)-1); i++ {
			if (*btn.icon.Data)[i] == 1 {
				ix, iy := scale*(i%icons.BITMAPWIDTH), scale*(i/icons.BITMAPWIDTH)
				for j := int32(0); j < scale; j++ {
					for k := int32(0); k < scale; k++ {
						window.SetPixel(iconX+ix+j, iconY+iy+k, clr)
					}
				}
			}
		}
	}
}
