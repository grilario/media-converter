package ffmpeg

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/grilario/video-converter/pkg/ffprobe"
	"github.com/grilario/video-converter/pkg/runner"
)

type Codec uint8

const (
	COPY Codec = iota
	H264
	H265
	VP9
	AV1
)

func (c Codec) String() string {
	switch c {
	case COPY:
		return "COPY"
	case H264:
		return "H264"
	case H265:
		return "H265"
	case VP9:
		return "VP9"
	case AV1:
		return "AV1"
	}

	return "unknown"
}

func (c Codec) encoder() string {
	switch c {
	case COPY:
		return "copy"
	case H264:
		return "libx264"
	case H265:
		return "libx265"
	case VP9:
		return "libvpx-vp9"
	case AV1:
		return "libsvtav1"
	}

	return "unknown"
}

func ListCodecs() []Codec {
	return []Codec{COPY, H264, H265, VP9, AV1}
}

type Stream struct {
	codecType     string
	codecName     string
	codecLongName string
	codecOut      Codec
	remove        bool
}

func (s Stream) Options() (string, string, Codec, bool) {
	return s.codecType, s.codecName, s.codecOut, s.remove
}

type Media struct {
	// path of input that will be handled by ffmpeg
	Input string
	// path of output that will be handled by ffmpeg
	Output string
	// order of array preserves the stream media order of input
	Streams []Stream
}

func NewMedia(probe ffprobe.MediaDetails, outPath string) (Media, error) {
	m := Media{
		Input:  probe.Filepath,
		Output: outPath,
	}

	for _, s := range probe.Streams {
		m.Streams = append(m.Streams, Stream{s.Codec_type, s.Codec_name, s.Codec_long_name, COPY, false})
	}

	return m, nil
}

type Config struct {
	Codec  Codec
	Remove bool
}

func (m *Media) UpdateStream(stream *Stream, config Config) {
	hasStream := false
	for i := range m.Streams {
		if &m.Streams[i] == stream {
			hasStream = true
		}
	}

	if hasStream {
		stream.codecOut = config.Codec
		stream.remove = config.Remove
	}
}

// Convert calls ffmpeg runner with arguments generated based in Streams field,
// it's send progress to channel (1 == 100%) and any errors occurred to channel error
func (m *Media) Convert(progress chan float64, errors chan error, runner runner.Runner) {
	duration, err := ffprobe.MediaDuration(m.Input, runner)
	if err != nil {
		errors <- err
		return
	}

	args := []string{"-y", "-loglevel", "warning", "-progress", "pipe:1", "-nostats", "-i", m.Input}

	for i, stream := range m.Streams {
		if stream.remove {
			continue
		}

		if stream.codecOut == COPY {
			args = append(args, "-map", fmt.Sprintf("0:%v", i), "-codec", "copy")
		} else {
			args = append(args, "-map", fmt.Sprintf("0:%v", i), "-codec", fmt.Sprintf("%v", stream.codecOut.encoder()))
		}

	}

	args = append(args, m.Output)

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
