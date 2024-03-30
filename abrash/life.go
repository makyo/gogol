package abrash

import (
	"math"
	"math/rand"

	"github.com/makyo/gogol/rle"
)

type cell struct {
	state     int
	neighbors int
}

type model struct {
	width  int
	height int
	field  [][]cell
}

// wrapPos wraps a cell position that would otherwise be outside of a rectangular grid.
func (m model) wrapPos(x, y int) (int, int) {
	return int(math.Abs(float64(m.width+x))) % m.width, int(math.Abs(float64(m.height+y))) % m.height
}

// addToNeighbors adds amount to all neighboring cells neighbors amount.
func (m model) addToNeighbors(amount, x, y int) {
	// West
	newX, newY := m.wrapPos(x-1, y)
	m.field[newY][newX].neighbors += amount

	// Northwest
	newX, newY = m.wrapPos(x-1, y-1)
	m.field[newY][newX].neighbors += amount

	// North
	newX, newY = m.wrapPos(x, y-1)
	m.field[newY][newX].neighbors += amount

	// Northeast
	newX, newY = m.wrapPos(x+1, y-1)
	m.field[newY][newX].neighbors += amount

	// East
	newX, newY = m.wrapPos(x+1, y)
	m.field[newY][newX].neighbors += amount

	// Southeast
	newX, newY = m.wrapPos(x+1, y+1)
	m.field[newY][newX].neighbors += amount

	// South
	newX, newY = m.wrapPos(x, y+1)
	m.field[newY][newX].neighbors += amount

	// Southwest
	newX, newY = m.wrapPos(x-1, y+1)
	m.field[newY][newX].neighbors += amount
}

// calculateNeighbors calculates alive neighbors for every cell in the model.
func (m model) calculateAllNeighbors() {
	for y, _ := range m.field {
		for x, c := range m.field[y] {
			if c.state == 1 {
				m.addToNeighbors(1, x, y)
			}
		}
	}
}

// makeAlive sets the cell state to alive and increments the neighbor count. It assumes the cell is dead.
func (m model) makeAlive(x, y int) {
	m.addToNeighbors(1, x, y)
	m.field[y][x].state = 1
}

// makeDead sets the cell state to dead and decrements the neighbor count. It assumes the cell is alive.
func (m model) makeDead(x, y int) {
	m.addToNeighbors(-1, x, y)
	m.field[y][x].state = 0
}

// nextGeneration evolves the field of automata one generation based on the rules of Conway's Game of Life.
func (m model) Next() {
	// Deep copy the field
	next := make([][]cell, m.height)
	for y, _ := range m.field {
		next[y] = make([]cell, m.width)
		for x, _ := range m.field[y] {
			next[y][x].state = m.field[y][x].state
			next[y][x].neighbors = m.field[y][x].neighbors
		}
	}

	// Loop over rows...
	for y, row := range next {

		// Loop over columns...
		for x, c := range row {

			// If a cell has zero alive neighbors, it's just going to stay dead.
			if c.neighbors == 0 && c.state == 0 {
				continue
			}
			if c.state == 1 {
				if c.neighbors != 2 && c.neighbors != 3 {
					m.makeDead(x, y)
				}
			} else if c.neighbors == 3 {
				m.makeAlive(x, y)
			}
		}
	}
}

// Populate generates a random field of automata, where each cell has a 1 in 5 chance of being alive.
func (m model) Populate() {
	for y, _ := range m.field {
		for x, _ := range m.field[y] {
			if rand.Intn(5) == 0 {
				m.field[y][x].state = 1
			}
		}
	}
	m.calculateAllNeighbors()
}

// Ingest sets the field to the given value.
func (m model) Ingest(f *rle.RLEField) {
	startX := (m.width - f.Width) / 2
	startY := (m.height - f.Height) / 2
	for y, row := range f.Field {
		for x, col := range row {
			if col {
				m.field[y+startY][x+startX].state = 1
			}
		}
	}
	m.calculateAllNeighbors()
}

func (m model) ToggleCell(x, y int) {
	if m.field[y][x].state == 1 {
		m.field[y][x].state = 0
	} else {
		m.field[y][x].state = 1
	}
}

// View builds the entire screen's worth of cells to be printed by returning a • for a living cell or a space for a dead cell.
func (m model) String() string {
	var frame string

	// Loop over rows...
	for _, row := range m.field {
		frame += "\n"

		// Loop over collumns...
		for _, col := range row {

			// Set the cell contents
			if col.state == 1 {
				frame += "•"
			} else {
				frame += " "
			}
		}

	}
	return frame
}

func New(width, height int) model {
	m := model{
		width:  width,
		height: height,
		field:  make([][]cell, height),
	}
	for i := 0; i < m.height; i++ {
		m.field[i] = make([]cell, m.width)
	}
	return m
}
