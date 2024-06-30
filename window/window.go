package window

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	ScreenWidth  int32 = 0
	ScreenHeight int32 = 0
	pixels       []byte
)

type Color struct {
	R, G, B byte
}

type Runnable interface {
	Init(*[]byte)
	Update() bool
	Draw()
}

func Visualise(name string, w, h int32, app Runnable) {
	ScreenHeight = h
	ScreenWidth = w
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Fprintf(os.Stderr, " Failed to initialise SDL : %s\n", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(name, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, " Failed to Create window : %s\n", err)
		os.Exit(2)
	}
	defer window.Destroy()

	renderer, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()

	// look into this pixel format
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, ScreenWidth, ScreenHeight)
	if err != nil {
		fmt.Fprintf(os.Stderr, " Failed to Create TEXTURE : %s\n", err)
		os.Exit(3)
	}
	defer tex.Destroy()
	pixels = make([]byte, ScreenHeight*ScreenWidth*4)

	quit := false
	app.Init(&pixels)

	for !quit {
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			if e.GetType() == sdl.QUIT {
				quit = true
			}
		}
		app.Update()
		app.Draw()

		tex.Update(nil, unsafe.Pointer(&pixels[0]), 4*int(ScreenWidth))
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		sdl.Delay(50)
	}

}

func SetPixel(x, y int, c *Color) {
	index := (x + (int(ScreenWidth) * y)) * 4
	if index+3 <= len(pixels)-1 && index >= 0 {
		pixels[index] = c.R
		pixels[index+1] = c.G
		pixels[index+2] = c.B
		//verify index +3 dedicated for alpha
	}
}
