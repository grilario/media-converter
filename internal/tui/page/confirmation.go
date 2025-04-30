package page

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grilario/video-converter/internal/tui/util"
)

var ConfirmationPage PageID = "confirmationPage"

type confirmationPage struct {
	cursor int
}

func NewConfirmationPage() tea.Model {
	return confirmationPage{
		cursor: 0,
	}
}

func (p confirmationPage) Init() tea.Cmd {
	return nil
}

func (c confirmationPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, util.DefaultKeyMap.Up, util.DefaultKeyMap.Down):
			if c.cursor == 0 {
				c.cursor = 1
			} else {
				c.cursor = 0
			}

		case key.Matches(msg, util.DefaultKeyMap.Next):
			// if cancel option
			if c.cursor == 1 {
				return c, util.CmdHandler(PageChangeMsg{StreamSelectionPage})
			}

			return c, util.CmdHandler(PageChangeMsg{ProgressPage})

		case key.Matches(msg, util.DefaultKeyMap.Back):
			return c, util.CmdHandler(PageChangeMsg{StreamSelectionPage})
		}
	}

	return c, nil
}

func (c confirmationPage) View() string {
	if c.cursor == 0 {
		return "\n   CONFIRM  Cancel   \n"
	}
	return "\n   Confirm  CANCEL   \n"
}
