package runner

import (
	"context"
	"os/exec"
	"path"
	"testing"
	"time"

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
	ctx := context.Background()
	go r.FFmpeg(ctx, []string{"-y", "-loglevel", "error", "-progress", "pipe:1", "-i", "../../tests/video.mp4", "-map", "0:0", "-c", "copy", out}, stdout, e)

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

func TestDefaultRunner_FFmpegCancel(t *testing.T) {
	r, _ := NewRunner()
	tempDir := t.TempDir()
	out := path.Join(tempDir, "out.mkv")

	stdout := make(chan []byte)
	e := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	go r.FFmpeg(ctx, []string{"-y", "-loglevel", "error", "-stream_loop", "-1", "-i", "../../tests/video.mp4", out}, stdout, e)

	for {
		select {
		case stdout := <-stdout:
			assert.NotNil(t, stdout)
		case err := <-e:
			var exitErr *exec.ExitError
			assert.ErrorAsf(t, err, &exitErr, "signal: killed")
			return
		}
	}
}
