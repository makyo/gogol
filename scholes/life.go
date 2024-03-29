package scholes

import (
	"math/rand"

	"github.com/makyo/gogol/base"
)

type model struct {
	wrap  bool
	width int
	field []int
}

func moveLeft(field []int, width int) []int {
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

func moveRight(field []int, width int) []int {
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

func moveUp(field []int, width int) []int {
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

func moveDown(field []int, width int) []int {
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

func sumShifted(a, b []int) []int {
	result := make([]int, len(a))
	for i, cell := range a {
		result[i] = cell + b[i]
	}
	return result
}

func sumField(field []int, width int) []int {
	result := make([]int, len(field))
	copy(result, field)

	// Move west, sum
	shifted := moveLeft(field, width)
	result = sumShifted(result, shifted)

	// Northwest
	shifted = moveUp(shifted, width)
	result = sumShifted(result, shifted)

	// North
	shifted = moveRight(shifted, width)
	result = sumShifted(result, shifted)

	// Northeast
	shifted = moveRight(shifted, width)
	result = sumShifted(result, shifted)

	// East
	shifted = moveDown(shifted, width)
	result = sumShifted(result, shifted)

	// Southeast
	shifted = moveDown(shifted, width)
	result = sumShifted(result, shifted)

	// South
	shifted = moveLeft(shifted, width)
	result = sumShifted(result, shifted)

	// Southwest
	shifted = moveLeft(shifted, width)
	result = sumShifted(result, shifted)

	return result
}

// Next evolves the field one generation.
func (m model) Next() base.Model {
	sum := sumField(m.field, m.width)
	for i, cell := range sum {
		if m.field[i] == 1 {
			if cell < 3 || cell > 4 {
				m.field[i] = 0
			}
		} else {
			if cell == 3 {
				m.field[i] = 1
			}
		}
	}
	return m
}

// Populate generates a random field of automata, where each cell has a 1 in 5 chance of being alive.
func (m model) Populate() base.Model {
	for i, _ := range m.field {
		if rand.Intn(5) == 0 {
			m.field[i] = 1
		}
	}
	return m
}

// ToggleCell toggles whether the given cell is alive or dead.
func (m model) ToggleCell(x, y int) base.Model {
	pos := y*m.width + x
	if m.field[pos] == 1 {
		m.field[pos] = 0
	} else {
		m.field[pos] = 1
	}
	return m
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
