package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/inpututil"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/solarlune/dngn"
)

type Game struct {
	Map            *dngn.Layout
	Tileset        *ebiten.Image
	Fontface       font.Face
	GenerationMode int
}

func NewGame() *Game {
	game := &Game{
		Map: dngn.NewLayout(80, 45),
	}
	game.Tileset, _, _ = ebitenutil.NewImageFromFile(game.GetPath("assets", "Tileset.png"), ebiten.FilterNearest)
	game.GenerateMap()

	fontData, _ := os.Open(game.GetPath("assets", "excel.ttf"))
	fontBytes, _ := ioutil.ReadAll(fontData)
	gameFont, _ := truetype.Parse(fontBytes)

	game.Fontface = truetype.NewFace(gameFont, &truetype.Options{
		Size: 10,
	})

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("dngn Example")

	return game
}

func (game *Game) GenerateMap() {

	// The full selection; we use this to tweak the map after generation.
	mapSelection := game.Map.Select()

	if game.GenerationMode == 0 {

		bspOptions := dngn.NewDefaultBSPOptions()
		bspOptions.SplitCount = 60
		bspOptions.MinimumRoomSize = 3

		bspRooms := game.Map.GenerateBSP(bspOptions)

		// Just generating BSP rooms is pretty good, but we can modify it a bit afterwards.
		// Below, we select a room to specify as the "start", and then effectively destroy
		// any rooms that are too far away (at least 5 hops away).

		start := bspRooms[0]

		for _, subroom := range bspRooms {

			subroomCenter := subroom.Center()
			center := game.Map.Center()

			margin := 10

			if subroomCenter.X > center.X-margin && subroomCenter.X < center.X+margin && subroomCenter.Y > center.Y-margin && subroomCenter.Y < center.Y+margin {
				start = subroom
				break
			}

		}

		for _, room := range bspRooms {

			hops := room.CountHopsTo(start)

			if hops < 0 || hops > 4 {
				// We're filtering out a little bit more on the width and height because the walls and doorways in GenerateBSP() are always on the top and left sides of each room.
				// By adding the right and bottom as well, we can nuke any doors that led into rooms we're deleting.
				mapSelection.FilterByArea(room.X, room.Y, room.W+1, room.H+1).Fill('x')
				room.Disconnect()
			}

		}

		// We could also remove unnecessary rooms by simply looping through bspRooms and calling bspRoom.Necessary() - a necessary room is determined to be a room
		// that affords entrance and exit to a room or set of rooms that have no other way to get in or out.

	} else if game.GenerationMode == 1 {

		game.Map.GenerateDrunkWalk(' ', 'x', 0.5)

	} else {
		game.Map.GenerateRandomRooms(' ', 'x', 6, 3, 3, 5, 5, true)

		// This selects the ground tiles that are between walls to place doors randomly. This isn't really good, but it at least
		// gets the idea across.
		mapSelection.FilterByRune(' ').FilterBy(func(x, y int) bool {
			return (game.Map.Get(x+1, y) == 'x' && game.Map.Get(x-1, y) == 'x') || (game.Map.Get(x, y-1) == 'x' && game.Map.Get(x, y+1) == 'x')
		}).FilterByPercentage(0.25).Fill('#')
	}

	// Fill the outer walls
	mapSelection.Remove(mapSelection.FilterByArea(1, 1, game.Map.Width-2, game.Map.Height-2)).Fill('x')

	// Add a different tile for an alternate floor
	mapSelection.FilterByRune(' ').FilterByPercentage(0.1).Fill('.')

	fmt.Println(game.Map.DataToString())

}

func (game *Game) Update(screen *ebiten.Image) error {

	var exit error

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		exit = errors.New("Quit")
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		game.GenerationMode = 0
		game.GenerateMap()
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		game.GenerationMode = 1
		game.GenerateMap()
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		game.GenerationMode = 2
		game.GenerateMap()
	}

	game.DrawTiles(screen)

	return exit

}

func (game *Game) Layout(width, height int) (int, int) {
	return 1280, 720
}

func (game *Game) GetPath(folders ...string) string {

	// Running apps from Finder in MacOS makes the working directory the home directory, which is nice, because
	// now I have to make this function to do what should be done anyway and give me a relative path starting from
	// the executable so that I can load assets from the assets directory. :,)

	cwd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)

	if strings.Contains(exeDir, "go-build") {
		// If the executable's directory contains "go-build", it's probably the result of a "go run" command, so just go with the CWD
		// as the "root" to base the path from
		exeDir = cwd
	}

	return filepath.Join(exeDir, filepath.Join(folders...))

}

func (game *Game) DrawTiles(screen *ebiten.Image) {

	roomSelect := game.Map.Select()

	for cell := range roomSelect.Cells {

		v := game.Map.Get(cell.X, cell.Y)

		left := game.Map.Get(cell.X-1, cell.Y) == v
		right := game.Map.Get(cell.X+1, cell.Y) == v
		up := game.Map.Get(cell.X, cell.Y-1) == v
		down := game.Map.Get(cell.X, cell.Y+1) == v
		rotation := 0.0

		// Tile graphic defaults to plain ground
		srcOffsetX := 0
		srcOffsetY := 16

		if v == ' ' || v == '#' {
			if game.Map.Get(cell.X, cell.Y-1) == 'x' {
				srcOffsetY = 0
			}
		}

		// Minor scratches on the ground
		if v == '.' {
			if game.Map.Get(cell.X, cell.Y-1) == 'x' {
				srcOffsetY = 0
			} else {
				srcOffsetY = 32
			}
		}

		// Wall
		if v == 'x' {

			num := 0
			if left {
				num++
			}
			if right {
				num++
			}
			if up {
				num++
			}
			if down {
				num++
			}

			if num == 0 {

				srcOffsetX = 48
				srcOffsetY = 16

			} else if num == 1 {

				srcOffsetX = 48
				srcOffsetY = 32

				if right {
					rotation = math.Pi
				} else if up {
					rotation = math.Pi / 2
				} else if down {
					rotation = -math.Pi / 2
				}

			} else if num == 2 {

				if left && right {
					srcOffsetX = 32
					srcOffsetY = 32
				} else if up && down {
					srcOffsetX = 32
					srcOffsetY = 32
					rotation = math.Pi / 2
				} else {

					srcOffsetX = 48
					srcOffsetY = 0

					if up && right {
						rotation = math.Pi / 2
					} else if right && down {
						rotation = math.Pi
					} else if down && left {
						rotation = -math.Pi / 2
					}

				}

			} else if num == 3 {
				srcOffsetX = 32
				srcOffsetY = 0

				if up && right && down {
					rotation = math.Pi / 2
				} else if right && down && left {
					rotation = math.Pi
				} else if down && left && up {
					rotation = -math.Pi / 2
				}

			} else if num == 4 {
				srcOffsetX = 32
				srcOffsetY = 16
			}

		}

		src := image.Rect(0, 0, 16, 16)
		src = src.Add(image.Point{srcOffsetX, srcOffsetY})

		tile := game.Tileset.SubImage(src).(*ebiten.Image)
		geoM := ebiten.GeoM{}
		geoM.Translate(-float64(src.Dx()/2), -float64(src.Dy()/2))

		geoM.Rotate(rotation)

		geoM.Translate(float64(src.Dx()/2), float64(src.Dy()/2))

		geoM.Translate(float64(cell.X*src.Dx()), float64(cell.Y*src.Dy()))
		screen.DrawImage(tile, &ebiten.DrawImageOptions{GeoM: geoM})

	}

	doors := roomSelect.FilterByRune('#')

	for d := range doors.Cells {

		src := image.Rect(16, 0, 32, 16)

		dstX, dstY := float64(d.X*src.Dx()), float64(d.Y*src.Dy())

		if game.Map.Get(d.X-1, d.Y) != 'x' && game.Map.Get(d.X+1, d.Y) != 'x' { // Horizontal door
			src = src.Add(image.Point{0, 16})
		} else {
			dstY += 4
		}

		tile := game.Tileset.SubImage(src).(*ebiten.Image)
		geoM := ebiten.GeoM{}
		geoM.Translate(dstX, dstY)
		screen.DrawImage(tile, &ebiten.DrawImageOptions{GeoM: geoM})

	}

	// screen.Fill(color.RGBA{255, 0, 0, 255})

	text.Draw(screen,
		"Press 1 to generate a BSP generation room.\n"+
			"Press 2 to generate a Drunk Walk room.\n"+
			"Press 3 to generate a Room Placement room.\n",
		game.Fontface, 41, 33, color.Black)

	text.Draw(screen,
		"Press 1 to generate a BSP generation room.\n"+
			"Press 2 to generate a Drunk Walk room.\n"+
			"Press 3 to generate a Room Placement room.\n",
		game.Fontface, 40, 32, color.White)

}

func main() {

	game := NewGame()
	ebiten.RunGame(game)

}
