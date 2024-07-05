package window

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	ScreenWidth  int32
	ScreenHeight int32
	pixels       []byte
)

func GetPixelsIndex(x, y int) int {
	return (x + (int(ScreenWidth) * y)) * 4
}
func GetPixelColor(x, y int) Color {
	index := GetPixelsIndex(x, y)
	return Color{pixels[index], pixels[index+1], pixels[index+2]}
}

func SetPixel(x, y int, c *Color) {

	index := GetPixelsIndex(x, y)
	if index+3 <= len(pixels)-1 && index >= 0 {
		pixels[index] = c.R
		pixels[index+1] = c.G
		pixels[index+2] = c.B
	}
}

type Runnable interface {
	Init([]byte)
	Update() (bool, bool)
	Render()
}

func Visualise(name string, w, h int32, app Runnable) {
	ScreenWidth = w
	ScreenHeight = h

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(name, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		os.Exit(2)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		os.Exit(3)
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, ScreenWidth, ScreenHeight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		os.Exit(4)
	}
	defer tex.Destroy()

	pixels = make([]byte, ScreenHeight*ScreenWidth*4)

	app.Init(pixels)
	quit := false
	update := false
	for !quit {

		update, quit = app.Update()
		if update {
			app.Render()

		}
		tex.Update(nil, unsafe.Pointer(&pixels[0]), 4*int(ScreenWidth))
		renderer.Clear()
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		//sdl.Delay(16)
	}
}
