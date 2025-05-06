package ffmpeg

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/grilario/video-converter/pkg/ffprobe"
	"github.com/grilario/video-converter/pkg/runner"
)

type Media struct {
	// path of input that will be handled by ffmpeg
	Input string
	// path of output that will be handled by ffmpeg
	Output string
	// order of array preserves the stream media order of input
	streams []Stream
}

func NewMedia(probe ffprobe.MediaDetails, outPath string) (Media, error) {
	m := Media{
		Input:  probe.Path,
		Output: outPath,
	}

	for _, s := range probe.Streams {
		stream := Stream{
			Tags:          Tags{s.Tags.Title, LanguageCodec(s.Tags.Language)},
			kind:          StreamKind(s.CodecType),
			codecName:     s.CodecName,
			codecLongName: s.CodecLongName,
			outCodec:      COPY,
			isAttchment:   s.IsAttachment,
			remove:        false,
		}

		m.streams = append(m.streams, stream)
	}

	return m, nil
}

func (m *Media) Streams() []*Stream {
	streams := []*Stream{}

	for i, stream := range m.streams {
		// remove streams that are attchments
		if stream.isAttchment {
			continue
		}

		// reference to original stream not copy of for range, it's is necessary to cofig method
		streams = append(streams, &m.streams[i])
	}

	return streams
}

type Config struct {
	Codec  Codec
	Remove bool
}

func (m *Media) ConfigStream(stream *Stream, config Config) {
	hasStream := false
	for i := range m.streams {
		if &m.streams[i] == stream {
			hasStream = true
		}
	}

	if hasStream {
		stream.outCodec = config.Codec
		stream.remove = config.Remove
	}
}

// Convert calls ffmpeg runner with arguments generated based in Streams field,
// it's send progress to channel (1 == 100%) and any errors occurred to channel error
func (m *Media) Convert(ctx context.Context, progress chan float64, errors chan error, runner runner.Runner) {
	duration, err := ffprobe.MediaDuration(m.Input, runner)
	if err != nil {
		errors <- err
		return
	}

	args := []string{"-y", "-loglevel", "warning", "-progress", "pipe:1", "-nostats", "-i", m.Input, "-dn"}
	mappers := []string{}
	codecs := []string{}

	for i, stream := range m.streams {
		if stream.remove || stream.isAttchment {
			continue
		}

		if stream.outCodec == COPY {
			mappers = append(mappers, "-map", fmt.Sprintf("0:%v", i))
			codecs = append(codecs, fmt.Sprintf("-c:%v", i), "copy")
		} else {
			mappers = append(mappers, "-map", fmt.Sprintf("0:%v", i))
			codecs = append(codecs, fmt.Sprintf("-c:%v", i), fmt.Sprintf("%v", stream.outCodec.encoder()))
		}
	}

	args = append(args, mappers...)
	args = append(args, codecs...)
	args = append(args, m.Output)

	stdout := make(chan []byte)
	ffErrs := make(chan error)

	go runner.FFmpeg(ctx, args, stdout, ffErrs)

	// capture progress stats
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

type StreamKind string

var (
	VideoStream    StreamKind = "video"
	AudioStream    StreamKind = "audio"
	SubtitleStream StreamKind = "subtitle"
)

func (k StreamKind) String() string {
	return string(k)
}

type Stream struct {
	Tags Tags

	kind          StreamKind
	codecName     string
	codecLongName string
	outCodec      Codec
	isAttchment   bool
	remove        bool
}

// Return name, description of entry codec
func (s Stream) EntryCodec() (string, string) {
	return s.codecName, ""
}

// Return name, description of current output codec
func (s Stream) OutCodec() (string, string) {
	return s.outCodec.String(), ""
}

func (s Stream) Kind() StreamKind {
	return s.kind
}

func (s Stream) ShouldRemoved() bool {
	return s.remove
}

type Tags struct {
	Title    string
	Language LanguageCodec
}
