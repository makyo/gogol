package scholes

import (
	"math/rand"
)

type model struct {
	width int
	field []int
}

// shiftLeft returns a copy of the field which has been moved left, wrapping around on the far edge.
func shiftLeft(field []int, width int) []int {
	shifted := make([]int, len(field))
	for i, cell := range field {
		if i == 0 {
			shifted[len(shifted)-1] = cell
		} else {
			shifted[i-1] = cell
		}
	}
	return shifted
}

// shiftRight returns a copy of the field which has been moved right, wrapping around on the far edge.
func shiftRight(field []int, width int) []int {
	shifted := make([]int, len(field))
	for i, cell := range field {
		if i == len(shifted)-1 {
			shifted[0] = cell
		} else {
			shifted[i+1] = cell
		}
	}
	return shifted
}

// shiftUp returns a copy of the field which has been moved up, wrapping around on the far edge.
func shiftUp(field []int, width int) []int {
	shifted := make([]int, len(field))
	for i, cell := range field {
		if i < width {
			shifted[len(shifted)-width-i] = cell
		} else {
			shifted[i-width] = cell
		}
	}
	return shifted
}

// shiftDown returns a copy of the field which has been moved down, wrapping around on the far edge.
func shiftDown(field []int, width int) []int {
	shifted := make([]int, len(field))
	for i, cell := range field {
		if i > len(shifted)-width-1 {
			shifted[len(shifted)-i] = cell
		} else {
			shifted[i+width] = cell
		}
	}
	return shifted
}

func where(field []int, which int) []int {
	result := make([]int, len(field))
	for i, cell := range field {
		if cell == which {
			result[i] = 1
		}
	}
	return result
}

// sumField sums the values of the field along with the field shifted one space to each of the cardinal and ordinal compass points.
func sumField(field []int, width int) []int {
	result := make([]int, len(field))

	west := shiftLeft(field, width)
	northwest := shiftUp(west, width)
	north := shiftRight(northwest, width)
	northeast := shiftRight(north, width)
	east := shiftDown(northeast, width)
	southeast := shiftDown(east, width)
	south := shiftLeft(southeast, width)
	southwest := shiftLeft(south, width)

	for i, cell := range field {
		result[i] = cell + west[i] + northwest[i] + north[i] + northeast[i] + east[i] + southeast[i] + south[i] + southwest[i]
	}

	return result
}

// Next evolves the field one generation.
func (m model) Next() {
	sum := sumField(m.field, m.width)
	threes := where(sum, 3)
	fours := where(sum, 4)
	for i, cell := range m.field {
		m.field[i] = threes[i] | fours[i]&cell
	}
}

func (m model) Ingest(field [][]int) {
	m.width = len(field[0])
	m.field = make([]int, m.width*len(field))
	for row, _ := range field {
		for col, _ := range field[row] {
			m.field[row*col] = field[row][col]
		}
	}
}

// Populate generates a random field of automata, where each cell has a 1 in 5 chance of being alive.
func (m model) Populate() {
	for i, _ := range m.field {
		if rand.Intn(5) == 0 {
			m.field[i] = 1
		}
	}
}

// ToggleCell toggles whether the given cell is alive or dead.
func (m model) ToggleCell(x, y int) {
	pos := y*m.width + x
	if m.field[pos] == 1 {
		m.field[pos] = 0
	} else {
		m.field[pos] = 1
	}
}

// View builds the entire screen's worth of cells to be printed by returning a • for a living cell or a space for a dead cell.
func (m model) String() string {
	var frame string

	// Loop over rows...
	for i, cell := range m.field {
		if cell == 1 {
			frame += "•"
		} else {
			frame += " "
		}
		if (i+1)%m.width == 0 {
			frame += "\n"
		}
	}

	return frame
}

func New(width, height int) model {
	return model{
		width: width,
		field: make([]int, width*height),
	}
}
