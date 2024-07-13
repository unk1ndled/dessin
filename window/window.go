package window

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"unsafe"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/harry1453/go-common-file-dialog/cfd"
	"github.com/harry1453/go-common-file-dialog/cfdutil"
)

var (
	ScreenWidth  int32
	ScreenHeight int32
	pixels       []byte

	tex *sdl.Texture
)

func GetPixelsIndex(x, y int32) int32 {
	return (x + ((ScreenWidth) * y)) * 4
}
func GetPixelColor(x, y int32) Color {
	index := GetPixelsIndex(x, y)
	return Color{pixels[index], pixels[index+1], pixels[index+2]}
}

func SetPixel(x, y int32, c *Color) {

	index := GetPixelsIndex(x, y)
	if index+3 <= int32(len(pixels))-1 && index >= 0 {
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

	// Load the icon image
	icon, err := img.Load("icon.jpg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load icon: %s\n", err)
		os.Exit(3)
	}
	defer icon.Free()

	// Set the window icon
	window.SetIcon(icon)

	tex, err = renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, ScreenWidth, ScreenHeight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %s\n", err)
		os.Exit(5)
	}
	defer tex.Destroy()

	pixels = make([]byte, ScreenHeight*ScreenWidth*4)

	app.Init(pixels)
	app.Render()
	tex.Update(nil, unsafe.Pointer(&pixels[0]), 4*int(ScreenWidth))
	renderer.Clear()
	renderer.Copy(tex, nil, nil)
	renderer.Present()

	quit := false
	update := false
	for !quit {
		update, quit = app.Update()
		if update {
			app.Render()
			tex.Update(nil, unsafe.Pointer(&pixels[0]), 4*int(ScreenWidth))
			renderer.Clear()
			renderer.Copy(tex, nil, nil)
			renderer.Present()
			sdl.Delay(5)
			continue
		}
		renderer.Present()
		// apparantly the delay really deceases cpu usage
		sdl.Delay(16)
	}
}

func SaveTextureAsImage(x, y, w, h int32, filename string) {
	//TODO : implement text input

	SavePNG(x, y, w, h, ScreenWidth*4)

}

func SavePNG(xOffset, yOffset, w, h, pitch int32) error {

	file, err := cfdutil.ShowSaveFileDialog(cfd.DialogConfig{
		DefaultExtension: "png",
	})
	if err == cfd.ErrorCancelled {
		log.Println("Dialog was cancelled by the user.")
		return nil
	} else if err != nil {
		log.Fatal(err)
	}

	imag := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	for y := int32(0); y < h; y++ {
		for x := int32(0); x < w; x++ {
			index := (y+yOffset)*pitch + (x+xOffset)*4
			r := pixels[index]
			g := pixels[index+1]
			b := pixels[index+2]
			imag.SetRGBA(int(x), int(y), color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, imag); err != nil {
		return fmt.Errorf("could not encode PNG: %v", err)
	}
	return nil
}

func OpenPNG(xOffset, yOffset, canvasW, canvasH int32) error {

	filename, err := cfdutil.ShowOpenFileDialog(cfd.DialogConfig{
		Title: "Open A File",
		Role:  "OpenFileExample",
		FileFilters: []cfd.FileFilter{
			{
				DisplayName: "Text Files (*.txt)",
				Pattern:     "*.txt",
			},
			{
				DisplayName: "Image Files (*.jpg, *.png)",
				Pattern:     "*.jpg;*.png",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
		SelectedFileFilterIndex: 2,
		FileName:                "file.txt",
		DefaultExtension:        "txt",
	})
	if err == cfd.ErrorCancelled {
		log.Println("Dialog was cancelled by the user.")
		return nil
	} else if err != nil {
		log.Fatal(err)
	}
	log.Printf("Chosen file: %s\n", filename)

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}

	img, tp, err := image.Decode(file)
	log.Println(tp)
	if err != nil {
		return fmt.Errorf("could not decode file: %v", err)
	}

	for y := int32(0); y < (canvasH); y++ {
		for x := int32(0); x < (canvasW); x++ {

			r, g, b, _ := img.At(int(x), int(y)).RGBA()
			SetPixel(x+xOffset, y+yOffset, &Color{byte(r), byte(g), byte(b)})
		}

	}
	return nil
}
