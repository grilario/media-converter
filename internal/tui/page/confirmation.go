package page

import (
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
		switch msg.String() {
		case "l", "left", "h", "right":
			if c.cursor == 0 {
				c.cursor = 1
			} else {
				c.cursor = 0
			}

		case "enter":
			// if cancel option
			if c.cursor == 1 {
				return c, util.CmdHandler(PageChangeMsg{StreamSelectionPage})
			}

			return c, util.CmdHandler(PageChangeMsg{ProgressPage})
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
