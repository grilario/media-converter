package ffprobe

import (
	"encoding/json"
	"strconv"

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
func Info(filepath string, runner runner.Runner) (MediaDetails, error) {
	args := []string{filepath, "-show_streams", "-loglevel", "error", "-print_format", "json"}
	buf, err := runner.FFprobe(args)
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

func MediaDuration(filename string, runner runner.Runner) (float64, error) {
	args := []string{filename, "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", "-loglevel", "error"}
	buf, err := runner.FFprobe(args)
	if err != nil {
		return 0, err
	}

	timeStr := string(buf[0:(len(buf) - 1)]) // remove \n from the end of the slice
	time, err := strconv.ParseFloat(timeStr, 64)
	if err != nil {
		return 0, err
	}

	duration := time * 1_000_000 // convert seconds to microsecond
	return duration, nil
}
