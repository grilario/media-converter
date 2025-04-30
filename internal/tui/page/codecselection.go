package page

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grilario/video-converter/internal/app"
	"github.com/grilario/video-converter/internal/tui/styles"
	"github.com/grilario/video-converter/internal/tui/util"
	"github.com/grilario/video-converter/pkg/ffmpeg"
)

var (
	audioCodecs = ffmpeg.ListAudioCodecs()
	videoCodecs = ffmpeg.ListVideoCodecs()
)

var CodecSelectionPage PageID = "codecSelection"

type codecSelectionPage struct {
	app            *app.App
	cursor         int
	selectedStream *ffmpeg.Stream
	choices        []ffmpeg.Codec
	nchoices       int // should len of choices + 1 its include remove and back options
}

func NewCodecSelectionPage(app *app.App, selectedStream *ffmpeg.Stream) tea.Model {
	choices := []ffmpeg.Codec{}
	switch selectedStream.Kind() {
	case "video":
		choices = videoCodecs

	case "audio":
		choices = audioCodecs
	}

	return codecSelectionPage{
		app:            app,
		cursor:         0,
		selectedStream: selectedStream,
		choices:        choices,
		nchoices:       len(choices) + 1,
	}
}

func (p codecSelectionPage) Init() tea.Cmd {
	return nil
}

func (p codecSelectionPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, util.DefaultKeyMap.Next):
			return p.choose()

		case key.Matches(msg, util.DefaultKeyMap.Back):
			back := p.nchoices
			p.cursor = back
			return p.choose()

		case key.Matches(msg, util.DefaultKeyMap.Down):
			if p.cursor > 0 {
				p.cursor--
			}

		case key.Matches(msg, util.DefaultKeyMap.Up):
			if p.cursor < p.nchoices {
				p.cursor++
			}
		}
	}

	return p, nil
}

func (p codecSelectionPage) choose() (tea.Model, tea.Cmd) {
	removeOption := p.nchoices - 1
	backOption := p.nchoices

	switch p.cursor {
	case removeOption:
		p.app.Media.ConfigStream(p.selectedStream, ffmpeg.Config{Remove: true})

	case backOption:
		// do nothing

	default:
		p.app.Media.ConfigStream(p.selectedStream, ffmpeg.Config{Codec: p.choices[p.cursor], Remove: false})
	}

	return p, util.CmdHandler(PageChangeMsg{StreamSelectionPage})
}

func (p codecSelectionPage) View() string {
	var view strings.Builder

	for i, codec := range p.choices {
		cursor := styles.GetCursor(p.cursor, i)

		fmt.Fprintf(&view, "%s %s \n", cursor, codec)
	}

	removeCursor := styles.GetCursor(p.cursor, p.nchoices-1)
	fmt.Fprintf(&view, "%s Remove \n", removeCursor)

	backCursor := styles.GetCursor(p.cursor, p.nchoices)
	fmt.Fprintf(&view, "%s Back \n", backCursor)

	return view.String()
}
