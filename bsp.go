package dngn

type BSPOptions struct {
	WallValue       rune // Rune value to use for walls
	SplitCount      int  // How many times to split the layout
	DoorValue       rune // Rune value to use for doors / doorways
	MinimumRoomSize int  // Minimum allowed size of each room within the generated BSP layout
}

func NewDefaultBSPOptions() BSPOptions {

	return BSPOptions{
		WallValue:       'x',
		SplitCount:      10,
		DoorValue:       '#',
		MinimumRoomSize: 4,
	}
}

// BSPRoom represents a room generated through Layout.GenerateBSP().
type BSPRoom struct {
	X, Y, W, H  int        // X, Y, Width, and Height of the BSPRoom.
	Connected   []*BSPRoom // The BSPRooms this room is connected to.
	Traversable bool       // Whether the BSPRoom is traversable when using CountHopsTo().
}

func NewBSPRoom(x, y, w, h int) *BSPRoom {
	return &BSPRoom{
		X:           x,
		Y:           y,
		W:           w,
		H:           h,
		Connected:   []*BSPRoom{},
		Traversable: true,
	}
}

// Area returns the area of the BSPRoom (width * height).
func (bsp *BSPRoom) Area() int {
	return bsp.W * bsp.H
}

// MinSize returns the minimum size of the room.
func (bsp *BSPRoom) MinSize() int {
	if bsp.W < bsp.H {
		return bsp.W
	}
	return bsp.H
}

func (bsp *BSPRoom) Center() Position {
	return Position{bsp.X + bsp.W/2, bsp.Y + bsp.H/2}
}

// CountHopsTo will count the number of hops to go from one room to another, by hopping through connected neighbors. If no traversable link between the two rooms found, CountHopsTo will return -1.
func (bsp *BSPRoom) CountHopsTo(room *BSPRoom) int {

	toCheck := append([]*BSPRoom{}, bsp)
	perRoomHopCount := map[*BSPRoom]int{
		bsp: 0,
	}

	for len(toCheck) > 0 {

		next := toCheck[0]

		if next == room {
			return perRoomHopCount[next]
		}

		toCheck = toCheck[1:]

		if !next.Traversable {
			continue
		}

		for _, connected := range next.Connected {

			if _, exists := perRoomHopCount[connected]; !exists {
				toCheck = append(toCheck, connected)
				perRoomHopCount[connected] = perRoomHopCount[next] + 1
			}

		}

	}

	return -1
}

// Disconnect removes the BSPRoom from any of its neighbors' Connected lists, breaking the link between them.
func (bsp *BSPRoom) Disconnect() {

	for _, neighbor := range bsp.Connected {
		for i, me := range neighbor.Connected {
			if me == bsp {
				neighbor.Connected = append(neighbor.Connected[:i], neighbor.Connected[i+1:]...)
				break
			}
		}
	}

	bsp.Connected = []*BSPRoom{}
}

// Necessary returns if the BSPRoom is necessary to facilitate traversal from its neighbors to the rest of the BSP Layout.
func (bsp *BSPRoom) Necessary() bool {

	// If you only have one neighbor, then you're necessary
	if len(bsp.Connected) == 1 {
		return true
	}

	bsp.Traversable = false

	for _, neighbor := range bsp.Connected {

		// If your neighbor is only connected to you, then you're necessary.
		if len(neighbor.Connected) <= 1 {
			bsp.Traversable = true
			return true
		}

		for _, otherNeighbor := range bsp.Connected {

			if otherNeighbor == bsp || otherNeighbor == neighbor {
				continue
			}

			if neighbor.CountHopsTo(otherNeighbor) < 0 {
				bsp.Traversable = true
				return true
			}

		}

	}

	bsp.Traversable = true
	return false

}
