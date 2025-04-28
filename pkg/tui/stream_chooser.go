package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grilario/video-converter/pkg/ffmpeg"
)

type StreamChooser struct {
	cursor   int
	choices  []*ffmpeg.Stream
	nchoices int
}

func NewStreamChooser(streams []*ffmpeg.Stream) StreamChooser {
	return StreamChooser{
		cursor:   0,
		choices:  streams,
		nchoices: len(streams), // contains apply choice
	}
}

func (c StreamChooser) Init() tea.Cmd {
	return nil
}

func (c StreamChooser) Update(msg Msg, app App) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "l", "right":
			return c.choose(app)

		case "up", "k":
			if c.cursor > 0 {
				c.cursor--
			}

		case "down", "j":
			if c.cursor < len(c.choices) {
				c.cursor++
			}
		}
	}

	app.page = c
	return app, nil
}

func (c StreamChooser) choose(app App) (tea.Model, tea.Cmd) {
	// case choice is confirm jump to tab confirm
	if c.cursor == c.nchoices {
		app.page = Confirm{}
		return app, nil
	}

	app.selectedStream = c.choices[c.cursor]
	app.page = NewCodecChooser(app.selectedStream)

	return app, nil
}

func (c StreamChooser) View(app App) string {
	if app.media == nil {
		return ""
	}

	var s strings.Builder

	for i, choice := range c.choices {
		cursor := getCursor(c.cursor, i)

		outCodec := "Delete"
		if !choice.ShouldRemoved() {
			name, _ := choice.OutCodec()
			outCodec = fmt.Sprintf("%s", name)
		}

		entryCodec, _ := choice.EntryCodec()

		fmt.Fprintf(&s, "%s %s (%s)   ->   %s \n", cursor, choice.Kind(), entryCodec, outCodec)
	}

	cursor := getCursor(c.cursor, c.nchoices)
	fmt.Fprintf(&s, "%s Confirm", cursor)

	return s.String()
}
