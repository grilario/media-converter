package ffprobe

import (
	"encoding/json"

	"github.com/grilario/video-converter/pkg/runner"
)

type Executor struct{}

type MediaDetails struct {
	Filepath string
	Streams  []Stream
}

type Stream struct {
	Index           int
	Codec_type      string
	Codec_name      string
	Codec_long_name string
}

// executes ffprobe getting information about file streams
func Info(filepath string, runner *runner.Runner) (MediaDetails, error) {
	args := []string{filepath, "-show_streams", "-loglevel", "error", "-print_format", "json"}
	buf, err := runner.FFprobe().Exec(args)
	if err != nil {
		return MediaDetails{}, err
	}

	var data MediaDetails
	err = json.Unmarshal([]byte(buf), &data)
	if err != nil {
		return MediaDetails{}, err
	}

	data.Filepath = filepath

	return data, nil
}
