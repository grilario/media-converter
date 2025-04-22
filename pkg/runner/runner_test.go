package runner

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunFfprobe(t *testing.T) {
	r, _ := NewRunner()

	tests := []struct {
		name string
		args []string
	}{
		{
			"should be return array bytes not empty",
			[]string{"../../tests/video.mp4", "-show_streams", "-loglevel", "error", "-print_format", "json"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FFprobe(tt.args)

			if assert.NoError(t, err) {
				assert.NotEmpty(t, got)
			}
		})
	}
}

func TestDefaultRunner_FFmpeg(t *testing.T) {
	r, _ := NewRunner()
	tempDir := t.TempDir()
	out := path.Join(tempDir, "out.mkv")

	stdout := make(chan []byte)
	e := make(chan error)
	go r.FFmpeg([]string{"-y", "-loglevel", "error", "-progress", "pipe:1", "-i", "../../tests/video.mp4", "-map", "0:0", "-c", "copy", out}, stdout, e)

	for {
		select {
		case stdout := <-stdout:
			assert.NotNil(t, stdout)
		case err := <-e:
			assert.NoError(t, err)
			return
		}
	}
}
