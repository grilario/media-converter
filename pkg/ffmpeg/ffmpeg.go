package ffmpeg

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/grilario/video-converter/pkg/ffprobe"
	"github.com/grilario/video-converter/pkg/runner"
)

type Media struct {
	InFilename  string
	OutFilename string
	// order of array streams preserves the stream media order
	Streams []*Stream
}

type Stream struct {
	codec_in_type  string
	codec_in_name  string
	codec_out_type string
	remove         bool
}

func NewMedia(input ffprobe.MediaDetails, outPath string) (Media, error) {
	m := Media{
		InFilename:  input.Filepath,
		OutFilename: outPath,
	}

	for _, stream := range input.Streams {
		m.Streams = append(m.Streams, &Stream{stream.Codec_type, stream.Codec_name, stream.Codec_type, false})
	}

	return m, nil
}

type Config struct {
	codec_type string
	remove     bool
}

func (m *Media) UpdateStream(index int, config Config) {
	stream := m.Streams[index]
	stream.codec_out_type = config.codec_type
	stream.remove = config.remove
}

// Convert calls ffmpeg runner with arguments generated based in Streams field,
// it's send progress to channel (1 == 100%) and any errors occurred to channel error
func (m *Media) Convert(progress chan float64, errors chan error, runner runner.Runner) {
	duration, err := ffprobe.MediaDuration(m.InFilename, runner)
	if err != nil {
		errors <- err
		return
	}

	args := []string{"-y", "-loglevel", "warning", "-progress", "pipe:1", "-nostats", "-i", m.InFilename}

	for i, stream := range m.Streams {
		if stream.remove {
			continue
		}

		if stream.codec_in_type == stream.codec_out_type {
			args = append(args, "-map", fmt.Sprintf("0:%v", i), "-codec", "copy")
		} else {
			args = append(args, "-map", fmt.Sprintf("0:%v", i), "-codec", fmt.Sprintf("%v", stream.codec_out_type))
		}

	}

	args = append(args, m.OutFilename)

	stdout := make(chan []byte)
	ffErrs := make(chan error)

	go runner.FFmpeg(args, stdout, ffErrs)

	for {
		select {
		case stdout := <-stdout:
			line := string(stdout)

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := parts[0]
			value := parts[1]

			if key == "out_time_ms" {
				partialTime, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					errors <- err
					return
				}

				progress <- math.Abs(float64(partialTime)) / duration
			}

			if key == "progress" && value == "end" {
				progress <- 1
			}
		case err := <-ffErrs:
			errors <- err
			return
		}
	}
}
