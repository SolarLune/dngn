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

	"github.com/SolarLune/dngn"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	Room           *dngn.Room
	Tileset        *ebiten.Image
	Fontface       font.Face
	GenerationMode int
}

func NewGame() *Game {
	game := &Game{
		Room: dngn.NewRoom(40, 23),
	}
	game.Tileset, _, _ = ebitenutil.NewImageFromFile(game.GetPath("example", "assets", "Tileset.png"), ebiten.FilterNearest)
	game.GenerateRoom()

	fontData, _ := os.Open(game.GetPath("example", "assets", "excel.ttf"))
	fontBytes, _ := ioutil.ReadAll(fontData)
	gameFont, _ := truetype.Parse(fontBytes)

	game.Fontface = truetype.NewFace(gameFont, &truetype.Options{
		Size: 10,
	})

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("dngn Example")

	return game
}

func (game *Game) GenerateRoom() {

	game.Room.Select().Fill(' ')

	mapSelection := game.Room.Select()

	switch game.GenerationMode {
	case 0:
		mapSelection.RemoveSelection(mapSelection.ByArea(1, 1, game.Room.Width-2, game.Room.Height-2)).Fill('x')
		game.Room.GenerateBSP('x', '#', 20)
		mapSelection.ByRune(' ').ByPercentage(0.1).Fill('.') // '.' is the alternate floor
	case 1:
		mapSelection.Fill('x')
		game.Room.GenerateDrunkWalk(' ', 0.5)
	case 2:
		mapSelection.Fill('x')
		game.Room.GenerateRandomRooms('R', 'H', 6, 3, 3, 5, 5, true, false)

		// A hallway space with three room spaces next to it as well as two walls that are non-diagonal
		// Desired door placement
		//      R
		// HHHHDR
		//      R
		doorLocations := game.Room.Select().ByRune('H').ByNeighbor('R', 3, true).By(func(x, y int) bool {
			return (game.Room.Get(x+1, y) == 'x' && game.Room.Get(x-1, y) == 'x') || (game.Room.Get(x, y+1) == 'x' && game.Room.Get(x, y-1) == 'x')
		})

		doorLocations.ByPercentage(.5).Fill('#')
		game.Room.Select().ByRune('H').Fill(' ')
		game.Room.Select().ByRune('R').Fill(' ')

	default:
		mapSelection.Fill('x')
		game.Room.GenerateRandomRooms(' ', ' ', 6, 3, 3, 5, 5, true, true)

		// This selects the ground tiles that are between walls to place doors randomly. This isn't really good, but it at least
		// gets the idea across.
		mapSelection.ByRune(' ').By(func(x, y int) bool {
			return (game.Room.Get(x+1, y) == 'x' && game.Room.Get(x-1, y) == 'x') || (game.Room.Get(x, y-1) == 'x' && game.Room.Get(x, y+1) == 'x')
		}).ByPercentage(0.25).Fill('#')
	}
	fmt.Println(game.Room.DataToString())
}

func (game *Game) Update(screen *ebiten.Image) error {

	var exit error

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		exit = errors.New("Quit")
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		game.GenerationMode = 0
		game.GenerateRoom()
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		game.GenerationMode = 1
		game.GenerateRoom()
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		game.GenerationMode = 3
		game.GenerateRoom()
	} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
		game.GenerationMode = 2
		game.GenerateRoom()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		game.Room.Select().Degrade('x')
	}

	game.DrawTiles(screen)

	return exit

}

func (game *Game) Layout(width, height int) (int, int) {
	return 640, 360
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

	roomSelect := game.Room.Select()

	for _, cell := range roomSelect.Cells {

		x, y := cell[0], cell[1]

		v := game.Room.Get(x, y)
		left := game.Room.Get(x-1, y) == v
		right := game.Room.Get(x+1, y) == v
		up := game.Room.Get(x, y-1) == v
		down := game.Room.Get(x, y+1) == v
		rotation := 0.0

		// Tile graphic defaults to plain ground
		srcOffsetX := 0
		srcOffsetY := 16

		if v == ' ' || v == '#' {
			if game.Room.Get(x, y-1) == 'x' {
				srcOffsetY = 0
			}
		}

		// Minor scratches on the ground
		if v == '.' {
			if game.Room.Get(x, y-1) == 'x' {
				srcOffsetY = 0
			} else {
				srcOffsetY = 16
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

		geoM.Translate(float64(x*src.Dx()), float64(y*src.Dy()))
		screen.DrawImage(tile, &ebiten.DrawImageOptions{GeoM: geoM})

	}

	doors := roomSelect.ByRune('#')

	for _, d := range doors.Cells {

		x, y := d[0], d[1]

		src := image.Rect(16, 0, 32, 16)

		dstX, dstY := float64(x*src.Dx()), float64(y*src.Dy())

		if game.Room.Get(x-1, y) != 'x' && game.Room.Get(x+1, y) != 'x' { // Horizontal door
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
			"Press 3 to generate a Room Placement room.\n"+
			"Press A to degrade the map.",
		game.Fontface, 41, 33, color.Black)

	text.Draw(screen,
		"Press 1 to generate a BSP generation room.\n"+
			"Press 2 to generate a Drunk Walk room.\n"+
			"Press 3 to generate a Room Placement room with diagonals.\n"+
			"Press 4 to generate a Room Placement room without diagonals\n"+
			"Press A to degrade the map.",
		game.Fontface, 40, 32, color.White)

}
