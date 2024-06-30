package main

import (
	"github.com/unk1ndled/draw/dessin"
	"github.com/unk1ndled/draw/window"
)

func main() {
	window.Visualise("dessin", 400, 400, dessin.NewCanvas())
}
