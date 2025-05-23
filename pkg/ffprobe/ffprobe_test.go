package ffprobe

import (
	"testing"

	"github.com/grilario/video-converter/pkg/runner"
	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	r, _ := runner.NewRunner()
	videoFile := "../../tests/video.mp4"

	type args struct {
		filepath string
		runner   runner.Runner
	}
	tests := []struct {
		name    string
		args    args
		want    MediaDetails
		wantErr bool
	}{
		{
			"must return streams information about media",
			args{videoFile, r},
			MediaDetails{
				Path:    videoFile,
				Streams: []Stream{{0, "video", "h264", "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10", Tags{Language: "und"}, false}},
			},
			false,
		},
		{
			"must return an error",
			args{"unknown", r},
			MediaDetails{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Info(tt.args.filepath, tt.args.runner)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			if assert.NoError(t, err) {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func TestMediaDuration(t *testing.T) {
	t.Run("must be return media duration", func(t *testing.T) {
		r, _ := runner.NewRunner()

		duration, err := MediaDuration("../../tests/video.mp4", r)
		if assert.NoError(t, err) {
			assert.Equal(t, duration, 3067000.0)
		}
	})
}
