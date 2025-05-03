package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grilario/video-converter/internal/app"
	"github.com/grilario/video-converter/internal/tui"
	"github.com/grilario/video-converter/internal/tui/page"
	"github.com/grilario/video-converter/pkg/ffmpeg"
	"github.com/grilario/video-converter/pkg/ffprobe"
	"github.com/grilario/video-converter/pkg/runner"
)

var input = flag.String("i", "", "Media input")
var output = flag.String("o", "", "Media Ouput")

func main() {
	flag.Parse()

	if len(*input) == 0 || len(*output) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	r, _ := runner.NewRunner()

	d, err := ffprobe.Info(*input, r)
	if err != nil {
		panic(err)
	}

	m, err := ffmpeg.NewMedia(d, *output)
	if err != nil {
		panic(err)
	}

	app := app.New(m, r)
	p := tea.NewProgram(tui.New(app))

	app.Notify = func(progress float64, err error) {
		p.Send(page.ProgressMsg{Progress: progress, Error: err})
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		app.Run()
	}()

	go func() {
		defer wg.Done()

		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	}()

	wg.Wait()
}
