package main

import (
	"github.com/unk1ndled/draw/dessin"
	"github.com/unk1ndled/draw/window"
)

func main() {
	window.Visualise("dessin", 700, 540, dessin.NewPaint())
}

// type fum struct{}

// func (f *fum) fourtimes(n int) int {
// 	return (4 * n)
// }
// func foo(fn func(int) int) {
// 	fn(1)
// }
// func main() {
// 	f := fum{}
// 	foo(f.fourtimes)
// }
