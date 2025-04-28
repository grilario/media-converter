package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grilario/video-converter/pkg/ffmpeg"
)

var (
	audioCodecs = ffmpeg.ListAudioCodecs()
	videoCodecs = ffmpeg.ListVideoCodecs()
)

type CodecChooser struct {
	cursor   int
	choices  []ffmpeg.Codec
	nchoices int // should len of choices + 1 its include remove and back options
}

func NewCodecChooser(selectedStream *ffmpeg.Stream) CodecChooser {
	choices := []ffmpeg.Codec{}
	switch selectedStream.Kind() {
	case "video":
		choices = videoCodecs

	case "audio":
		choices = audioCodecs
	}

	return CodecChooser{
		cursor:   0,
		choices:  choices,
		nchoices: len(choices) + 1,
	}
}

func (c CodecChooser) Update(msg Msg, app App) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "l", "right":
			return c.choose(app)

		case "esc", "h", "left":
			back := c.nchoices
			c.cursor = back
			return c.choose(app)

		case "up", "k":
			if c.cursor > 0 {
				c.cursor--
			}

		case "down", "j":
			if c.cursor < c.nchoices {
				c.cursor++
			}
		}
	}

	app.page = c
	return app, nil
}

func (c CodecChooser) choose(app App) (tea.Model, tea.Cmd) {
	removeOption := c.nchoices - 1
	backOption := c.nchoices

	switch c.cursor {
	case removeOption:
		app.media.ConfigStream(app.selectedStream, ffmpeg.Config{Remove: true})

	case backOption:
		// do nothing

	default:
		app.media.ConfigStream(app.selectedStream, ffmpeg.Config{Codec: c.choices[c.cursor], Remove: false})
	}

	app.selectedStream = nil
	app.page = NewStreamChooser(app.media.Streams())
	return app, nil
}

func (c CodecChooser) View(app App) string {
	var view strings.Builder

	for i, codec := range c.choices {
		cursor := getCursor(c.cursor, i)

		fmt.Fprintf(&view, "%s %s \n", cursor, codec)
	}

	removeCursor := getCursor(c.cursor, c.nchoices-1)
	fmt.Fprintf(&view, "%s Remove \n", removeCursor)

	backCursor := getCursor(c.cursor, c.nchoices)
	fmt.Fprintf(&view, "%s Back \n", backCursor)

	return view.String()
}
