package page

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/grilario/video-converter/internal/app"
	"github.com/grilario/video-converter/internal/tui/util"
)

var ProgressPage PageID = "progressPage"

type ProgressKeyMap struct {
	Quit key.Binding
}

var progressKeyMap = ProgressKeyMap{
	Quit: util.DefaultKeyMap.Quit,
}

type progressPage struct {
	app      *app.App
	progress float64

	headerStyle    lipgloss.Style
	indicatorStyle lipgloss.Style
	indicator      progress.Model

	helpStyle lipgloss.Style
	help      help.Model
}

func NewProgressPage(app *app.App) tea.Model {
	return progressPage{
		app:      app,
		progress: 0.0,

		headerStyle:    lipgloss.NewStyle().Bold(true).Margin(1),
		indicatorStyle: lipgloss.NewStyle().Margin(1, 2),
		indicator:      progress.New(),

		helpStyle: lipgloss.NewStyle().PaddingLeft(1).MarginTop(1).MarginBottom(1),
		help:      help.New(),
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

func (p progressPage) View() string {
	header := p.headerStyle.Render("Progress:")
	indicator := p.indicatorStyle.Render(p.indicator.ViewAs(p.progress))
	help := p.helpStyle.Render(p.help.ShortHelpView(util.KeyMapToSlice(progressKeyMap)))

	return header + indicator + help
}
