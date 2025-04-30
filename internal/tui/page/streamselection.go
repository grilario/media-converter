package page

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
}

func NewStreamSelectionPage(app *app.App) tea.Model {
	choices := app.Media.Streams()

	return streamSelectionPage{
		app:      app,
		cursor:   0,
		choices:  choices,
		nchoices: len(choices), // contains apply choice
	}
}

func (p streamSelectionPage) Init() tea.Cmd {
	return nil
}

func (p streamSelectionPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "l", "right":
			return p.choose()

		case "up", "k":
			if p.cursor > 0 {
				p.cursor--
			}

		case "down", "j":
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
	var s strings.Builder
	for i, choice := range p.choices {
		cursor := styles.GetCursor(p.cursor, i)

		outCodec := "Delete"
		if !choice.ShouldRemoved() {
			name, _ := choice.OutCodec()
			outCodec = fmt.Sprintf("%s", name)
		}

		entryCodec, _ := choice.EntryCodec()

		fmt.Fprintf(&s, "%s %s (%s)   ->   %s \n", cursor, choice.Kind(), entryCodec, outCodec)
	}

	cursor := styles.GetCursor(p.cursor, p.nchoices)
	fmt.Fprintf(&s, "%s Confirm", cursor)

	return s.String()
}
