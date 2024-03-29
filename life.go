package main

import (
	"flag"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/makyo/gogol/base"
	"github.com/makyo/gogol/naive1d"
	"github.com/makyo/gogol/naive2d"
	"github.com/makyo/gogol/scholes"
)

type tickMsg time.Time

type model struct {
	base base.Model
}

var (
	wrapFlag = flag.Bool("wrap", false, "Wrap the grid at the edges, treating it like a toroid")
	algoFlag = flag.String("algo", "naive1d", "Which algorithm to use (life1d, life2d)")
	width    = 0
	height   = 0
)

func getModel(width, height int, wrap bool) model {
	m := model{}
	switch *algoFlag {
	case "naive1d":
		m.base = naive1d.New(width, height, wrap)
	case "naive2d":
		m.base = naive2d.New(width, height, wrap)
	case "scholes":
		m.base = scholes.New(width, height, wrap)
	}
	return m
}

// tick updates the model every 1/10 second.
func tick() tea.Cmd {
	return tea.Tick(time.Second/10, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init initializes the model. Since this requires a first WindowSizeMsg, we just send a tickMsg.
func (m model) Init() tea.Cmd {
	return tick()
}

// Update updates the state of the model based on various types of messages.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Key press messages
	case tea.KeyMsg:
		switch msg.String() {

		// Quit on Escape or Ctrl+C
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		// Regenerate the field on Ctrl+R
		case "ctrl+r":
			m = getModel(width, height, *wrapFlag)
			return m, nil
		}

	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			m.base = m.base.ToggleCell(msg.X, msg.Y)
		}
		return m, nil

	// Window size messages — we receive one initially, and then again on every resize
	case tea.WindowSizeMsg:

		// Reset the field to the correct size
		width = msg.Width
		height = msg.Height
		m = getModel(width, height, *wrapFlag)
		m.base = m.base.Populate()

	// Tick messages
	case tickMsg:

		// Evolve the next generation
		m.base = m.base.Next()
		return m, tick()
	}
	return m, nil
}

// View builds the entire screen's worth of cells to be printed by returning a • for a living cell or a space for a dead cell.
func (m model) View() string {
	return m.base.String()
}

func main() {
	flag.Parse()
	p := tea.NewProgram(getModel(width, height, *wrapFlag), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
