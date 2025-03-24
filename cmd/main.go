package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grilario/video-converter/pkg/command"
)

// menus
const (
	OutputFormat = iota
	SelectTrack
	SelectTrackCodec
)

type model struct {
	input         command.Input
	output        command.Output
	menu          int32
	selectedTrack command.StreamOut
	cursor        int
	choices       int
}

func initialModel() model {
	return model{
		input:   command.Input{},
		menu:    OutputFormat,
		choices: 3,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "backspace", "esc":
			if m.menu != OutputFormat {
				m.menu -= 1
			}

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < m.choices-1 {
				m.cursor++
			}
		}
	}

	switch m.menu {
	case OutputFormat:
		return updateSelectedOutput(msg, m)

	case SelectTrack:
		return updateSelectedTrack(msg, m)

	case SelectTrackCodec:
		return updateSelectedCodec(msg, m)
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString("Qual formato do arquivo de saída?\n")

	choices := []string{"MKV", "MP4", "WEBM"}

	for i, choice := range choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		fmt.Fprintf(&s, "%s %s\n", cursor, choice)
	}

	return s.String()
}

func updateSelectedOutput(msg tea.Msg, m model) (tea.Model, tea.Cmd) {

	return m, nil
}

func updateSelectedTrack(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	return m, nil
}

func updateSelectedCodec(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	return m, nil
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
