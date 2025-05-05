package ffmpeg

import (
	"context"
	"path"
	"testing"

	"github.com/grilario/video-converter/pkg/ffprobe"
	"github.com/grilario/video-converter/pkg/runner"
	"github.com/stretchr/testify/assert"
)

func TestNewMedia(t *testing.T) {
	t.Run("must be create new media", func(t *testing.T) {
		input := ffprobe.MediaDetails{
			Path:    "simple.mkv",
			Streams: []ffprobe.Stream{{Index: 0, CodecType: "video", CodecName: "h264", CodecLongName: "h264 mpeg"}},
		}

		media, err := NewMedia(input, "any")
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, input.Path, media.Input)
		assert.Equal(t, "any", media.Output)
		assert.Len(t, media.Streams(), len(input.Streams))
	})
}

func TestMedia_UpdateStream(t *testing.T) {
	t.Run("must be set out type equal H265", func(t *testing.T) {
		input := ffprobe.MediaDetails{
			Path:    "simple.mkv",
			Streams: []ffprobe.Stream{{Index: 0, CodecType: "video", CodecName: "h264", CodecLongName: "h264 mpeg"}},
		}

		m, err := NewMedia(input, "any")
		if !assert.NoError(t, err) {
			return
		}

		m.ConfigStream(&m.streams[0], Config{Codec: H265})

		assert.Equal(t, "H265", m.streams[0].outCodec.String())
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

		m.ConfigStream(&m.streams[0], Config{Codec: H264})

		p := make(chan float64)
		e := make(chan error)
		go m.Convert(context.Background(), p, e, r)

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

func TestLanguageCodecs_Init(t *testing.T) {
	t.Run("asserts length of map language codecs is 184", func(t *testing.T) {
		assert.Len(t, LanguageCodecs, 184)
	})

	t.Run("LanguageCodecs should contains most used codes", func(t *testing.T) {
		assert.Equal(t, "Portuguese", LanguageCodecs["por"])
		assert.Equal(t, "English", LanguageCodecs["eng"])
		assert.Equal(t, "Spanish; Castilian", LanguageCodecs["spa"])
	})
}
