package main

import (
	"strconv"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var screenWidth int32 = 640
var screenHeight int32 = 360

var renderer *sdl.Renderer
var window *sdl.Window
var avgFramerate int

var drawHelpText = true

func main() {

	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	ttf.Init()
	defer ttf.Quit()

	window, renderer, _ = sdl.CreateWindowAndRenderer(screenWidth, screenHeight, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)

	window.SetResizable(true)

	renderer.SetLogicalSize(screenWidth, screenHeight)

	fpsMan := &gfx.FPSmanager{}
	gfx.InitFramerate(fpsMan)
	gfx.SetFramerate(fpsMan, 60)

	// Change this to one of the other World structs to change the world and see different tests

	var world WorldInterface = &World1{}

	world.Create()

	running := true

	var frameCount int
	var framerateDelay uint32

	for running {

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				keyboard.ReportEvent(event.(*sdl.KeyboardEvent))
			}
		}

		keyboard.Update()

		if keyboard.KeyPressed(sdl.K_F2) {
			gfx.SetFramerate(fpsMan, 10)
		} else if keyboard.KeyPressed(sdl.K_F3) {
			gfx.SetFramerate(fpsMan, 60)
		}

		if keyboard.KeyPressed(sdl.K_F1) {
			drawHelpText = !drawHelpText
		}

		if keyboard.KeyPressed(sdl.K_ESCAPE) {
			running = false
		}

		world.Update()

		renderer.SetDrawColor(20, 30, 40, 255)

		renderer.Clear()

		world.Draw()

		framerateDelay += gfx.FramerateDelay(fpsMan)

		if framerateDelay >= 1000 {
			avgFramerate = frameCount
			framerateDelay -= 1000
			frameCount = 0
			// fmt.Println(avgFramerate, " FPS")
		}

		frameCount++

		DrawText(screenWidth-32, 0, strconv.Itoa(avgFramerate))

		renderer.Present()

	}

}

func DrawText(x, y int32, textLines ...string) {

	sy := y

	for _, text := range textLines {

		font, _ := ttf.OpenFont("assets/ARCADEPI.TTF", 12)
		defer font.Close()

		var surf *sdl.Surface

		surf, _ = font.RenderUTF8Solid(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})

		textSurface, _ := renderer.CreateTextureFromSurface(surf)
		defer textSurface.Destroy()

		_, _, w, h, _ := textSurface.Query()

		textSurface.SetAlphaMod(100)
		renderer.Copy(textSurface, &sdl.Rect{X: 0, Y: 0, W: w, H: h}, &sdl.Rect{X: x, Y: sy, W: w, H: h})

		sy += 16

	}

}
