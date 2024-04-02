package abrash1d

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
	field  []cell
}

// wrapPos wraps a cell position that would otherwise be outside of a rectangular grid.
func (m model) wrapPos(pos int) int {
	return int(math.Abs(float64(pos))) % len(m.field)
}

// addToNeighbors to all neighboring cells.
func (m model) addToNeighbors(pos int) {
	// West
	newPos := m.wrapPos(pos - 1)
	m.field[newPos] += 0x1

	// Northwest
	newPos = m.wrapPos(pos - m.width - 1)
	m.field[newPos] += 0x1

	// North
	newPos = m.wrapPos(pos - m.width)
	m.field[newPos] += 0x1

	// Northeast
	newPos = m.wrapPos(pos - m.width + 1)
	m.field[newPos] += 0x1

	// East
	newPos = m.wrapPos(pos + 1)
	m.field[newPos] += 0x1

	// Southeast
	newPos = m.wrapPos(pos + m.width + 1)
	m.field[newPos] += 0x1

	// South
	newPos = m.wrapPos(pos + m.width)
	m.field[newPos] += 0x1

	// Southwest
	newPos = m.wrapPos(pos + m.width - 1)
	m.field[newPos] += 0x1
}

// subtractFromNeighbors subtracts from all neighboring cells.
func (m model) subtractFromNeighbors(pos int) {
	// West
	newPos := m.wrapPos(pos - 1)
	m.field[newPos] -= 0x1

	// Northwest
	newPos = m.wrapPos(pos - m.width - 1)
	m.field[newPos] -= 0x1

	// North
	newPos = m.wrapPos(pos - m.width)
	m.field[newPos] -= 0x1

	// Northeast
	newPos = m.wrapPos(pos - m.width + 1)
	m.field[newPos] -= 0x1

	// East
	newPos = m.wrapPos(pos + 1)
	m.field[newPos] -= 0x1

	// Southeast
	newPos = m.wrapPos(pos + m.width + 1)
	m.field[newPos] -= 0x1

	// South
	newPos = m.wrapPos(pos + m.width)
	m.field[newPos] -= 0x1

	// Southwest
	newPos = m.wrapPos(pos + m.width - 1)
	m.field[newPos] -= 0x1
}

// calculateNeighbors calculates alive neighbors for every cell in the model.
func (m model) calculateAllNeighbors() {
	for i, c := range m.field {
		if c.state() {
			m.addToNeighbors(i)
		}
	}
}

// makeAlive sets the cell state to alive and increments the neighbor count. It assumes the cell is dead.
func (m model) makeAlive(pos int) {
	m.addToNeighbors(pos)
	m.field[pos] = m.field[pos].vivify()
}

// makeDead sets the cell state to dead and decrements the neighbor count. It assumes the cell is alive.
func (m model) makeDead(pos int) {
	m.subtractFromNeighbors(pos)
	m.field[pos] = m.field[pos].kill()
}

// nextGeneration evolves the field of automata one generation based on the rules of Conway's Game of Life.
func (m model) Next() {
	// Deep copy the field
	next := make([]cell, m.width*m.height)
	for i, _ := range m.field {
		next[i] = m.field[i]
	}

	// Loop over the field.
	for i, c := range next {

		// If a cell has zero alive neighbors and is already dead, it's just going to stay dead.
		if c == 0 {
			continue
		}
		if c.state() {
			if c.neighbors() != 0x2 && c.neighbors() != 0x3 {
				m.makeDead(i)
			}
		} else if c.neighbors() == 0x3 {
			m.makeAlive(i)
		}
	}
}

// Populate generates a random field of automata, where each cell has a 1 in 5 chance of being alive.
func (m model) Populate() {
	for i, _ := range m.field {
		m.field[i] = 0x0
		if rand.Intn(5) == 0 {
			m.field[i] = m.field[i].vivify()
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
				pos := (y+startY)*m.width + x + startX
				m.field[pos] = m.field[pos].vivify()
			}
		}
	}
	m.calculateAllNeighbors()
}

func (m model) ToggleCell(x, y int) {
	pos := y*m.width + x
	if m.field[pos].state() {
		m.field[pos] = m.field[pos].vivify()
	} else {
		m.field[pos] = m.field[pos].kill()
	}
}

// View builds the entire screen's worth of cells to be printed by returning a • for a living cell or a space for a dead cell.
func (m model) String() string {
	var frame string

	// Loop the field.
	for i, c := range m.field {

		// Set the cell contents
		if c.state() {
			frame += "•"
		} else {
			frame += " "
		}
		if i%m.width == 0 {
			frame += "\n"
		}

	}
	return frame
}

func New(width, height int) model {
	m := model{
		width:  width,
		height: height,
		field:  make([]cell, width*height),
	}
	return m
}
