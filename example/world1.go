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

	world.GameMap = dngn.NewRoom(0, 0, 40, 23, nil)

	if WorldGenerationMode == 0 {
		world.GameMap.GenerateBSP(1, 5)
		world.GameMap.OpenIntoNeighbors(0)
	} else {
		world.GameMap.Fill(1)
		world.GameMap.GenerateDrunkWalk(0, .5)
		world.GameMap.OpenIntoNeighbors(0)
	}

	world.GameMap.Outline(1)

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
	} else if keyboard.KeyPressed(sdl.K_a) {
		world.GameMap.Degrade(0, 1)
	}
}

func (world *World1) Draw() {

	DrawTiles(world.GameMap, world.Tileset)

	if drawHelpText {

		DrawText(0, 0,
			"Press F1 to toggle debug mode drawing",
			"Press 1 to generate a BSP room",
			"Press 2 to generate a Drunk Walk room",
			"Press A to degrade the map.")

	}
}

func (world *World1) Destroy() {

}
