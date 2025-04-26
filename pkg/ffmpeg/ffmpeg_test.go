package ffmpeg

import (
	"path"
	"testing"

	"github.com/grilario/video-converter/pkg/ffprobe"
	"github.com/grilario/video-converter/pkg/runner"
	"github.com/stretchr/testify/assert"
)

func TestNewMedia(t *testing.T) {
	t.Run("must be create new media", func(t *testing.T) {
		input := ffprobe.MediaDetails{
			Filepath: "simple.mkv",
			Streams:  []ffprobe.Stream{{Index: 0, Codec_type: "video", Codec_name: "h264", Codec_long_name: "h264 mpeg"}},
		}

		media, err := NewMedia(input, "any")
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, input.Filepath, media.Input)
		assert.Equal(t, "any", media.Output)
		assert.Len(t, media.Streams, len(input.Streams))
	})
}

func TestMedia_UpdateStream(t *testing.T) {
	t.Run("must be set out type equal H.265", func(t *testing.T) {
		input := ffprobe.MediaDetails{
			Filepath: "simple.mkv",
			Streams:  []ffprobe.Stream{{Index: 0, Codec_type: "video", Codec_name: "h264", Codec_long_name: "h264 mpeg"}},
		}

		m, err := NewMedia(input, "any")
		if !assert.NoError(t, err) {
			return
		}

		m.UpdateStream(&m.Streams[0], Config{Codec: H265})

		assert.Equal(t, "H.265", m.Streams[0].codecOut.String())
	})
}

func TestMedia_Convert(t *testing.T) {
	t.Run("must be convert video", func(t *testing.T) {
		tempDir := t.TempDir()
		video := "../../tests/video.mp4"
		r, _ := runner.NewRunner()

		d, err := ffprobe.Info(video, r)
		if !assert.NoError(t, err) {
			return
		}

		m, err := NewMedia(d, path.Join(tempDir, "out.mp4"))
		if !assert.NoError(t, err) {
			return
		}

		m.UpdateStream(&m.Streams[0], Config{Codec: H264})

		p := make(chan float64)
		e := make(chan error)
		go m.Convert(p, e, r)

		for {
			select {
			case p := <-p:
				assert.NotNil(t, p)
			case err := <-e:
				assert.NoError(t, err)
				return
			}
		}
	})
}
