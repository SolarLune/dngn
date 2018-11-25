package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"

	"github.com/SolarLune/dngn"
)

type WorldInterface interface {
	Create()
	Update()
	Draw()
	Destroy()
}

var WorldGenerationMode = 0

type World1 struct {
	GameMap *dngn.Room
	Tileset *sdl.Texture
}

func (world *World1) Create() {

	world.Tileset, _ = img.LoadTexture(renderer, "assets/Tileset.png")

	// Base room; this is what we're messing with
	world.GameMap = dngn.NewRoom(40, 23)

	mapSelection := world.GameMap.Select()

	if WorldGenerationMode == 0 {

		mapSelection.RemoveSelection(mapSelection.ByArea(1, 1, world.GameMap.Width-2, world.GameMap.Height-2)).Fill('x')
		world.GameMap.GenerateBSP('x', '#', 20)
		mapSelection.ByRune(' ').ByPercentage(0.1).Fill('.') // '.' is the alternate floor

	} else if WorldGenerationMode == 1 {

		mapSelection.Fill('x')
		world.GameMap.GenerateDrunkWalk(' ', 0.5)

	} else {
		mapSelection.Fill('x')
		world.GameMap.GenerateRoomPlacer(' ', 6, 3, 3, 5, 5, true)

		// This selects the ground tiles that are between walls to place doors randomly. This isn't really good, but it at least
		// gets the idea across.
		mapSelection.ByRune(' ').By(func(x, y int) bool {
			return (world.GameMap.Get(x+1, y) == 'x' && world.GameMap.Get(x-1, y) == 'x') || (world.GameMap.Get(x, y-1) == 'x' && world.GameMap.Get(x, y+1) == 'x')
		}).ByPercentage(0.25).Fill('#')
	}

	fmt.Println(world.GameMap.DataToString())

}

func (world *World1) Update() {
	if keyboard.KeyPressed(sdl.K_1) {
		WorldGenerationMode = 0
		world.Destroy()
		world.Create()
	} else if keyboard.KeyPressed(sdl.K_2) {
		WorldGenerationMode = 1
		world.Destroy()
		world.Create()
	} else if keyboard.KeyPressed(sdl.K_3) {
		WorldGenerationMode = 2
		world.Destroy()
		world.Create()
	}
	if keyboard.KeyPressed(sdl.K_a) {
		world.GameMap.Select().Degrade(1)
	}
}

func (world *World1) Draw() {

	DrawTiles(world.GameMap, world.Tileset)

	if drawHelpText {

		DrawText(0, 0,
			"Press F1 to toggle debug mode drawing",
			"Press 1 to generate a BSP generation room",
			"Press 2 to generate a Drunk Walk room",
			"Press 3 to generate a Room Placement room",
			"Press A to degrade the map.")

	}
}

func (world *World1) Destroy() {

}
