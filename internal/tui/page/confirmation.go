package page

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/grilario/video-converter/internal/tui/util"
)

var ConfirmationPage PageID = "confirmationPage"

type ConfirmKeyMap struct {
	Right key.Binding
	Left  key.Binding
	Enter key.Binding
}

var confirmKeyMap = ConfirmKeyMap{
	Right: key.NewBinding(
		key.WithKeys("k", "right"),
		key.WithHelp("k/🡒", "move right"),
	),
	Left: key.NewBinding(
		key.WithKeys("j", "left"),
		key.WithHelp("k/🡐", "move left"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "choose"),
	),
}

type confirmationPage struct {
	cursor int

	headerStyle         lipgloss.Style
	choiceStyle         lipgloss.Style
	selectedChoiceStyle lipgloss.Style
	contentStyle        lipgloss.Style

	helpStyles    lipgloss.Style
	helpContainer help.Model
}

func NewConfirmationPage() tea.Model {
	return confirmationPage{
		cursor: 0,

		headerStyle:         lipgloss.NewStyle().Bold(true).Margin(1),
		choiceStyle:         lipgloss.NewStyle().Bold(true).Padding(0, 2).Background(lipgloss.Color("#5A56E0")),
		selectedChoiceStyle: lipgloss.NewStyle().Bold(true).Underline(true).Padding(0, 2).Background(lipgloss.Color("#03A64A")),
		contentStyle:        lipgloss.NewStyle().Margin(1),

		helpStyles:    lipgloss.NewStyle().Margin(1),
		helpContainer: help.New(),
	}
}

func (p confirmationPage) Init() tea.Cmd {
	return nil
}

func (p confirmationPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, confirmKeyMap.Right):
			p.cursor = 1

		case key.Matches(msg, confirmKeyMap.Left):
			p.cursor = 0

		case key.Matches(msg, confirmKeyMap.Enter):
			// if cancel option
			if p.cursor == 1 {
				return p, util.CmdHandler(PageChangeMsg{StreamSelectionPage})
			}

			return p, util.CmdHandler(PageChangeMsg{ProgressPage})
		}
	}

	return p, nil
}

func (p confirmationPage) View() string {
	var view strings.Builder
	var confirm, cancel string

	switch p.cursor {
	case 0:
		confirm = p.selectedChoiceStyle.Render("Confirm")
		cancel = p.choiceStyle.Render("Cancel")

	case 1:
		confirm = p.choiceStyle.Render("Confirm")
		cancel = p.selectedChoiceStyle.Render("Cancel")
	}

	view.WriteString(p.headerStyle.Render("Do you want to continue?"))
	view.WriteString(p.contentStyle.Render(confirm + "  " + cancel))
	view.WriteString(p.helpStyles.Render(p.helpContainer.ShortHelpView(util.KeyMapToSlice(confirmKeyMap))))

	return view.String()
}
