package command

import (
	"fmt"
	"strings"
)

// ffmpeg -i INPUT.mkv -map 0:1 -c copy OUTPUT.mp4
// select stream of firts input with index 1 and copy to output

var ()

type Input struct {
	Filename string
	Streams  []Stream
}

type Stream struct {
	Name  string
	Codec string

	index      int
	inputIndex int
}

type Output struct {
	Filename string
	Streams  []StreamOut
}

type StreamOut struct {
	ref *Stream
	Stream
}

func (out *Output) AddStream(input *Stream, codec string) {
	stream := StreamOut{
		ref: input,
		Stream: Stream{
			Name:  input.Name,
			Codec: codec,
		},
	}

	out.Streams = append(out.Streams, stream)
}

func (out *Output) RemoveStream() {

}

func (out *Output) GenerateCommand() string {
	var c strings.Builder

	for _, stream := range out.Streams {
		// select input stream to be used from input index and stream index
		fmt.Fprintf(&c, "-map %d:%d ", stream.ref.inputIndex, stream.ref.index)

		if stream.ref.Codec == stream.Codec {
			c.WriteString("-c copy")
		} else {
			fmt.Fprintf(&c, "-c %s ", stream.Codec)
		}
	}

	return c.String()
}
