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

		mapSelection.Shrink(true).Invert().Fill(1) // Outline
		world.GameMap.GenerateBSP(1, 2, 20)
		mapSelection.ByValue(0).ByPercentage(0.1).Fill(3) // 3 is the alternate floor

	} else if WorldGenerationMode == 1 {

		mapSelection.Fill(1)
		world.GameMap.GenerateDrunkWalk(0, 0.5)

	} else {
		mapSelection.Fill(1)
		world.GameMap.GenerateRoomPlacer(0, 6, 3, 3, 5, 5, true)
		mapSelection.ByValue(0).By(func(x, y int) bool {
			return (world.GameMap.Get(x+1, y) == 1 && world.GameMap.Get(x-1, y) == 1) || (world.GameMap.Get(x, y-1) == 1 && world.GameMap.Get(x, y+1) == 1)
		}).ByPercentage(0.25).Fill(2)
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
		world.GameMap.Select().Degrade(0, 1)
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
