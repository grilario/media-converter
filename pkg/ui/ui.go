package ui

import (
	// "fmt"
	// "strings"
	//
	tea "github.com/charmbracelet/bubbletea"
	"github.com/grilario/video-converter/pkg/ffmpeg"
	"github.com/grilario/video-converter/pkg/runner"
)

type Msg any

type Model interface {
	Update(Msg, App) (tea.Model, tea.Cmd)
	View(App) string
}

type App struct {
	media          *ffmpeg.Media
	currentTab     Model
	selectedStream *ffmpeg.Stream
	workerChannel  WorkerChannel
}

func NewApp(media *ffmpeg.Media, channel WorkerChannel) App {
	return App{
		media:          media,
		currentTab:     StreamChooser{cursor: 0},
		selectedStream: nil,
		workerChannel:  channel,
	}
}

func (a App) Init() tea.Cmd {
	return nil
}

func (app App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return app, tea.Quit
		}
	}

	newApp, cmd := app.currentTab.Update(msg, app)

	return newApp, cmd
}

func (app App) View() string {
	return app.currentTab.View(app)
}

type WorkerChannel chan ffmpeg.Media
type ProgressMsg struct {
	progress float64
	error    error
}

func ConverterWorker(program *tea.Program, msg WorkerChannel, runner runner.Runner) {
	media := <-msg
	close(msg)

	p := make(chan float64)
	e := make(chan error)

	go media.Convert(p, e, runner)

	for {
		select {
		case p := <-p:
			program.Send(ProgressMsg{progress: p, error: nil})

		case err := <-e:
			program.Send(ProgressMsg{progress: 1, error: err})
		}
	}
}
