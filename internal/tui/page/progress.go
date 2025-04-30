package page

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grilario/video-converter/internal/app"
)

var ProgressPage PageID = "progressPage"

type progressPage struct {
	app      *app.App
	progress float64
}

func NewProgressPage(app *app.App) tea.Model {
	return progressPage{
		app:      app,
		progress: 0.0,
	}
}

func (p progressPage) Init() tea.Cmd {
	p.app.Commandc <- app.StartConversion

	return nil
}

func (p progressPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ProgressMsg:
		if msg.Error != nil {
			panic(msg.Error) // todo! improve error handling
		}

		p.progress = msg.Progress
	}
	return p, nil
}

func (c progressPage) View() string {
	return fmt.Sprintf(" Progress: %3.f%%   ", c.progress*100)
}
