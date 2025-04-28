package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grilario/video-converter/pkg/ffmpeg"
	"github.com/grilario/video-converter/pkg/ffprobe"
	"github.com/grilario/video-converter/pkg/runner"
	"github.com/grilario/video-converter/pkg/tui"
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

	w := make(tui.WorkerChannel)

	p := tea.NewProgram(tui.NewApp(&m, w))
	go tui.ConverterWorker(p, w, r)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
