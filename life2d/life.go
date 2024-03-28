package life2d

import (
	"math"
	"math/rand"

	"github.com/makyo/gogol/base"
)

type model struct {
	wrap  bool
	width int
	field [][]int
}

// wrapPos wraps a cell position that would otherwise be outside of a rectangular grid.
func (m model) wrapPos(pos int) int {
	return int(math.Abs(float64(pos))) % len(m.field)
}

// nextGeneration evolves the field of automata one generation based on the rules of Conway's Game of Life.
func (m model) Next() base.Model {
	// Create a new field based on the existing one.
	next := make([][]int, len(m.field))
	for i := 0; i < len(m.field); i++ {
		next[i] = make([]int, len(m.field[0]))
	}

	// Loop over rows...
	for y, _ := range m.field {

		// Loop over columns...
		for x, _ := range m.field[y] {
			neighborCount := 0

			// Count the adjacent living cells on the row above.
			if y-1 >= 0 {
				if x-1 >= 0 {
					neighborCount += m.field[y-1][x-1]
				}
				neighborCount += m.field[y-1][x]
				if x+1 < len(m.field[y]) {
					neighborCount += m.field[y-1][x+1]
				}
			}

			// Count the adjacent cells to either side.
			if x-1 >= 0 {
				neighborCount += m.field[y][x-1]
			}
			if x+1 < len(m.field[y]) {
				neighborCount += m.field[y][x+1]
			}

			// Count the adjacent cells on the row below.
			if y+1 < len(m.field) {
				if x-1 >= 0 {
					neighborCount += m.field[y+1][x-1]
				}
				neighborCount += m.field[y+1][x]
				if x+1 < len(m.field[y]) {
					neighborCount += m.field[y+1][x+1]
				}
			}

			// Evolve the current cell by the following rules:
			//
			// 1. A dead cell becomes live if it's surrounded by exactly three living cells to represent breeding.
			// 2. A living cell dies of loneliness if it has 0 or 1 neighbors.
			// 3. A living cell dies of overcrowding if it has more than 3 neighbors.
			// 4. A living cell stays alive if it has 2 or 3 neighbors.
			if m.field[y][x] == 0 {
				if neighborCount == 3 {
					next[y][x] = 1
				} else {
					next[y][x] = 0
				}
			} else if m.field[y][x] == 1 {
				if neighborCount < 2 || neighborCount > 3 {
					next[y][x] = 0
				} else {
					next[y][x] = 1
				}
			}
		}
	}
	m.field = next
	return m
}

// generateField generates a random field of automata, where each cell has a 1 in 5 chance of being alive.
func (m model) Populate() base.Model {
	for y, _ := range m.field {
		for x, _ := range m.field[y] {
			if rand.Intn(5) == 0 {
				m.field[y][x] = 1
			}
		}
	}
	return m
}

func (m model) ToggleCell(x, y int) base.Model {
	if m.field[y][x] == 1 {
		m.field[y][x] = 0
	} else {
		m.field[y][x] = 1
	}
	return m
}

// View builds the entire screen's worth of cells to be printed by returning a • for a living cell or a space for a dead cell.
func (m model) String() string {
	var frame string

	// Loop over rows...
	for _, row := range m.field {

		// Loop over collumns...
		for _, col := range row {

			// Set the cell contents
			if col == 1 {
				frame += "•"
			} else {
				frame += " "
			}
		}

		// Newline at the end of every row.
		frame += "\n"
	}
	return frame
}

func New(width, height int, wrap bool) model {
	m := model{
		wrap:  wrap,
		field: make([][]int, height),
	}
	for i := 0; i < height; i++ {
		m.field[i] = make([]int, width)
	}
	return m
}
