package main

import (
	"github.com/unk1ndled/draw/dessin"
	"github.com/unk1ndled/draw/window"
)

func main() {
	window.Visualise("dessin", 700, 540, dessin.NewPaint())
}
