package app

import (
	"context"

	"github.com/grilario/video-converter/pkg/ffmpeg"
	"github.com/grilario/video-converter/pkg/runner"
)

type App struct {
	Media    ffmpeg.Media
	Commandc chan WorkerMsg
	Runner   runner.Runner
	Notify   func(progress float64, err error)

	// worker context
	ctx context.Context
}

type WorkerMsg uint8

const (
	StartConversion WorkerMsg = iota
	CancelConversion
)

// start ther worker and wait for the worker resources on cancellations
func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go a.runWorker(ctx, cancel)

	<-ctx.Done()
}

// goroutine that handles media conversion
func (a *App) runWorker(ctx context.Context, cancel context.CancelFunc) {
	for {
		select {
		case cmd := <-a.Commandc:
			switch cmd {
			case CancelConversion:
				cancel()

			case StartConversion:
				go a.convert(ctx)
			}

		case <-ctx.Done():
			return
		}
	}
}

func (a *App) convert(ctx context.Context) {
	progress := make(chan float64)
	err := make(chan error)

	go a.Media.Convert(ctx, progress, err, a.Runner)

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

func New(media ffmpeg.Media, runner runner.Runner) *App {
	commandCh := make(chan WorkerMsg)
	app := &App{
		Media:    media,
		Commandc: commandCh,
		Runner:   runner,
	}

	return app
}
