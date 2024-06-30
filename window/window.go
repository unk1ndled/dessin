package window

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/unk1ndled/draw/util"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	ScreenWidth  int32
	ScreenHeight int32
	pixels       []byte
	buffer       *util.Deque[[]byte]
)

type Color struct {
	R, G, B byte
}

func SetPixel(x, y int, c *Color) {
	index := (x + (int(ScreenWidth) * y)) * 4
	if index+3 <= len(pixels)-1 && index >= 0 {
		pixels[index] = c.R
		pixels[index+1] = c.G
		pixels[index+2] = c.B
	}
}

type Runnable interface {
	Init([]byte)
	Update() bool
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
	buffer = util.NewDeque[[]byte]()

	app.Init(pixels)
	quit := false
	ctrlPressed, zPressed := false, false

	prevz := false
	prevctrl := false

	for !quit {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch t := e.(type) {
			case *sdl.QuitEvent:
				quit = true
			case *sdl.KeyboardEvent:
				switch t.Type {
				case sdl.KEYDOWN:
					switch t.Keysym.Scancode {
					case sdl.SCANCODE_LCTRL, sdl.SCANCODE_RCTRL:
						prevctrl = ctrlPressed
						ctrlPressed = true
					case sdl.SCANCODE_Z:
						prevz = zPressed
						zPressed = true
					}
				case sdl.KEYUP:
					switch t.Keysym.Scancode {
					case sdl.SCANCODE_LCTRL, sdl.SCANCODE_RCTRL:
						prevctrl = ctrlPressed
						ctrlPressed = false
					case sdl.SCANCODE_Z:
						prevz = zPressed
						zPressed = false
					}
				}
			}
		}

		if !zPressed && prevctrl && prevz {
			undo()
			fmt.Println("ctrl + z")
			prevctrl = false
			prevz = false
		}

		if app.Update() {
			updateBuffer()
		}
		app.Render()

		tex.Update(nil, unsafe.Pointer(&pixels[0]), 4*int(ScreenWidth))
		renderer.Clear()
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		//sdl.Delay(16)
	}
}

//TODO : rethink coupling
func updateBuffer() {
	if buffer.Size() == 10 {
		buffer.PopFront()
	}
	buffer.PushBack(append([]byte(nil), pixels...))
	fmt.Println("pushed")
}

func undo() {
	if buffer.Size() > 0 {
		pixels, _ = buffer.PopBack()
		fmt.Println("redid")
	}

}
