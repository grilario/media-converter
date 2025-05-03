package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/grilario/video-converter/internal/app"
	"github.com/grilario/video-converter/internal/tui/page"
	"github.com/grilario/video-converter/internal/tui/util"
	"github.com/grilario/video-converter/pkg/ffmpeg"
)

type appModel struct {
	app            *app.App
	currentPage    page.PageID
	selectedStream *ffmpeg.Stream
	pages          map[page.PageID]tea.Model
}

func New(app *app.App) tea.Model {

	return appModel{
		app:            app,
		currentPage:    page.StreamSelectionPage,
		selectedStream: nil,
		pages: map[page.PageID]tea.Model{
			page.StreamSelectionPage: page.NewStreamSelectionPage(app),
			page.ConfirmationPage:    page.NewConfirmationPage(),
			page.ProgressPage:        page.NewProgressPage(app),
		},
	}
}

func (m appModel) Init() tea.Cmd {
	return nil
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, util.DefaultKeyMap.Quit):
			m.app.Commandc <- app.CancelConversion

			return m, tea.Quit
		}

	case page.PageChangeMsg:
		cmd := m.moveToPage(msg.ID)
		return m, cmd

	case page.SelectedStreamChangeMsg:
		m.selectedStream = msg.SelectedStream
		return m, cmd
	}

	m.pages[m.currentPage], cmd = m.pages[m.currentPage].Update(msg)
	return m, cmd
}

func (m *appModel) moveToPage(pageID page.PageID) tea.Cmd {
	if pageID == page.CodecSelectionPage {
		m.pages[pageID] = page.NewCodecSelectionPage(m.app, m.selectedStream)
	}

	m.currentPage = pageID
	cmd := m.pages[m.currentPage].Init()

	return cmd
}

func (m appModel) View() string {
	if m.pages[m.currentPage] == nil {
		return string(m.currentPage)
	}

	return m.pages[m.currentPage].View()
}
