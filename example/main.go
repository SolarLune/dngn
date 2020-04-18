package main

import "github.com/hajimehoshi/ebiten"

func main() {

	game := NewGame()
	ebiten.RunGame(game)

}
