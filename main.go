package main

import (
	"embed"
	"fmt"
	"image"
	"image/color"

	"github.com/clktmr/n64/drivers/cartfs"
	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/drivers/display"
	"github.com/clktmr/n64/drivers/draw"
	"github.com/clktmr/n64/fonts/gomono12"
	_ "github.com/clktmr/n64/machine"
	"github.com/clktmr/n64/rcp/serial/joybus"
	"github.com/clktmr/n64/rcp/texture"

	"github.com/clktmr/n64/rcp/video"
)

var (
	//go:embed gopher-anim.CI8
	_storageFiles embed.FS
	storageFiles  cartfs.FS = cartfs.Embed(_storageFiles)
)

type Controllers [4]controller.Controller

const (
	blowsToWin          = 8
	amountOfControllers = 4
	playerOneIndex      = 0
	banner              = "Sopre a fita pro gopher poder jogar!\n"
	spriteSize          = 128
	spriteSheetFile     = "gopher-anim.CI8"
)

var face = gomono12.NewFace()
var background = &image.Uniform{color.RGBA{0x7f, 0x7f, 0xaf, 0x0}}

func main() {
	// Setup video output
	video.Setup(false)

	// Allocate framebuffer
	display := display.NewDisplay(image.Pt(320, 240), video.BPP16)

	// start controller polling
	ctrlsChan := make(chan Controllers)
	go ReadControllers(ctrlsChan)

	// Load cartridge files
	gopherFile, err := storageFiles.Open(spriteSheetFile)
	if err != nil {
		panic(err)
	}

	gopherTexture, err := texture.Load(gopherFile)
	if err != nil {
		panic(err)
	}

	gopherRectangle := image.Rect(0, 0, spriteSize, spriteSize)
	blows := 0

	for {
		fb := display.Swap() // Blocks until next VBlank

		draw.Src.Draw(fb, fb.Bounds(), background, fb.Bounds().Min)

		ctrlsInput := <-ctrlsChan
		player1Ctrl := ctrlsInput[playerOneIndex]

		textarea := fb.Bounds().Inset(15)
		pt := textarea.Min.Add(image.Pt(0, int(face.Ascent)))

		bannerText := GetBanner(blows)
		pt = draw.DrawText(fb, textarea, face, pt, image.Black, nil, bannerText)

		blows = HandleBlowInput(player1Ctrl, blows)
		gopherFrame := SetFrame(player1Ctrl, blows)

		draw.Over.Draw(fb, gopherRectangle.Add(pt), gopherTexture, gopherFrame)
		draw.Flush() // Blocks until everything is drawn
	}
}

func ReadControllers(ctrlsChan chan Controllers) {
	var states [4]controller.Controller
	for {
		controller.Poll(&states)
		ctrlsChan <- states
	}
}

func HandleBlowInput(ctrl controller.Controller, blows int) int {
	if blows < blowsToWin && IsButtonPressed(ctrl, joybus.ButtonA) {
		blows++
	}
	return blows
}

func SetFrame(ctrl controller.Controller, blows int) image.Point {
	frame := image.Point{} // default gopher
	if blows < blowsToWin {
		if IsButtonDown(ctrl, joybus.ButtonA) {
			frame.X += spriteSize // blowing gopher
		}
	} else {
		frame.X += spriteSize * 2 // happy gopher
	}
	return frame
}

func IsButtonPressed(ctrl controller.Controller, button joybus.ButtonMask) bool {
	return ctrl.Pressed()&button != 0
}

func IsButtonDown(ctrl controller.Controller, button joybus.ButtonMask) bool {
	return ctrl.Down()&button != 0
}

func GetBanner(blows int) []byte {
	fixBannerBytes := []byte(banner)
	switch blows {
	case 0:
		return fmt.Appendf(fixBannerBytes, "Sopradas.: %d, Vamos lá!\n", blows)
	case 1, 2, 3:
		return fmt.Appendf(fixBannerBytes, "Sopradas.: %d, Continue Assim!\n", blows)
	case 4, 5:
		return fmt.Appendf(fixBannerBytes, "Sopradas.: %d, Quase lá!\n", blows)
	case 6, 7:
		return fmt.Appendf(fixBannerBytes, "Sopradas.: %d, Ultimo Esforço\n", blows)
	case 8:
		return fmt.Appendf(fixBannerBytes, "Sopradas.: %d, Aeeee Vamos Jogar\n", blows)
	default:
		return fixBannerBytes
	}
}
