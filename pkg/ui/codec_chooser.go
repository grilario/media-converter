package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grilario/video-converter/pkg/ffmpeg"
)

var (
	choices  = ffmpeg.ListCodecs()
	nchoices = len(choices) + 1
	Remove   = nchoices - 1
	Back     = nchoices
)

type CodecChooser struct {
	cursor int
}

func (c CodecChooser) Update(msg Msg, app App) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "l", "right":
			return c.choose(app)

		case "esc", "h", "left":
			c.cursor = Back
			return c.choose(app)

		case "up", "k":
			if c.cursor > 0 {
				c.cursor--
			}

		case "down", "j":
			if c.cursor < nchoices {
				c.cursor++
			}
		}
	}

	app.currentTab = c
	return app, nil
}

func (c CodecChooser) choose(app App) (tea.Model, tea.Cmd) {
	switch c.cursor {
	case Remove:
		app.media.UpdateStream(app.selectedStream, ffmpeg.Config{Remove: true})

	case Back:
		// do nothing

	default:
		app.media.UpdateStream(app.selectedStream, ffmpeg.Config{Codec: choices[c.cursor], Remove: false})
	}

	app.selectedStream = nil
	app.currentTab = StreamChooser{}
	return app, nil
}

func (c CodecChooser) View(app App) string {
	var s strings.Builder

	for i, codec := range choices {
		cursor := c.getCursor(i)

		fmt.Fprintf(&s, "%s %s \n", cursor, codec)
	}

	removeCursor := c.getCursor(nchoices - 1)
	fmt.Fprintf(&s, "%s Remove \n", removeCursor)

	backCursor := c.getCursor(nchoices)
	fmt.Fprintf(&s, "%s Back \n", backCursor)

	return s.String()
}

func (c CodecChooser) getCursor(current int) string {
	cursor := " "
	if c.cursor == current {
		cursor = ">"
	}

	return cursor
}
