package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunFfprobe(t *testing.T) {
	r, _ := NewRunner()
	runner := r.FFprobe()

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
			got, err := runner.Exec(tt.args)

			if assert.NoError(t, err) {
				assert.NotEmpty(t, got)
			}
		})
	}
}
