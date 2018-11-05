package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/SolarLune/dngn/dngn"
)

// A pure neighbor-checking test.
type World2 struct {
	GameMap *dngn.Room
}

func (world *World2) Create() {

	world.GameMap = dngn.NewRoom(0, 0, 20, 20, nil)

	dngn.NewRoom(4, 8, 4, 4, world.GameMap)
	dngn.NewRoom(8, 15, 4, 4, world.GameMap)
	dngn.NewRoom(8, 12, 4, 4, world.GameMap)
	dngn.NewRoom(4, 12, 4, 4, world.GameMap)

	// fmt.Println(world.GameMap.Children[1].Neighbors())

	// world.GameMap.GenerateBSP(1, 3)

	// world.GameMap.Fill(1).GenerateDrunkWalk(0, .5)

	world.GameMap.Children[0].Fill(1)
	world.GameMap.Children[1].Fill(2)
	world.GameMap.Children[2].Fill(1)
	world.GameMap.Children[3].Fill(1)
	fmt.Println(world.GameMap.DataToString())

	// world.GameMap.Degrade(1, 0)

}

func (world *World2) Update() {

	change := false

	if keyboard.KeyPressed(sdl.K_UP) {
		world.GameMap.Children[1].Y--
		change = true
	}

	if keyboard.KeyPressed(sdl.K_DOWN) {
		world.GameMap.Children[1].Y++
		change = true
	}

	if keyboard.KeyPressed(sdl.K_LEFT) {
		world.GameMap.Children[1].X--
		change = true
	}

	if keyboard.KeyPressed(sdl.K_RIGHT) {
		world.GameMap.Children[1].X++
		change = true
	}

	if change {
		world.GameMap.Fill(0)
		world.GameMap.Children[1].Fill(2)
		world.GameMap.Children[0].Fill(1)
		world.GameMap.Children[2].Fill(1)
		world.GameMap.Children[3].Fill(1)
		for _, neighbor := range world.GameMap.Children[1].Neighbors() {
			neighbor.Fill(4)
		}
		fmt.Println(world.GameMap.DataToString())
		fmt.Println(world.GameMap.Children[1].Neighbors())
	}

}

func (world *World2) Draw() {}

func (world *World2) Destroy() {}
