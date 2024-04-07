package prestafford2

import (
	"math"
	"math/rand"

	"github.com/makyo/gogol/rle"
)

type model struct {
	width   int
	height  int
	field   []cell
	changes []int
}

// wrapPos wraps a cell position that would otherwise be outside of a rectangular grid.
func (m *model) wrapPos(pos int) int {
	return int(math.Abs(float64(pos))) % len(m.field)
}

func (m *model) state(index, pos int) bool {
	switch pos {
	case 0:
		return m.field[index].leftState()
	case 1:
		return m.field[index].middleState()
	case 2:
		return m.field[index].rightState()
	}
	return false
}

// makeAlive sets the cell state to alive and increments the neighbor count. It is idemptotent, assumes the cell is dead, and if it became alive on the current generation, appends the position to the list of changes.
func (m *model) makeAlive(index, pos int) bool {
	switch pos {
	case 0:
		if m.field[index].leftState() {
			return false
		}
		// East
		m.field[m.wrapPos(index-1)] += rightCountone

		// Northeast
		m.field[m.wrapPos(index-m.width-1)] += rightCountone

		// North and Northwest are in one triplet
		m.field[m.wrapPos(index-m.width)] += leftCountone + middleCountone

		// West — We don't need to do because that's the middle cell, which accounts for the left cell being set

		// Southwest and South are in one triplet
		m.field[m.wrapPos(index+m.width)] += leftCountone + middleCountone

		// Southeast
		m.field[m.wrapPos(index+m.width-1)] += rightCountone

		m.field[index] = m.field[index].setLeftState(true)
		m.changes = append(m.changes, index, m.wrapPos(index-1), m.wrapPos(index-m.width-1), m.wrapPos(index+m.width), m.wrapPos(index-m.width), m.wrapPos(index+m.width-1))
	case 1:
		if m.field[index].middleState() {
			return false
		}
		// Northern triplet
		m.field[m.wrapPos(index-m.width)] += leftCountone + middleCountone + rightCountone

		// Southern triplet
		m.field[m.wrapPos(index+m.width)] += leftCountone + middleCountone + rightCountone

		m.field[index] = m.field[index].setMiddleState(true)
		m.changes = append(m.changes, index, m.wrapPos(index+m.width), m.wrapPos(index-m.width))
	case 2:
		if m.field[index].rightState() {
			return false
		}
		// West
		m.field[m.wrapPos(index+1)] += leftCountone

		// Northwest
		m.field[m.wrapPos(index-m.width+1)] += leftCountone

		// North and Northeast are in one triplet
		m.field[m.wrapPos(index-m.width)] += rightCountone + middleCountone

		// East — We don't need to do because that's the middle cell, which accounts for the left cell being set

		// Southeast and South are in one triplet
		m.field[m.wrapPos(index+m.width)] += rightCountone + middleCountone

		// Southwest
		m.field[m.wrapPos(index+m.width+1)] += leftCountone

		m.field[index] = m.field[index].setRightState(true)
		m.changes = append(m.changes, index, m.wrapPos(index+1), m.wrapPos(index-m.width+1), m.wrapPos(index+m.width), m.wrapPos(index-m.width), m.wrapPos(index+m.width+1))
	default:
		return false
	}
	m.changes = append(m.changes, index, m.wrapPos(index-m.width), m.wrapPos(index+m.width))
	return true
}

// makeDead sets the cell state to dead and decrements the neighbor count. It is idemptotent, assumes the cell is alive, and if it became dead on the current generation, appends the position to the list of changes.
func (m *model) makeDead(index, pos int) bool {
	switch pos {
	case 0:
		if !m.field[index].leftState() {
			return false
		}
		// East
		m.field[m.wrapPos(index-1)] -= rightCountone

		// Northeast
		m.field[m.wrapPos(index-m.width-1)] -= rightCountone

		// North and Northwest are in one triplet
		m.field[m.wrapPos(index-m.width)] -= leftCountone + middleCountone

		// West — We don't need to do because that's the middle cell, which accounts for the left cell being set

		// Southwest and South are in one triplet
		m.field[m.wrapPos(index+m.width)] -= leftCountone + middleCountone

		// Southeast
		m.field[m.wrapPos(index+m.width-1)] -= rightCountone

		m.field[index] = m.field[index].setLeftState(false)
		m.changes = append(m.changes, index, m.wrapPos(index-1), m.wrapPos(index-m.width-1), m.wrapPos(index+m.width), m.wrapPos(index-m.width), m.wrapPos(index+m.width-1))
	case 1:
		if !m.field[index].middleState() {
			return false
		}
		// Northern triplet
		m.field[m.wrapPos(index-m.width)] -= leftCountone + middleCountone + rightCountone

		// Southern triplet
		m.field[m.wrapPos(index+m.width)] -= leftCountone + middleCountone + rightCountone

		m.field[index] = m.field[index].setMiddleState(false)
		m.changes = append(m.changes, index, m.wrapPos(index+m.width), m.wrapPos(index-m.width))
	case 2:
		if !m.field[index].rightState() {
			return false
		}
		// West
		m.field[m.wrapPos(index+1)] -= leftCountone

		// Northwest
		m.field[m.wrapPos(index-m.width+1)] -= leftCountone

		// North and Northeast are in one triplet
		m.field[m.wrapPos(index-m.width)] -= rightCountone + middleCountone

		// East — We don't need to do because that's the middle cell, which accounts for the left cell being set

		// Southeast and South are in one triplet
		m.field[m.wrapPos(index+m.width)] -= rightCountone + middleCountone

		// Southwest
		m.field[m.wrapPos(index+m.width+1)] -= leftCountone

		m.field[index] = m.field[index].setRightState(false)
		m.changes = append(m.changes, index, m.wrapPos(index+1), m.wrapPos(index-m.width+1), m.wrapPos(index+m.width), m.wrapPos(index-m.width), m.wrapPos(index+m.width+1))
	default:
		return false
	}
	return true
}

// nextGeneration evolves the field of automata one generation based on the rules of Conway's Game of Life.
func (m *model) Next() {
	// Deep copy the field and changelist
	currentChanges := []int{}

	// Loop over the list of changed triplets.
	for _, change := range m.changes {
		t := m.field[change]

		lc, mc, rc := t.leftNeighbors(), t.middleNeighbors(), t.rightNeighbors()
		t = t.setLeftNext(lc == 3 || t.leftState() && lc == 2)
		t = t.setMiddleNext(mc == 3 || t.middleState() && mc == 2)
		t = t.setRightNext(rc == 3 || t.rightState() && rc == 2)

		if t.changed() {
			currentChanges = append(currentChanges, change)
			m.field[change] = t
		}
	}
	m.changes = []int{}
	for _, change := range currentChanges {
		t := m.field[change]

		if !t.changed() {
			continue
		}

		if t.leftNext() {
			m.makeAlive(change, 0)
		} else {
			m.makeDead(change, 0)
		}

		if t.middleNext() {
			m.makeAlive(change, 1)
		} else {
			m.makeDead(change, 1)
		}

		if t.rightNext() {
			m.makeAlive(change, 2)
		} else {
			m.makeDead(change, 2)
		}
	}
}

// Populate generates a random field of automata, where each cell has a 1 in 5 chance of being alive.
func (m *model) Populate() {
	for i, _ := range m.field {
		m.field[i] = cell(0)
		for c := 0; c < 3; c++ {
			if rand.Intn(5) == 0 {
				m.makeAlive(i, c)
			}
		}
	}
}

// Ingest sets the field to the given value.
func (m *model) Ingest(f *rle.RLEField) {
	startX := (m.width - f.Width) / 2
	startY := (m.height - f.Height) / 2
	for y, row := range f.Field {
		for x, col := range row {
			if col {
				pos := (y+startY)*m.width + x + startX
				m.makeAlive(pos/3, pos%3)
			}
		}
	}
}

func (m *model) ToggleCell(x, y int) {
	pos := y*m.width + x
	if m.state(pos/3, pos%3) {
		m.makeDead(pos/3, pos%3)
	} else {
		m.makeAlive(pos/3, pos%3)
	}
}

// View builds the entire screen's worth of cells to be printed by returning a • for a living cell or a space for a dead cell.
func (m *model) String() string {
	var frame string

	// Loop the field.
	for i, _ := range m.field {

		// Set the cell contents
		for j := 0; j < 3; j++ {
			if m.state(i, j) {
				frame += "•"
			} else {
				frame += "."
			}
		}
		if (i+1)%m.width == 0 {
			frame += "\n"
		}

	}
	return frame
}

func New(width, height int) *model {
	m := &model{
		width:   width,
		height:  height,
		field:   make([]cell, width*height),
		changes: []int{},
	}
	return m
}
