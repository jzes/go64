package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/drivers/display"
	"github.com/clktmr/n64/drivers/draw"
	"github.com/clktmr/n64/fonts/gomono12"
	_ "github.com/clktmr/n64/machine"
	"github.com/clktmr/n64/rcp/video"
)

var face = gomono12.NewFace()
var background = &image.Uniform{color.RGBA{0x7f, 0x7f, 0xaf, 0x0}}

func main() {
	// Enable video output
	video.Setup(false)

	// Allocate framebuffer
	display := display.NewDisplay(image.Pt(320, 240), video.BPP16)

	controlsChan := make(chan [4]controller.Controller)
	go func() {
		var states [4]controller.Controller
		for {
			controller.Poll(&states)
			controlsChan <- states
		}
	}()

	for {
		fb := display.Swap() // Blocks until next VBlank

		textarea := fb.Bounds().Inset(15)
		pt := textarea.Min.Add(image.Pt(0, int(face.Ascent)))

		draw.Src.Draw(fb, fb.Bounds(), background, fb.Bounds().Min)

		text := fmt.Appendln(nil, "Vanessinha é muito lindona !s2")
		input := <-controlsChan
		text = fmt.Appendln(text, input[0].Down())
		pt = draw.DrawText(fb, textarea, face, pt, image.Black, nil, text)

		draw.Flush() // Blocks until everything is drawn
	}
}
