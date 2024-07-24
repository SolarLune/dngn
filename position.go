package dngn

import "math"

// Position represents a cell within a Selection.
type Position struct {
	X, Y int
}

// DistanceTo returns the distance from one Position to another.
func (position Position) DistanceTo(other Position) float64 {
	return math.Sqrt(float64(math.Pow(float64(other.X-position.X), 2) + math.Pow(float64(other.Y-position.Y), 2)))
}
