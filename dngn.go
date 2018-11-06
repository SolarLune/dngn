/*
Package dngn is a simple random map generation library primarily made to be used for 2D games. It features a simple API,
and a couple of different means to generate maps. The easiest way to kick things off when using dngn is to simply create a Room
to represent your overall game map, which can then be manipulated or have a Generate function run on it to actually generate the
content on the map.
*/
package dngn

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Room represents a singular Room in the dungeon, or the root that all other Rooms are parented to.
// X and Y are the position of the Room in the layout. This is basically only used for children Rooms where they can occupy a
// space within a larger dungeon generation.
// Width and Height are the width and height of the Room in the layout. If the Room is a parent-less root Room, then this
// determines the size of the overall Data structure backing the Room layout.
// Children is a list of the Rooms the Room has as direct children.
// Parent is the Room that this Room is a child of. If this is the root Room, then this would be nil.
// Data is the core underlying data structure representing the dungeon. It's a 2D array of numbers. If the Room is a root Room,
// then this will be a pointer to the array itself. If this room is a child, then this will be a pointer to the parent's
// pointer. In other words, all Rooms point to the same Data array to simplify editing it.
// Depth is the depth of the Room in the layout. A Room with no Parent will have a Depth of 1, and so on down the line.
// Seed is the seed of the Room to use when doing random generation using the Generate* functions below.
// CustomSeed indicates whether the Seed was customized - if not, then it will default to using the time of the system to have
// random generation each time you use Generate* functions.
type Room struct {
	ID            string
	X, Y          int
	Width, Height int
	Children      []*Room
	Parent        *Room
	Data          [][]int
	Depth         int
	Seed          int64
	CustomSeed    bool
}

// NewRoom returns a new Room. x, y, w, and h get passed to the Room's X, Y, Width, and Height values. If this is the Root room,
// then the X and Y values are essentially ignored. parent is the parent Room if this Room is to be a child Room. If it is a
// root, then this should be nil.
func NewRoom(x, y, w, h int, parent *Room) *Room {

	r := &Room{X: x, Y: y, Width: w, Height: h}
	r.Children = make([]*Room, 0)
	if parent != nil {
		// This is a child Room, so its Data array can reference the Parent's Data array.
		r.Parent = parent
		r.Parent.Children = append(r.Parent.Children, r)

		for i, child := range r.Parent.Children {
			if child == r {
				r.ID = r.Parent.ID + "." + strconv.Itoa(i)
			}
		}

		r.Data = r.Parent.Data
		r.Depth = r.Parent.Depth + 1
	} else {
		// No parent, so this is a root Room that should "own the Data".
		r.Data = [][]int{}
		r.ID = "0"
		for y := 0; y < h; y++ {
			r.Data = append(r.Data, []int{})
			for x := 0; x < w; x++ {
				r.Data[y] = append(r.Data[y], 0)
			}
		}
		r.Depth = 1
	}
	return r

}

// GenerateBSP generates a map in the bounds of the Room specified using BSPs (binary space partitions). It splits the Room up horizontall or vertically,
// along a random line between 20% and 80% in, and then continues splitting it up into smaller chunks until the maximum depth specified is reached.
// More depth means more complexity, generally, though it is random.
// Link: http://www.roguebasin.com/index.php?title=Basic_BSP_Dungeon_generation
func (room *Room) GenerateBSP(fillValue, maxDepth int) *Room {

	if room.CustomSeed {
		rand.Seed(room.Seed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	for true {

		for _, l := range room.Leaves() {
			l.Split(randBool(), 0.2+rand.Float32()*.6)
			if l.Children[0].MinimumSize() <= 2 || l.Children[1].MinimumSize() <= 2 {
				l.ClearChildren()
				continue
			}
		}

		if room.MaxDepth() >= maxDepth {
			break
		}

	}

	room.Border(fillValue)

	return room

}

// GenerateDrunkWalk generates a map in the bounds of the Room specified using drunk walking.
// Link: http://www.roguebasin.com/index.php?title=Random_Walk_Cave_Generation
func (room *Room) GenerateDrunkWalk(fillValue int, percentageFilled float32) *Room {

	if room.CustomSeed {
		rand.Seed(room.Seed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	sx := room.X + rand.Intn(room.Width)
	sy := room.Y + rand.Intn(room.Height)
	fillCount := float32(0)

	totalArea := float32(room.Area())

	for true {

		cell := room.Get(sx, sy)

		if cell != fillValue {
			room.Set(sx, sy, fillValue)
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

		if sx < room.X {
			sx = room.X
		} else if sx >= room.X+room.Width {
			sx = room.X + room.Width - 1
		}

		if sy < room.Y {
			sy = room.Y
		} else if sy >= room.Y+room.Height {
			sy = room.Y + room.Height - 1
		}

		if fillCount/totalArea >= percentageFilled {
			break
		}

	}

	return room

}

// Set sets the value provided in the Room's Data. A convenience function stand-in for "room.Data[y][x] = value".
func (room *Room) Set(x, y, value int) {
	room.Data[y][x] = value
}

// Get returns the value provided in the Room's Data. A convenience function stand-in for "room.Data[y][x]".
func (room *Room) Get(x, y int) int {
	if x < 0 {
		x = 0
	} else if x >= room.X+room.Width {
		x = room.Width - 1
	}
	if y < 0 {
		y = 0
	} else if y >= room.Y+room.Height {
		y = room.Height - 1
	}

	return int(room.Data[y][x])
}

// SetSeed sets a custom seed for random generation.
func (room *Room) SetSeed(seed int64) {
	room.CustomSeed = true
	room.Seed = seed
}

// ClearSeed clears a custom seed set for random generation. When using a clear seed, random generation functions will use the
// system's Unix time.
func (room *Room) ClearSeed() {
	room.CustomSeed = false
}

// IsIn returns if the specified position is in the Room.
func (room *Room) IsIn(x, y int) bool {
	return x >= room.X && x < room.X+room.Width && y >= room.Y && y < room.Y+room.Height
}

// Center returns the center position of the Room.
func (room *Room) Center() (int, int) {
	return room.X + room.Width/2, room.Y + room.Height/2
}

// OpenInto sets the value specified at a random position on the border between the two Rooms if they are neighbors, creating
// a doorway.
func (room *Room) OpenInto(other *Room, toValue int) bool {

	doorwayChoices := make([]int, 0)

	room.Run(func(x, y int) int {

		if y > room.Y && y < room.Y+room.Height-1 && y > other.Y && y < other.Y+room.Height-1 {

			if other.IsIn(x+1, y) { // Right edge of this room, right neighbor's the other Room
				doorwayChoices = append(doorwayChoices, x, y, x+1, y)
			}
			if other.IsIn(x-1, y) {
				doorwayChoices = append(doorwayChoices, x, y, x-1, y)
			}

		}

		if x > room.X && x < room.X+room.Width-1 && x > other.X && x < other.X+room.Width-1 {

			if other.IsIn(x, y+1) {
				doorwayChoices = append(doorwayChoices, x, y, x, y+1)
			}
			if other.IsIn(x, y-1) {
				doorwayChoices = append(doorwayChoices, x, y, x, y-1)
			}

		}

		return room.Get(x, y)
	})

	if len(doorwayChoices) > 0 {

		ci := rand.Intn(len(doorwayChoices)/4) * 4

		room.Set(doorwayChoices[ci], doorwayChoices[ci+1], toValue)
		other.Set(doorwayChoices[ci+2], doorwayChoices[ci+3], toValue)

		return true

	}

	return false

}

// OpenIntoNeighbors automatically loops through all leaf recursive children to open a doorway between the border of the child
// and its neighbors.
func (room *Room) OpenIntoNeighbors(value int) {

	connected := make([]*Room, 0)

	for _, leaf := range room.Leaves() {

		for _, neighbor := range leaf.Neighbors() {

			if neighbor.IsLeaf() {

				if contains(connected, neighbor) {
					continue
				}

				leaf.OpenInto(neighbor, value)

			}

		}

		connected = append(connected, leaf)

	}

}

// Run runs the provided function on each cell in the Room's Data array. The function should take the X and Y position of each
// cell in the Room's Data array, and returns an int to put in the array in that position.
func (room *Room) Run(function func(x, y int) int) *Room {

	for y := room.Y; y < room.Y+room.Height; y++ {
		for x := room.X; x < room.X+room.Width; x++ {
			room.Set(x, y, function(x, y))
		}
	}

	for _, sub := range room.Children { // Make this optional / changeable?
		sub.Run(function)
	}

	return room

}

// Fill fills the Room's Data array with the value provided.
func (room *Room) Fill(value int) *Room {
	room.Run(func(x int, y int) int {
		return value
	})
	return room
}

// Split splits the Room into two different sub-rooms. vertical controls whether the split is vertical or horizontal, and
// percentage is how "far" into the room the split should happen.
func (room *Room) Split(vertical bool, percentage float32) *Room {

	room.ClearChildren()

	if vertical {
		rWidth := int(float32(room.Width) * percentage)
		NewRoom(room.X, room.Y, rWidth, room.Height, room)
		NewRoom(room.X+rWidth, room.Y, room.Width-rWidth, room.Height, room)
	} else {
		rHeight := int(float32(room.Height) * percentage)
		NewRoom(room.X, room.Y, room.Width, rHeight, room)
		NewRoom(room.X, room.Y+rHeight, room.Width, room.Height-rHeight, room)
	}

	return room

}

// ClearChildren clears the room's Children list, making it a leaf Room.
func (room *Room) ClearChildren() {
	room.Children = make([]*Room, 0)
}

// Border outlines all children leaf rooms with the value provided on just two sides, creating a border between them.
func (room *Room) Border(value int) *Room {

	for _, leaf := range room.Leaves() {
		leaf.Run(func(x int, y int) int {
			if (x >= leaf.X && x <= leaf.X+leaf.Width) && (y == leaf.Y) || (y >= leaf.Y && y <= leaf.Y+leaf.Height) && (x == leaf.X) {
				return value
			}
			return leaf.Get(x, y)
		})
	}

	return room

}

// Outline outlines the Room's edge on all four sides with the value provided.
func (room *Room) Outline(value int) *Room {

	room.Run(func(x int, y int) int {
		if (x >= room.X && x <= room.X+room.Width) && (y == room.Y) || (y >= room.Y && y <= room.Y+room.Height) && (x == room.X) ||
			(x >= room.X && x <= room.X+room.Width) && (y == room.Y+room.Height-1) || (y >= room.Y && y <= room.Y+room.Height) && (x == room.X+room.Width-1) {
			return value
		}
		return room.Get(x, y)
	})

	return room

}

// Degrade applies a formula that randomly degrades the cells in the room that have the from value to the to value if their
// neighbors have the to value. The more neighbors that have the to value, the more likely the cell will degrade.
func (room *Room) Degrade(from, to int) *Room {

	// Basically, if a cell has 1 neighbor being the value, 15% chance to turn into the value, 2 sides = 25%, 3 sides = 50%

	room.Run(func(x, y int) int {
		c := room.Get(x, y)
		multiplier := 0
		if c == from {
			if room.Get(x-1, y) == to {
				multiplier++
			}
			if room.Get(x+1, y) == to {
				multiplier++
			}
			if room.Get(x, y-1) == to {
				multiplier++
			}
			if room.Get(x, y+1) == to {
				multiplier++
			}
		}

		if rand.Float32() <= float32(multiplier)*.025 {
			c = to
		}
		return c
	})

	return room

}

// Remap changes all values from the "from" value specified to the "to" value specified.
func (room *Room) Remap(from, to int) *Room {

	room.Run(func(x, y int) int {

		cv := room.Get(x, y)

		if cv == from {
			return to
		}

		return cv

	})

	return room

}

// Right returns the right edge of the room (i.e. room.X + room.Width).
func (room *Room) Right() int {
	return room.X + room.Width - 1
}

// Bottom returns the bottom edge of the room (i.e. room.Y + room.Height).
func (room *Room) Bottom() int {
	return room.Y + room.Height - 1
}

// Neighbors returns all Rooms that are neighbors (share a border) with the calling Room. Note that these include non-leaf
// Rooms, as well, so be aware of that.
func (room *Room) Neighbors() []*Room {

	neighbors := make([]*Room, 0)

	if room.Parent == nil {
		return neighbors
	}

	for _, r := range room.GetRoot().ChildrenRecursive() {

		if r == room {
			continue
		}

		if r.X == room.Right()+1 && r.Bottom() >= room.Y && r.Y <= room.Bottom() {
			neighbors = append(neighbors, r)
			continue
		}

		if r.Right() == room.X-1 && r.Bottom() >= room.Y && r.Y <= room.Bottom() {
			neighbors = append(neighbors, r)
			continue
		}

		if r.Y == room.Bottom()+1 && r.Right() >= room.X && r.X <= room.Right() {
			neighbors = append(neighbors, r)
			continue
		}

		if r.Bottom() == room.Y-1 && r.Right() >= room.X && r.X <= room.Right() {
			neighbors = append(neighbors, r)
			continue
		}

	}

	return neighbors

}

// Sibling returns the sibling Room in a BSP generation map (i.e. the "other" child of the Room's parent).
func (room *Room) Sibling() *Room {
	var other *Room
	if room.Parent != nil {
		other = room.Parent.Children[0]
		if other == room {
			other = room.Parent.Children[1]
		}
	}
	return other
}

// Find allows you to easily walk down a Room's list of children to get to a specific node / Room. Each index passed is used to
// get the calling Room's children, and continues down until Find get to the end of the indices. In other words, it
// turns "room.Children[0].Children[1].Children[0].Children[0].Fill(1)" into "room.Find(0, 1, 0, 0).Fill(1)". An index of -1
// can be used to go back up to a parent.
func (room *Room) Find(indices ...int) *Room {

	current := room

	for _, i := range indices {
		if i >= 0 {
			current = current.Children[i]
		} else {
			current = current.Parent
		}
	}

	return current

}

// IsLeaf returns if the Room is a leaf room (a room that has no children Rooms).
func (room *Room) IsLeaf() bool {
	return len(room.Children) == 0
}

// GetRoot walks up the Room's parents recursively until it can't any longer, returning the root Room.
func (room *Room) GetRoot() *Room {
	base := room
	for true {
		if base.Parent != nil {
			base = base.Parent
		} else {
			break
		}
	}
	return base
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

// MaxDepth returns the maximum overall depth of the map, starting from the referenced room and going down its Children.
func (room *Room) MaxDepth() int {

	depth := 0

	for _, node := range room.ChildrenRecursive() {

		if node.Depth > depth {
			depth = node.Depth
		}

	}

	return depth - (room.Depth + 1)

}

// ChildrenRecursive returns a list of all rooms that are children recursively going down the map.
func (room *Room) ChildrenRecursive() []*Room {

	rooms := make([]*Room, 0)
	rooms = append(rooms, room)

	checks := make([]*Room, 0)
	checks = append(checks, room)

	for true {

		for _, r := range checks {

			checks = checks[1:]
			checks = append(checks, r.Children...)

			if !contains(rooms, r) {
				rooms = append(rooms, r)
			}

		}
		if len(checks) == 0 {
			break
		}
	}

	return rooms
}

// Leaves returns a list of the calling Room's recursive children that are leaf rooms (Rooms that have no children rooms).
func (room *Room) Leaves() []*Room {

	leaves := make([]*Room, 0)

	for _, r := range room.ChildrenRecursive() {

		if len(r.Children) == 0 {
			leaves = append(leaves, r)
		}

	}

	return leaves
}

// DataToString returns the underlying data of the overall Room layout in an easily understood visual format.
// 0's turn into blank spaces when using DataToString, and the column is shown at the left of the map.
func (room *Room) DataToString() string {

	s := ""

	// s += "   "
	// for x := 0; x < len(room.Data[0]); x++ {
	// 	if x%5 == 0 || x == len(room.Data[0])-1 {
	// 		s += fmt.Sprintf("%3d", x)
	// 		// s += strconv.Itoa(x) + " "
	// 	} else {

	// 		s += "  "
	// 	}
	// }
	// s += "\n"

	s += "\n"

	for y := 0; y < len(room.Data); y++ {
		s += fmt.Sprintf("%3d  |", y)
		for x := 0; x < len(room.Data[y]); x++ {
			if room.Data[y][x] == 0 {
				s += fmt.Sprintf("  ")
			} else {
				s += fmt.Sprintf("%v ", room.Data[y][x])
			}
		}
		s += "|\n"
	}

	return s

}

func randBool() bool {
	r := rand.Float32()
	if r >= .5 {
		return true
	}
	return false
}

func randChoiceInt(options ...int) int {
	return options[rand.Intn(len(options))]
}

func contains(rooms []*Room, room *Room) bool {

	for _, r := range rooms {
		if r == room {
			return true
		}
	}
	return false

}
