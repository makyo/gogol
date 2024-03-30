package abrashbitwise

import (
	"math"
	"math/rand"

	"github.com/makyo/gogol/rle"
)

// Bitwise operations are not my strong suit, so lots of comments ahead to help me understand.

// A cell is just a byte, which holds both the neighbor count and the state.
type cell byte

const (
	// Store the current state of the cell in the fourth bit (so 0000x000). This is mostly syntactic sugar, meaning...
	state = 4

	// ...if we shift 1 over 4 spots, we get the fourth bit (16).
	statebit = 1 << state

	// We can thus use everything less than that (15) for the count of neighbors.
	countbit = 0xf
)

// state gets the state of the cell by AND-ing the fourth bit. The result will be 0 (so, dead) for anything where 00001000 is not true.
func (c cell) state() bool {
	return (c & statebit) != 0
}

// neighbors gets the count of neighbors by AND-ing everything below the fourth bit (00000111). The result will be the remainder. Essentually mod statebit.
func (c cell) neighbors() cell {
	return c & countbit
}

// vivify makes the cell alive by OR-ing the fourth bit (so, noop if it's already alive).
func (c cell) vivify() cell {
	return c | statebit
}

// kill makes the cell dead by AND-ing the bit with NOT the fourth bit (so, AND-ing with false).
func (c cell) kill() cell {
	return c &^ statebit
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

// addToNeighbors to all neighboring cells.
func (m model) addToNeighbors(x, y int) {
	// West
	newX, newY := m.wrapPos(x-1, y)
	m.field[newY][newX] += 0x1

	// Northwest
	newX, newY = m.wrapPos(x-1, y-1)
	m.field[newY][newX] += 0x1

	// North
	newX, newY = m.wrapPos(x, y-1)
	m.field[newY][newX] += 0x1

	// Northeast
	newX, newY = m.wrapPos(x+1, y-1)
	m.field[newY][newX] += 0x1

	// East
	newX, newY = m.wrapPos(x+1, y)
	m.field[newY][newX] += 0x1

	// Southeast
	newX, newY = m.wrapPos(x+1, y+1)
	m.field[newY][newX] += 0x1

	// South
	newX, newY = m.wrapPos(x, y+1)
	m.field[newY][newX] += 0x1

	// Southwest
	newX, newY = m.wrapPos(x-1, y+1)
	m.field[newY][newX] += 0x1
}

// subtractFromNeighbors subtracts from all neighboring cells.
func (m model) subtractFromNeighbors(x, y int) {
	// West
	newX, newY := m.wrapPos(x-1, y)
	m.field[newY][newX] -= 0x1

	// Northwest
	newX, newY = m.wrapPos(x-1, y-1)
	m.field[newY][newX] -= 0x1

	// North
	newX, newY = m.wrapPos(x, y-1)
	m.field[newY][newX] -= 0x1

	// Northeast
	newX, newY = m.wrapPos(x+1, y-1)
	m.field[newY][newX] -= 0x1

	// East
	newX, newY = m.wrapPos(x+1, y)
	m.field[newY][newX] -= 0x1

	// Southeast
	newX, newY = m.wrapPos(x+1, y+1)
	m.field[newY][newX] -= 0x1

	// South
	newX, newY = m.wrapPos(x, y+1)
	m.field[newY][newX] -= 0x1

	// Southwest
	newX, newY = m.wrapPos(x-1, y+1)
	m.field[newY][newX] -= 0x1
}

// calculateNeighbors calculates alive neighbors for every cell in the model.
func (m model) calculateAllNeighbors() {
	for y, _ := range m.field {
		for x, c := range m.field[y] {
			if c.state() {
				m.addToNeighbors(x, y)
			}
		}
	}
}

// makeAlive sets the cell state to alive and increments the neighbor count. It assumes the cell is dead.
func (m model) makeAlive(x, y int) {
	m.addToNeighbors(x, y)
	m.field[y][x] = m.field[y][x].vivify()
}

// makeDead sets the cell state to dead and decrements the neighbor count. It assumes the cell is alive.
func (m model) makeDead(x, y int) {
	m.subtractFromNeighbors(x, y)
	m.field[y][x] = m.field[y][x].kill()
}

// nextGeneration evolves the field of automata one generation based on the rules of Conway's Game of Life.
func (m model) Next() {
	// Deep copy the field
	next := make([][]cell, m.height)
	for y, _ := range m.field {
		next[y] = make([]cell, m.width)
		for x, _ := range m.field[y] {
			next[y][x] = m.field[y][x]
		}
	}

	// Loop over rows...
	for y, row := range next {

		// Loop over columns...
		for x, c := range row {

			// If a cell has zero alive neighbors and is already dead, it's just going to stay dead.
			if c == 0 {
				continue
			}
			if c.state() {
				if c.neighbors() != 0x2 && c.neighbors() != 0x3 {
					m.makeDead(x, y)
				}
			} else if c.neighbors() == 0x3 {
				m.makeAlive(x, y)
			}
		}
	}
}

// Populate generates a random field of automata, where each cell has a 1 in 5 chance of being alive.
func (m model) Populate() {
	for y, _ := range m.field {
		for x, _ := range m.field[y] {
			m.field[y][x] = 0x0
			if rand.Intn(5) == 0 {
				m.field[y][x] = m.field[y][x].vivify()
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
				m.field[y+startY][x+startX] = m.field[y+startY][x+startX].vivify()
			}
		}
	}
	m.calculateAllNeighbors()
}

func (m model) ToggleCell(x, y int) {
	if m.field[y][x].state() {
		m.field[y][x] = m.field[y][x].vivify()
	} else {
		m.field[y][x] = m.field[y][x].kill()
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
			if col.state() {
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
