package page

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/grilario/video-converter/internal/app"
	"github.com/grilario/video-converter/internal/tui/styles"
	"github.com/grilario/video-converter/internal/tui/util"
	"github.com/grilario/video-converter/pkg/ffmpeg"
)

var StreamSelectionPage PageID = "streamSelection"

type streamSelectionPage struct {
	app      *app.App
	cursor   int
	choices  []*ffmpeg.Stream
	nchoices int

	choiceStyle       lipgloss.Style
	choiceSafeStyle   lipgloss.Style
	choiceWarnStyle   lipgloss.Style
	choiceDangerStyle lipgloss.Style
	contentStyle      lipgloss.Style

	helpStyles    lipgloss.Style
	helpContainer help.Model
}

func NewStreamSelectionPage(app *app.App) tea.Model {
	choices := app.Media.Streams()

	return streamSelectionPage{
		app:      app,
		cursor:   0,
		choices:  choices,
		nchoices: len(choices), // contains apply choice

		choiceStyle:       lipgloss.NewStyle().Bold(true),
		choiceSafeStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#37A603")),
		choiceWarnStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#ffe224")),
		choiceDangerStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#ff2423")),
		contentStyle:      lipgloss.NewStyle().Margin(1, 2).MarginBottom(1),

		helpStyles:    lipgloss.NewStyle().Margin(1).MarginTop(1),
		helpContainer: help.New(),
	}
}

func (p streamSelectionPage) Init() tea.Cmd {
	return nil
}

func (p streamSelectionPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, util.DefaultKeyMap.Next):
			return p.choose()

		case key.Matches(msg, util.DefaultKeyMap.Up):
			if p.cursor > 0 {
				p.cursor--
			}

		case key.Matches(msg, util.DefaultKeyMap.Down):
			if p.cursor < len(p.choices) {
				p.cursor++
			}
		}
	}

	return p, nil
}

func (p streamSelectionPage) choose() (tea.Model, tea.Cmd) {
	// case choice is confirm jump to tab confirm
	if p.cursor == p.nchoices {
		return p, util.CmdHandler(PageChangeMsg{ConfirmationPage})
	}

	selectedStreamChangeMsg := util.CmdHandler(SelectedStreamChangeMsg{p.choices[p.cursor]})
	pageChangeMsg := util.CmdHandler(PageChangeMsg{CodecSelectionPage})

	return p, tea.Sequence(selectedStreamChangeMsg, pageChangeMsg)
}

func (p streamSelectionPage) View() string {
	var view, choices strings.Builder

	for i, choice := range p.choices {
		cursor := styles.GetCursor(p.cursor, i)

		outCodec, _ := choice.OutCodec()

		switch {
		case choice.ShouldRemoved():
			outCodec = p.choiceDangerStyle.Render("Delete")

		case outCodec == ffmpeg.COPY.String():
			outCodec = p.choiceSafeStyle.Render(outCodec)

		default:
			outCodec = p.choiceWarnStyle.Render(outCodec)
		}

		entryCodec, _ := choice.EntryCodec()
		entry := fmt.Sprintf("%s - %s", choice.Kind(), entryCodec)
		entry = lipgloss.NewStyle().Width(20).Render(entry)

		out := lipgloss.PlaceHorizontal(15, lipgloss.Right, outCodec)

		choices.WriteString(p.choiceStyle.Render(cursor, entry, "🡒", out) + "\n")
	}

	cursor := styles.GetCursor(p.cursor, p.nchoices)
	choices.WriteString(fmt.Sprintf("%s %s", cursor, p.choiceStyle.Render("Confirm")))

	view.WriteString(p.contentStyle.Bold(true).Render("Stream Selection"))
	view.WriteString(p.contentStyle.Render(choices.String()))
	view.WriteString(p.helpStyles.Render(p.helpContainer.ShortHelpView(util.KeyMapToSlice(util.DefaultKeyMap))))

	return view.String()
}
