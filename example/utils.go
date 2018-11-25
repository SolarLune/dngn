package main

import (
	"github.com/SolarLune/dngn"
	"github.com/veandco/go-sdl2/sdl"
)

func DrawTiles(room *dngn.Room, tileset *sdl.Texture) {

	roomSelect := room.Select()

	for _, cell := range roomSelect.Cells {

		x, y := cell[0], cell[1]
		v := room.Get(x, y)
		left := room.Get(x-1, y) == v
		right := room.Get(x+1, y) == v
		up := room.Get(x, y-1) == v
		down := room.Get(x, y+1) == v

		src := &sdl.Rect{0, 0, 16, 16}
		dst := &sdl.Rect{int32(x) * src.W, int32(y) * src.H, src.W, src.H}
		rotation := 0.0

		if v == ' ' || v == '#' {
			src.Y = src.H
			if room.Get(x, y-1) == 'x' {
				src.Y -= src.H
			}
		}

		if v == '.' {
			if room.Get(x, y-1) == 'x' {
				src.Y = 0
			} else {
				src.Y = 32
			}
		}

		if v == 'x' { // Wall

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

				src.X = src.W * 3
				src.Y = src.H * 1

			} else if num == 1 {

				src.X = src.W * 3
				src.Y = src.H * 2

				if right {
					rotation = 180
				} else if up {
					rotation = 90
				} else if down {
					rotation = -90
				}

			} else if num == 2 {

				if left && right {
					src.X = src.W * 2
					src.Y = src.H * 2
				} else if up && down {
					src.X = src.W * 2
					src.Y = src.H * 2
					rotation = 90
				} else {

					src.X = src.W * 3
					src.Y = src.H * 0

					if up && right {
						rotation = 90
					} else if right && down {
						rotation = 180
					} else if down && left {
						rotation = -90
					}

				}

			} else if num == 3 {
				src.X = src.W * 2

				if up && right && down {
					rotation = 90
				} else if right && down && left {
					rotation = 180
				} else if down && left && up {
					rotation = -90
				}

			} else if num == 4 {
				src.X = src.W * 2
				src.Y = src.H
			}

		}

		renderer.CopyEx(tileset, src, dst, rotation, &sdl.Point{src.W / 2, src.H / 2}, sdl.FLIP_NONE)

	}

	doors := roomSelect.ByRune('#')

	for _, d := range doors.Cells {

		x, y := d[0], d[1]
		src := &sdl.Rect{16, 0, 16, 16}
		dst := &sdl.Rect{int32(x) * src.W, int32(y) * src.H, src.W, src.H}

		if room.Get(x, y-1) != 'x' && room.Get(x, y+1) != 'x' { // Vertical door
			dst.Y += 4
		} else {
			src.Y += src.H
		}

		renderer.CopyEx(tileset, src, dst, 0, &sdl.Point{src.W / 2, src.H / 2}, sdl.FLIP_NONE)

	}

}
