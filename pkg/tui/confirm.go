package tui

import tea "github.com/charmbracelet/bubbletea"

type Confirm struct {
	cursor int
}

func (c Confirm) Update(msg Msg, app App) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "l", "left", "h", "right":
			if c.cursor == 0 {
				c.cursor = 1
			} else {
				c.cursor = 0
			}

		case "enter":
			app.workerChannel <- *app.media
			app.page = Progress{}
			return app, nil
		}
	}

	app.page = c
	return app, nil
}

func (c Confirm) View(app App) string {
	if c.cursor == 0 {
		return "\n   CONFIRM  Cancel   \n"
	}
	return "\n   Confirm  CANCEL   \n"
}
