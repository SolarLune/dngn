/*
Package dngn is a simple random map generation library primarily made to be used for 2D games. It features a simple API,
and a couple of different means to generate maps. The easiest way to kick things off when using dngn is to simply create a Room
to represent your overall game map, which can then be manipulated or have a Generate function run on it to actually generate the
content on the map.
*/
package dngn

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Room represents a dungeon map.
// Width and Height are the width and height of the Room in the layout. This determines the size of the overall Data structure
// backing the Room layout.
// Data is the core underlying data structure representing the dungeon. It's a 2D array of runes.
// Seed is the seed of the Room to use when doing random generation using the Generate* functions below.
// CustomSeed indicates whether the Seed was customized - if not, then it will default to using the time of the system to have
// random generation each time you use Generate* functions.
type Room struct {
	Width, Height int
	Data          [][]rune
	seed          int64
	CustomSeed    bool
}

// NewRoom returns a new Room with the specified width and height.
func NewRoom(width, height int) *Room {

	r := &Room{Width: width, Height: height}
	r.Data = [][]rune{}
	for y := 0; y < height; y++ {
		r.Data = append(r.Data, []rune{})
		for x := 0; x < width; x++ {
			r.Data[y] = append(r.Data[y], ' ')
		}
	}

	return r

}

// NewRoomFromRuneArrays creates a new Room with the data contained in the provided rune arrays.
func NewRoomFromRuneArrays(arrays [][]rune) *Room {

	r := &Room{Width: len(arrays[0]), Height: len(arrays)}
	r.Data = [][]rune{}
	for y := 0; y < len(arrays); y++ {
		r.Data = append(r.Data, []rune{})
		for x := 0; x < len(arrays[0]); x++ {
			r.Data[y] = append(r.Data[y], arrays[y][x])
		}
	}

	return r

}

// NewRoomFromStringArray creates a new Room with the data contained in the provided string array.
func NewRoomFromStringArray(array []string) *Room {

	r := &Room{Width: len(array[0]), Height: len(array)}
	r.Data = [][]rune{}
	for y := 0; y < len(array); y++ {
		asRunes := []rune(array[y])
		r.Data = append(r.Data, []rune{})
		for x := 0; x < len(array[0]); x++ {
			r.Data[y] = append(r.Data[y], asRunes[x])
		}
	}

	return r

}

// GenerateBSP generates a map using BSP (binary space partitioning) generation, drawing lines of wallValue runes horizontally and
// vertically across, partitioning the room into pieces. It also will place single cells of doorValue on the walls, creating
// doorways. Link: http://www.roguebasin.com/index.php?title=Basic_BSP_Dungeon_generation
// BUG: GenerateBSP doesn't handle pre-existing doorValues well (i.e. if 0's already exist on the map and you try to use a
// doorValue of 0 to indicate not to place doors, it bugs out; either it might place a wall where a doorway exists, or it
// won't place anything at all because everywhere could be a door). A workaround is to use another value that you turn into
// 0's later.
func (room *Room) GenerateBSP(wallValue, doorValue rune, numSplits int) {

	type subroom struct {
		X, Y, W, H int
	}

	subMinSize := func(subroom subroom) int {
		if subroom.W < subroom.H {
			return subroom.W
		}
		return subroom.H
	}

	subSplit := func(parent subroom, vertical bool) (subroom, subroom, bool) {

		splitPercentage := 0.2 + rand.Float32()*0.6

		if vertical {

			splitCX := int(float32(parent.W) * splitPercentage)
			splitCX2 := parent.W - splitCX
			a, b := subroom{parent.X, parent.Y, splitCX, parent.H}, subroom{parent.X + splitCX, parent.Y, splitCX2, parent.H}

			if subMinSize(a) <= 2 || subMinSize(b) <= 2 {
				return a, b, false
			}

			// Line is attempting to start on a door
			if doorValue != wallValue && doorValue != 0 && (room.Get(parent.X+splitCX, parent.Y) == doorValue || room.Get(parent.X+splitCX, parent.Y+parent.H) == doorValue) {
				return a, b, false
			}

			room.DrawLine(parent.X+splitCX, parent.Y+1, parent.X+splitCX, parent.Y+parent.H-1, wallValue, 1, false)

			// Place door
			for i := 0; i < 100; i++ {
				ry := parent.Y + 1 + rand.Intn(parent.H-1)
				if room.Get(parent.X+splitCX-1, ry) == wallValue || room.Get(parent.X+splitCX+1, ry) == wallValue {
					continue
				}
				room.Set(parent.X+splitCX, ry, doorValue)
				break
			}

			return a, b, true
		}

		splitCY := int(float32(parent.H) * splitPercentage)
		splitCY2 := parent.H - splitCY
		a, b := subroom{parent.X, parent.Y, parent.W, splitCY}, subroom{parent.X, parent.Y + splitCY, parent.W, splitCY2}

		if subMinSize(a) <= 2 || subMinSize(b) <= 2 {
			return a, b, false
		}

		// Line is attempting to start on a door
		if doorValue != wallValue && doorValue != 0 && (room.Get(parent.X, parent.Y+splitCY) == doorValue || room.Get(parent.X+parent.W, parent.Y+splitCY) == doorValue) {
			return a, b, false
		}

		room.DrawLine(parent.X+1, parent.Y+splitCY, parent.X+parent.W-1, parent.Y+splitCY, wallValue, 1, false)

		// Create doors somewhere in the lines
		for i := 0; i < 100; i++ {
			rx := parent.X + 1 + rand.Intn(parent.W-1)
			if room.Get(rx, parent.Y+splitCY-1) == wallValue || room.Get(rx, parent.Y+splitCY+1) == wallValue {
				continue
			}
			room.Set(rx, parent.Y+splitCY, doorValue)
			break
		}

		return a, b, true

	}

	rooms := []subroom{subroom{0, 0, room.Width, room.Height}}

	if room.CustomSeed {
		rand.Seed(room.seed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	splitCount := 0

	i := 0
	for true {

		splitChoice := rooms[rand.Intn(len(rooms))]

		// Do the split

		a, b, success := subSplit(splitChoice, randBool())

		i++

		if i >= numSplits*10 { // Assume it's done to avoid just hanging the system
			break
		}

		if !success {
			continue
		} else {

			rooms = append(rooms, a, b)

			for i, r := range rooms {
				if r == splitChoice {
					rooms = append(rooms[:i], rooms[i+1:]...)
					break
				}
			}

		}

		splitCount++

		if splitCount >= numSplits {
			break
		}

	}

}

// GenerateRandomRooms generates a map using random room creation.
// roomFillRune is the rune to fill the rooms generated with.
// hallwayFillRune is the rune to fill the hallways with
// roomCount is how many rooms to place
// roomMinWidth and Height are how small they can be, minimum, while roomMaxWidth and Height are how large
//
// connectRooms determines if the algorithm should also attempt to connect the rooms using pathways between each room.
// allowDiagonal determines the shape of the connections. If true, then a direct connection will be made. If false, then only horizontal
// or vertical connections are made.
//
// The function returns the positions of each room created.
func (room *Room) GenerateRandomRooms(roomFillRune rune, hallwayFillRune rune, roomCount, roomMinWidth, roomMinHeight, roomMaxWidth, roomMaxHeight int, connectRooms bool, allowDiagonal bool) [][]int {

	if room.CustomSeed {
		rand.Seed(room.seed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	roomPositions := make([][]int, 0)

	for i := 0; i < roomCount; i++ {
		sx := rand.Intn(room.Width)
		sy := rand.Intn(room.Height)

		roomPositions = append(roomPositions, []int{sx, sy})

		roomW := roomMinWidth + rand.Intn(roomMaxWidth-roomMinWidth)
		roomH := roomMinHeight + rand.Intn(roomMaxHeight-roomMinHeight)

		drawRoom := func(x, y int) bool {
			dx := int(math.Abs(float64(sx) - float64(x)))
			dy := int(math.Abs(float64(sy) - float64(y)))
			if dx < roomW && dy < roomH {
				room.Set(x, y, roomFillRune)
			}
			return true
		}

		room.Select().By(drawRoom)

	}

	// Rooms are drawn. Save them for later use after connections are made
	roomMap := make(map[int]Selection)
	for idx, r := range roomPositions {
		roomMap[idx] = room.SelectContiguous(r[0], r[1])
	}

	if connectRooms && allowDiagonal {
		for p := 0; p < len(roomPositions); p++ {

			if p < len(roomPositions)-1 {

				x := roomPositions[p][0]
				y := roomPositions[p][1]

				x2 := roomPositions[p+1][0]
				y2 := roomPositions[p+1][1]

				room.DrawLine(x, y, x2, y2, hallwayFillRune, 1, true)
			}
		}
	}

	if connectRooms && !allowDiagonal {
		for a := range roomPositions {
			var (
				x  int
				y  int
				x2 int
				y2 int
			)
			for b := range roomPositions {
				if a == b {
					break
				}
				var connected bool = false
				for _, room1Coord := range roomMap[a].Cells {
					for _, room2Coord := range roomMap[b].Cells {
						if room1Coord[0] == room2Coord[0] || room1Coord[1] == room2Coord[1] {
							x = room1Coord[0]
							y = room1Coord[1]
							x2 = room2Coord[0]
							y2 = room2Coord[1]
							room.DrawLine(x, y, x2, y2, hallwayFillRune, 1, false)
							connected = true
							break
						}
					}
					if connected {
						break
					}
				}
			}
		}
	}
	// Fix up the rooms to clear out errant hallways
	for _, room := range roomMap {
		room.Fill(roomFillRune)
	}
	return roomPositions

}

// GenerateDrunkWalk generates a map in the bounds of the Room specified using drunk walking. It will pick a random point in the
// Room and begin walking around at random, placing fillRune in the Room, until at least percentageFilled (0.0 - 1.0) of the Room
// is filled. Note that it only counts values placed in the cell, not instances where it moves over a cell that already has the
// value being placed.
// Link: http://www.roguebasin.com/index.php?title=Random_Walk_Cave_Generation
func (room *Room) GenerateDrunkWalk(fillRune rune, percentageFilled float32) {

	if room.CustomSeed {
		rand.Seed(room.seed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	sx := rand.Intn(room.Width)
	sy := rand.Intn(room.Height)
	fillCount := float32(0)

	totalArea := float32(room.Area())

	for true {

		cell := room.Get(sx, sy)

		if cell != fillRune {
			room.Set(sx, sy, fillRune)
			fillCount++
		}

		dir := rand.Intn(4)

		if dir == 0 {
			sx++
		} else if dir == 1 {
			sx--
		} else if dir == 2 {
			sy++
		} else if dir == 3 {
			sy--
		}

		if sx < 0 {
			sx = 0
		} else if sx >= room.Width {
			sx = room.Width - 1
		}

		if sy < 0 {
			sy = 0
		} else if sy >= room.Height {
			sy = room.Height - 1
		}

		if fillCount/totalArea >= percentageFilled {
			break
		}

	}

}

// Rotate rotates the entire room 90 degrees clockwise.
func (room *Room) Rotate() {

	newData := make([][]rune, 0)

	for y := 0; y < len(room.Data[0]); y++ {
		newData = append(newData, []rune{})
		for x := 0; x < len(room.Data); x++ {
			nx := (room.Height - x) - 1
			newData[y] = append(newData[y], room.Data[nx][y])
		}
	}

	room.Data = newData
	room.Height = len(room.Data)
	room.Width = len(room.Data[0])

}

// CopyFrom copies the data from the other Room into this Room's data. x and y are the position of the other Room's data in the
// destination (calling) Room.
func (room *Room) CopyFrom(other *Room, x, y int) {

	for cy := 0; cy < room.Height; cy++ {
		for cx := 0; cx < room.Width; cx++ {
			if cx >= x && cy >= y && cx-x < other.Width && cy-y < other.Height {
				room.Set(cx, cy, other.Get(cx-x, cy-y))
			}
		}
	}

}

// DrawLine is used to draw a line from x, y, to x2, y2, placing the rune specified by fillRune in the cells between those points (including)
// in those points themselves, as well. thickness controls how thick the line is. If stagger is on, then the line will stagger it's
// vertical movement, allowing a 1-thickness line to actually be pass-able if an object was only able to move in cardinal directions
// and the line had a diagonal slope.
func (room *Room) DrawLine(x, y, x2, y2 int, fillRune rune, thickness int, stagger bool) {

	dx := int(math.Abs(float64(x2 - x)))
	dy := int(math.Abs(float64(y2 - y)))
	slope := float32(0)
	xAxis := true

	if dx != 0 {
		slope = float32(dy) / float32(dx)
	}
	length := dx

	if dy > dx {
		xAxis = false
		if dy != 0 {
			slope = float32(dx) / float32(dy)
		}
		length = dy
	}

	sx := float32(x)
	sy := float32(y)

	set := func(x, y int) {
		for fx := 0; fx < thickness; fx++ {
			for fy := 0; fy < thickness; fy++ {
				room.Set(x+fx-thickness/2, y+fy-thickness/2, fillRune)
			}
		}
	}

	for c := 0; c < length+1; c++ {

		set(int(math.Round(float64(sx))), int(math.Round(float64(sy))))

		mx := int(math.Round(float64(sx)))

		if xAxis {
			if x2 > x {
				sx++
			} else {
				sx--
			}
			if y2 > y {
				sy += slope
			} else {
				sy -= slope
			}
		} else {
			if y2 > y {
				sy++
			} else {
				sy--
			}
			if x2 > x {
				sx += slope
			} else {
				sx -= slope
			}
		}

		if stagger {
			set(mx, int(math.Round(float64(sy))))
		}

	}

}

// Set sets the rune provided in the Room's Data. A convenience function stand-in for "room.Data[y][x] = value".
func (room *Room) Set(x, y int, char rune) {

	if x < 0 {
		x = 0
	} else if x > room.Width-1 {
		x = room.Width - 1
	}

	if y < 0 {
		y = 0
	} else if y > room.Height-1 {
		y = room.Height - 1
	}

	room.Data[y][x] = char

}

// Get returns the rune in the specified position in the Room's Data array. A convenience function stand-in for "value := room.Data[y][x]".
func (room *Room) Get(x, y int) rune {

	if x < 0 || x >= room.Width || y < 0 || y >= room.Height {
		return 0
	}

	return room.Data[y][x]
}

// SetSeed sets a custom seed for random generation.
func (room *Room) SetSeed(seed int64) {
	room.CustomSeed = true
	room.seed = seed
}

// ClearSeed clears a custom seed set for random generation. When using a clear seed, random generation functions will use the
// system's Unix time.
func (room *Room) ClearSeed() {
	room.CustomSeed = false
}

// Center returns the center position of the Room.
func (room *Room) Center() (int, int) {
	return room.Width / 2, room.Height / 2
}

// Resize resizes the room to be of the width and height provided. Note that resizing to a smaller Room is destructive (and so,
// data will be lost if resizing to a smaller Room).
func (room *Room) Resize(width, height int) *Room {

	room.Width = width
	room.Height = height

	data := make([][]rune, 0)

	for y := 0; y < height; y++ {

		data = append(data, []rune{})

		for x := 0; x < width; x++ {

			if len(room.Data) > y && len(room.Data[y]) > x {
				data[y] = append(data[y], room.Get(x, y))
			} else {
				data[y] = append(data[y], 0)
			}

		}

	}

	room.Data = data

	return room

}

// Area returns the overall size of the Room by multiplying the width by the height.
func (room *Room) Area() int {
	return room.Width * room.Height
}

// MinimumSize returns the minimum distance (width or height) for the Room.
func (room *Room) MinimumSize() int {
	if room.Width < room.Height {
		return room.Width
	}
	return room.Height
}

// DataToString returns the underlying data of the overall Room layout in an easily understood visual format.
// 0's turn into blank spaces when using DataToString, and the column is shown at the left of the map.
func (room *Room) DataToString() string {

	s := fmt.Sprintf("  W:%d H:%d\n", room.Width, room.Height)

	for y := 0; y < len(room.Data); y++ {
		s += fmt.Sprintf("%3d  |", y)
		for x := 0; x < len(room.Data[y]); x++ {
			// s += " " + string(room.Data[y][x])
			s += fmt.Sprintf("%v ", string(room.Data[y][x]))
		}
		s += "|\n"
	}

	return s

}

// Select generates a Selection containing all of the cells of the Room.
func (room *Room) Select() Selection {

	cells := make([][]int, 0)

	for y := 0; y < room.Height; y++ {
		for x := 0; x < room.Width; x++ {
			cells = append(cells, []int{x, y})
		}
	}

	return Selection{Room: room, Cells: cells}

}

// SelectContiguous generates a Selection containing all contiguous (connected) cells featuring the same value as the cell in
// the position provided.
func (room *Room) SelectContiguous(x, y int) Selection {

	cells := make([][]int, 0)

	toBeChecked := make([][]int, 0)
	toBeChecked = append(toBeChecked, []int{x, y})

	conValue := room.Get(x, y)

	checked := make([][]int, 0)

	hasBeenChecked := func(x, y int) bool {

		for _, pos := range checked {
			if pos[0] == x && pos[1] == y {
				return true
			}
		}
		return false

	}

	for true {

		sx, sy := toBeChecked[0][0], toBeChecked[0][1]

		if !hasBeenChecked(sx, sy) {

			checked = append(checked, []int{sx, sy})
			cells = append(cells, []int{sx, sy})

			if sx > 0 && room.Get(sx-1, sy) == conValue && !hasBeenChecked(sx-1, sy) {
				toBeChecked = append(toBeChecked, []int{sx - 1, sy})
			}
			if sx < room.Width && room.Get(sx+1, sy) == conValue && !hasBeenChecked(sx+1, sy) {
				toBeChecked = append(toBeChecked, []int{sx + 1, sy})
			}
			if sy > 0 && room.Get(sx, sy-1) == conValue && !hasBeenChecked(sx, sy-1) {
				toBeChecked = append(toBeChecked, []int{sx, sy - 1})
			}
			if sy < room.Height && room.Get(sx, sy+1) == conValue && !hasBeenChecked(sx, sy+1) {
				toBeChecked = append(toBeChecked, []int{sx, sy + 1})
			}

		}

		toBeChecked = toBeChecked[1:]

		if len(toBeChecked) == 0 {
			break
		}

	}

	return Selection{Room: room, Cells: cells}

}

// A Selection represents a selection of cell positions in the Room's data array, and can be filtered down and manipulated
// using the functions on the Selection struct.
type Selection struct {
	Room  *Room
	Cells [][]int
}

// ByRune filters the Selection down to the cells that have the character (rune) provided.
func (selection Selection) ByRune(value rune) Selection {

	return selection.By(func(x, y int) bool {
		if selection.Room.Get(x, y) == value {
			return true
		}
		return false
	})

}

// ByPercentage selects the provided percentage of the cells curently in the Selection.
func (selection Selection) ByPercentage(percentage float32) Selection {

	cells := make([][]int, 0)

	for i := 0; i < int(float32(len(selection.Cells))*percentage); i++ {

		c := selection.Cells[rand.Intn(len(selection.Cells))]

		choseCell := false

		for _, o := range cells {
			if c[0] == o[0] && c[1] == o[1] {
				i--
				choseCell = true
			}
		}

		if !choseCell {
			cells = append(cells, c)
		}

	}

	selection.Cells = cells

	return selection

}

// ByArea filters down a selection by only selecting the cells that have X and Y values between X, Y, and X+W and Y+H.
// It crops the selection, basically.
func (selection Selection) ByArea(x, y, w, h int) Selection {

	return selection.By(func(cx, cy int) bool {
		return cx >= x && cy >= y && cx <= x+w-1 && cy <= y+h-1
	})

}

// AddSelection adds the cells in the other Selection to the current one if they're not already in it.
func (selection Selection) AddSelection(other Selection) Selection {

	sel := Selection{Room: selection.Room}
	cells := make([][]int, 0)

	for _, c := range selection.Cells {
		cells = append(cells, c)
	}

	for _, c2 := range other.Cells {

		contained := false

		for _, c1 := range selection.Cells {

			if c1[0] == c2[0] && c1[1] == c2[1] {
				contained = true
				break
			}

		}

		if !contained {
			cells = append(cells, c2)
		}

	}

	sel.Cells = cells

	return sel

}

// RemoveSelection removes the cells in the other Selection from the current one if they are already in it.
func (selection Selection) RemoveSelection(other Selection) Selection {

	sel := Selection{Room: selection.Room}
	cells := make([][]int, 0)

	for _, c1 := range selection.Cells {

		inside := false

		for _, c2 := range other.Cells {

			if c1[0] == c2[0] && c1[1] == c2[1] {
				inside = true
				break
			}

		}

		if !inside {
			cells = append(cells, c1)
		}

	}

	sel.Cells = cells

	return sel

}

// ByNeighbor selects the cells that are surrounded at least by minNeighborCount neighbors with a value of
// neighborValue. If diagonals is true, then diagonals are also checked.
func (selection Selection) ByNeighbor(neighborValue rune, minNeighborCount int, diagonals bool) Selection {

	return selection.By(func(x, y int) bool {

		n := 0

		if selection.Room.Get(x-1, y) == neighborValue {
			n++
		}
		if selection.Room.Get(x+1, y) == neighborValue {
			n++
		}
		if selection.Room.Get(x, y-1) == neighborValue {
			n++
		}
		if selection.Room.Get(x, y+1) == neighborValue {
			n++
		}

		if diagonals {
			if selection.Room.Get(x-1, y-1) == neighborValue {
				n++
			}
			if selection.Room.Get(x+1, y-1) == neighborValue {
				n++
			}
			if selection.Room.Get(x-1, y+1) == neighborValue {
				n++
			}
			if selection.Room.Get(x+1, y+1) == neighborValue {
				n++
			}
		}

		return n >= minNeighborCount

	})

}

// By simply takes a function that takes the X and Y values of each cell position contained in the Selection, and returns a
// boolean to indicate whether to include that cell in the Selection or not. This allows you to easily make custom filtering
// functions to filter down the cells in a Selection.
func (selection Selection) By(filterFunc func(x, y int) bool) Selection {

	cells := make([][]int, 0)

	for _, c := range selection.Cells {
		if filterFunc(c[0], c[1]) {
			cells = append(cells, []int{c[0], c[1]})
		}
	}

	selection.Cells = cells

	return selection

}

// Expand expands the selection outwards by the distance value provided. Diagonal indicates if the expansion should happen
// diagonally as well, or just on the cardinal 4 directions.
func (selection Selection) Expand(distance int, diagonal bool) Selection {

	cells := make([][]int, 0)

	inCells := func(x, y int) bool {
		for _, c := range cells {
			if x == c[0] && y == c[1] {
				return true
			}
		}
		return false
	}

	for _, cell := range selection.Cells {
		x, y := cell[0], cell[1]
		if !inCells(x, y) {
			cells = append(cells, []int{x, y})
		}

		if !inCells(x-distance, y) {
			cells = append(cells, []int{x - distance, y})
		}
		if !inCells(x+distance, y) {
			cells = append(cells, []int{x + distance, y})
		}
		if !inCells(x, y-distance) {
			cells = append(cells, []int{x, y - distance})
		}
		if !inCells(x, y+distance) {
			cells = append(cells, []int{x, y + distance})
		}
		if diagonal {
			if !inCells(x-distance, y-distance) {
				cells = append(cells, []int{x - distance, y - distance})
			}
			if !inCells(x+distance, y+distance) {
				cells = append(cells, []int{x + distance, y + distance})
			}
			if !inCells(x+distance, y-distance) {
				cells = append(cells, []int{x + distance, y - distance})
			}
			if !inCells(x-distance, y+distance) {
				cells = append(cells, []int{x - distance, y + distance})
			}
		}
	}

	selection.Cells = cells

	return selection

}

// Shrink shrinks the selection by one.
func (selection Selection) Shrink(diagonal bool) Selection {

	type pair struct {
		X, Y int
	}

	counts := make(map[pair]int, 0)

	countCell := func(pos pair) {
		_, ok := counts[pos]
		if !ok {
			counts[pos] = 0
		}
		counts[pos]++
	}

	for _, c := range selection.Cells {

		x, y := c[0], c[1]

		countCell(pair{x + 1, y})
		countCell(pair{x - 1, y})
		countCell(pair{x, y + 1})
		countCell(pair{x, y - 1})
		if diagonal {
			countCell(pair{x + 1, y + 1})
			countCell(pair{x - 1, y - 1})
			countCell(pair{x - 1, y + 1})
			countCell(pair{x + 1, y - 1})
		}

	}

	p := pair{}

	if diagonal {
		return selection.By(func(x, y int) bool {
			p.X, p.Y = x, y
			return counts[p] == 8
		})
	}
	return selection.By(func(x, y int) bool {
		p.X, p.Y = x, y
		return counts[p] == 4
	})

}

// Invert inverts the selection (selects all non-selected cells from the Selection's source Room).
func (selection Selection) Invert() Selection {

	sel := selection.Room.Select()

	return sel.By(func(x, y int) bool {
		return !selection.Contains(x, y)
	})

}

// Contains returns a boolean indicating if the specified cell is in the list of cells contained in the selection.
func (selection *Selection) Contains(x, y int) bool {
	for _, c := range selection.Cells {
		if c[0] == x && c[1] == y {
			return true
		}
	}
	return false
}

// Fill fills the cells in the Selection with the rune provided.
func (selection Selection) Fill(char rune) Selection {
	return selection.By(func(x, y int) bool {
		selection.Room.Set(x, y, char)
		return true
	})
}

// Degrade applies a formula that randomly sets the cells in the selection to the provided char rune if their neighbors have
// that same rune value. The more neighbors that have the rune value, the more likely the selected cell will be set to it as well.
func (selection Selection) Degrade(char rune) Selection {

	// Basically, if a cell has 1 neighbor being the value, 15% chance to turn into the value, 2 sides = 25%, 3 sides = 50%

	return selection.By(func(x, y int) bool {
		c := selection.Room.Get(x, y)
		multiplier := 0
		if selection.Room.Get(x-1, y) == char {
			multiplier++
		}
		if selection.Room.Get(x+1, y) == char {
			multiplier++
		}
		if selection.Room.Get(x, y-1) == char {
			multiplier++
		}
		if selection.Room.Get(x, y+1) == char {
			multiplier++
		}

		if rand.Float32() <= float32(multiplier)*.025 {
			c = char
		}
		selection.Room.Set(x, y, c)
		return true
	})

}

func randBool() bool {
	r := rand.Float32()
	if r >= .5 {
		return true
	}
	return false
}
