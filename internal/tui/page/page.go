package page

import "github.com/grilario/video-converter/pkg/ffmpeg"

type PageID string

type PageChangeMsg struct {
	ID PageID
}

type SelectedStreamChangeMsg struct {
	SelectedStream *ffmpeg.Stream
}

type ProgressMsg struct {
	Progress float64
	Error    error
}
