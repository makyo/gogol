package naive1d

import (
	"math"
	"math/rand"
)

type model struct {
	wrap  bool
	width int
	field []int
}

// wrapPos wraps a cell position that would otherwise be outside of a rectangular grid.
func (m model) wrapPos(pos int) int {
	return int(math.Abs(float64(pos))) % len(m.field)
}

// nextGeneration evolves the field of automata one generation based on the rules of Conway's Game of Life.
func (m model) Next() {
	// Create a new field based on the existing one.
	next := make([]int, len(m.field))

	// Loop over the cells.
	for i, _ := range m.field {
		neighborCount := 0

		// Count the adjacent living cells on the row above.
		if i-m.width > 0 || m.wrap {
			if i%m.width >= 0 || m.wrap {
				neighborCount += m.field[m.wrapPos(i-m.width-1)]
			}
			neighborCount += m.field[m.wrapPos(i-m.width)]
			if (i+1)%m.width != 0 || m.wrap {
				neighborCount += m.field[m.wrapPos(i-m.width+1)]
			}
		}

		// Count the adjacent cells to either side.
		if i%m.width != 0 || m.wrap {
			neighborCount += m.field[m.wrapPos(i-1)]
		}
		if (i < len(m.field)-1 && (i+1)%m.width != 0) || m.wrap {
			neighborCount += m.field[m.wrapPos(i+1)]
		}

		// Count the adjacent cells on the row below.
		if i+m.width < len(m.field) || m.wrap {
			if i%m.width >= 0 || m.wrap {
				neighborCount += m.field[m.wrapPos(i+m.width-1)]
			}
			neighborCount += m.field[m.wrapPos(i+m.width)]
			if (i+1)%m.width != 0 || m.wrap {
				neighborCount += m.field[m.wrapPos(i+m.width+1)]
			}
		}

		// Evolve the current cell by the following rules:
		//
		// 1. A dead cell becomes live if it's surrounded by exactly three living cells to represent breeding.
		// 2. A living cell dies of loneliness if it has 0 or 1 neighbors.
		// 3. A living cell dies of overcrowding if it has more than 3 neighbors.
		// 4. A living cell stays alive if it has 2 or 3 neighbors.
		next[i] = m.field[i]
		if m.field[i] == 0 && neighborCount == 3 {
			next[i] = 1
			continue
		}
		if m.field[i] == 1 && (neighborCount < 2 || neighborCount > 3) {
			next[i] = 0
		}
	}
	for i, _ := range next {
		m.field[i] = next[i]
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

func (m model) Ingest(field [][]int) {
	m.width = len(field[0])
	for row, _ := range field {
		for col, _ := range field[row] {
			m.field[row*col] = field[row][col]
		}
	}
}

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

func New(width, height int, wrap bool) model {
	return model{
		wrap:  wrap,
		width: width,
		field: make([]int, width*height),
	}
}
