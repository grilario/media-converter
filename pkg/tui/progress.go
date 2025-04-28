package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Progress struct {
	progress float64
}

func (p Progress) Update(msg Msg, app App) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ProgressMsg:
		if msg.error != nil {
			panic(msg.error)
			// return app, tea.Quit
		}

		p.progress = msg.progress
		app.page = p
	}

	return app, nil
}

func (c Progress) View(app App) string {
	return fmt.Sprintf(" Progress: %3.f%%   ", c.progress*100)
}
