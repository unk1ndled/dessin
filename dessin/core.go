package dessin

import (
	"fmt"

	"github.com/unk1ndled/draw/window"
	"github.com/veandco/go-sdl2/sdl"
)

// implements sdl.Runnable
type Canvas struct {
	pixels *[]byte
}

func NewCanvas() *Canvas {
	return &Canvas{}
}
func (Canvas *Canvas) Init(p *[]byte) {
	Canvas.pixels = p
}
func (Canvas *Canvas) Update() bool {
	return true
}
func (Canvas *Canvas) Draw() {
	fmt.Println("drew")

	x, y, _ := sdl.GetMouseState()
	window.SetPixel(int(x), int(y), &window.Color{R: 255, G: 0, B: 0})
}
