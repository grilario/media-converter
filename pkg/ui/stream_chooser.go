package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type StreamChooser struct {
	cursor int
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
			if c.cursor < len(app.media.Streams) {
				c.cursor++
			}
		}
	}

	app.currentTab = c
	return app, nil
}

func (c StreamChooser) choose(app App) (tea.Model, tea.Cmd) {
	// case choice is confirm jump to tab confirm
	if c.cursor == len(app.media.Streams) {
		app.currentTab = Confirm{}
		return app, nil
	}

	app.selectedStream = &app.media.Streams[c.cursor]
	app.currentTab = CodecChooser{}

	return app, nil
}

func (c StreamChooser) View(app App) string {
	if app.media == nil {
		return ""
	}

	var s strings.Builder

	for i, choice := range app.media.Streams {
		cursor := c.getCursor(i)

		codecType, codecIn, codecOut, willRemoved := choice.Options()
		out := fmt.Sprintf("%s", codecOut)
		if willRemoved {
			out = "Delete"
		}
		fmt.Fprintf(&s, "%s %s (%s)   ->   %s \n", cursor, codecType, codecIn, out)
	}

	cursor := c.getCursor(len(app.media.Streams))
	fmt.Fprintf(&s, "%s Confirm", cursor)

	return s.String()
}

func (c StreamChooser) getCursor(current int) string {
	cursor := " "
	if c.cursor == current {
		cursor = ">"
	}

	return cursor
}
