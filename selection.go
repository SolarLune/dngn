package dngn

import "math/rand"

// A Selection represents a selection of cell positions in the Layout's data array, and can be filtered down and manipulated
// using the functions on the Selection struct. You can use Selections to manipulate a
type Selection struct {
	Layout *Layout
	Cells  map[Position]bool
}

func (selection *Selection) Clone() Selection {
	newSelection := Selection{
		Layout: selection.Layout,
		Cells:  map[Position]bool{},
	}
	for key := range selection.Cells {
		newSelection.Cells[key] = true
	}
	return newSelection
}

// FilterByRune filters the Selection down to the cells that have the character (rune) provided.
func (selection Selection) FilterByRune(value rune) Selection {
	return selection.FilterBy(func(x, y int) bool {
		return selection.Layout.Get(x, y) == value
	})
}

// All returns a selection with all cells from the Layout selected.
func (selection Selection) All() Selection {
	newSelection := selection.Clone()
	for y := 0; y < len(newSelection.Layout.Data); y++ {
		for x := 0; x < len(newSelection.Layout.Data[y]); x++ {
			newSelection.Cells[Position{x, y}] = true
		}
	}
	return newSelection
}

// None returns a selection with no selected cells from the Layout.
func (selection Selection) None() Selection {
	newSelection := Selection{
		Layout: selection.Layout,
		Cells:  map[Position]bool{},
	}
	return newSelection
}

// FilterByPercentage selects the provided percentage (from 0 - 1) of the cells curently in the Selection.
func (selection Selection) FilterByPercentage(percentage float32) Selection {

	return selection.FilterBy(func(x, y int) bool {
		return rand.Float32() <= percentage
	})

}

// FilterByArea filters down a selection by only selecting the cells that have X and Y values between X, Y, and X+W and Y+H.
// It crops the selection, basically.
func (selection Selection) FilterByArea(x, y, w, h int) Selection {

	return selection.FilterBy(func(cx, cy int) bool {
		return cx >= x && cy >= y && cx <= x+w-1 && cy <= y+h-1
	})

}

// Add returns a clone of the current Selection with the cells in the other Selection.
func (selection Selection) Add(other Selection) Selection {

	newSelection := selection.Clone()

	for position := range other.Cells {
		newSelection.Cells[position] = true
	}

	return newSelection

}

// Remove returns a clone of the current Selection without the cells in the other Selection.
func (selection Selection) Remove(other Selection) Selection {

	newSelection := selection.Clone()

	for position := range other.Cells {
		delete(newSelection.Cells, position)
	}

	return newSelection

}

// FilterByNeighbor returns a filtered Selection of the cells that are surrounded at least by minNeighborCount neighbors with a value of
// neighborValue. If diagonals is true, then diagonals are also checked. If atMost is true, then FilterByNeighbor will only
// work if there's at MOST that many neighbors.
func (selection Selection) FilterByNeighbor(neighborValue rune, minNeighborCount int, diagonals bool, atMost bool) Selection {

	return selection.FilterBy(func(x, y int) bool {

		n := 0

		if selection.Layout.Get(x-1, y) == neighborValue {
			n++
		}
		if selection.Layout.Get(x+1, y) == neighborValue {
			n++
		}
		if selection.Layout.Get(x, y-1) == neighborValue {
			n++
		}
		if selection.Layout.Get(x, y+1) == neighborValue {
			n++
		}

		if diagonals {
			if selection.Layout.Get(x-1, y-1) == neighborValue {
				n++
			}
			if selection.Layout.Get(x+1, y-1) == neighborValue {
				n++
			}
			if selection.Layout.Get(x-1, y+1) == neighborValue {
				n++
			}
			if selection.Layout.Get(x+1, y+1) == neighborValue {
				n++
			}
		}

		if atMost {
			return n <= minNeighborCount
		}

		return n >= minNeighborCount

	})

}

// FilterBy takes a function that takes the X and Y values of each cell position contained in the Selection, and returns a
// boolean to indicate whether to include that cell in the Selection or not. If the result is true, the cell is included in the Selection;
// Otherwise, it is filtered out. This allows you to easily make custom filtering functions to filter down the cells in a Selection.
func (selection Selection) FilterBy(filterFunc func(x, y int) bool) Selection {

	// Note that while we're assigning the Cells variable of selection here directly,
	// because this function doesn't take a pointer notation, we're operating on a copy
	// of the selection, not the original.

	newSelection := selection.Clone()

	cells := map[Position]bool{}

	for c := range newSelection.Cells {
		if filterFunc(c.X, c.Y) {
			cells[c] = true
		}
	}

	newSelection.Cells = cells

	return newSelection

}

// Select attempts to select a number of the cells contained within the selection and returns them. If there's fewer cells in the selection,
// then it will simply return the entirety of the selection.
func (selection Selection) Select(num int) []Position {

	cells := []Position{}

	for cell := range selection.Cells {
		cells = append(cells, cell)
		if len(cells) >= num {
			return cells
		}
	}

	return cells

}

// Expand expands the selection outwards by the distance value provided. Diagonal indicates if the expansion should happen
// diagonally as well, or just on the cardinal 4 directions. If a negative value is given for distance, it shrinks the selection.
func (selection Selection) Expand(distance int, diagonal bool) Selection {

	newSelection := selection.Clone()

	if distance == 0 {
		return newSelection
	}

	shrinking := false
	if distance < 0 {
		shrinking = true
		distance *= -1
	}

	for i := 0; i < distance; i++ {

		// We can't loop through the cells while modifying them, so we'll make a copy after each iteration.
		cells := map[Position]bool{}

		for c := range newSelection.Cells {
			cells[c] = true
		}

		toRemove := []Position{}

		for cp := range cells {

			if shrinking {

				if !newSelection.Contains(cp.X-1, cp.Y) || !newSelection.Contains(cp.X+1, cp.Y) || !newSelection.Contains(cp.X, cp.Y-1) || !newSelection.Contains(cp.X, cp.Y+1) {

					if !diagonal || !newSelection.Contains(cp.X-1, cp.Y-1) || !newSelection.Contains(cp.X+1, cp.Y-1) || !newSelection.Contains(cp.X-1, cp.Y+1) || !newSelection.Contains(cp.X+1, cp.Y+1) {

						toRemove = append(toRemove, cp)

					}

				}

			} else {

				newSelection.AddPosition(cp.X-1, cp.Y)
				newSelection.AddPosition(cp.X+1, cp.Y)
				newSelection.AddPosition(cp.X, cp.Y-1)
				newSelection.AddPosition(cp.X, cp.Y+1)

				if diagonal {
					newSelection.AddPosition(cp.X-1, cp.Y-1)
					newSelection.AddPosition(cp.X-1, cp.Y+1)
					newSelection.AddPosition(cp.X+1, cp.Y-1)
					newSelection.AddPosition(cp.X+1, cp.Y+1)
				}

			}

		}

		for _, c := range toRemove {
			newSelection.RemovePosition(c.X, c.Y)
		}

	}

	return newSelection

}

// Invert inverts the selection (selects all non-selected cells from the Selection's source Map).
func (selection Selection) Invert() Selection {

	inverted := selection.Layout.Select()

	return inverted.FilterBy(func(x, y int) bool {
		return !selection.Contains(x, y)
	})

}

// Contains returns a boolean indicating if the specified cell is in the list of cells contained in the selection.
func (selection *Selection) Contains(x, y int) bool {
	for c := range selection.Cells {
		if c.X == x && c.Y == y {
			return true
		}
	}
	return false
}

// Fill fills the cells in the Selection with the rune provided.
func (selection Selection) Fill(char rune) Selection {
	return selection.FilterBy(func(x, y int) bool {
		selection.Layout.Set(x, y, char)
		return true
	})
}

// AddPosition adds a specific position to the Selection. If the position lies outside of the layout's area, then it's removed.
func (selection *Selection) AddPosition(x, y int) {

	if x < 0 || y < 0 || x >= selection.Layout.Width || y >= selection.Layout.Height {
		return
	}

	selection.Cells[Position{x, y}] = true

}

// RemovePosition removes a specific position from the Selection.
func (selection *Selection) RemovePosition(x, y int) {
	delete(selection.Cells, Position{x, y})
}
