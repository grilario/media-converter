package app

import (
	"github.com/grilario/video-converter/pkg/ffmpeg"
	"github.com/grilario/video-converter/pkg/runner"
)

type App struct {
	Media    ffmpeg.Media
	Commandc chan WorkerMsg
	Runner   runner.Runner
	Notify   func(progress float64, err error)
}

type WorkerMsg uint8

const (
	StartConversion WorkerMsg = iota
)

func (a *App) converterWorker() {
	for {
		cmd := <-a.Commandc
		if cmd != StartConversion {
			continue // current implementation only start conversion
		}

		progress := make(chan float64)
		err := make(chan error)

		go a.Media.Convert(progress, err, a.Runner)

		for {
			select {
			case progress := <-progress:
				if a.Notify != nil {
					a.Notify(progress, nil)
				}

			case err := <-err:
				if a.Notify != nil {
					a.Notify(1.0, err)
				}
			}
		}
	}
}

func New(media ffmpeg.Media, runner runner.Runner) *App {
	commandCh := make(chan WorkerMsg)
	app := &App{
		Media:    media,
		Commandc: commandCh,
		Runner:   runner,
	}

	go app.converterWorker()

	return app
}
