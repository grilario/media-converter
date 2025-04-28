package ffprobe

import (
	"encoding/json"
	"strconv"

	"github.com/grilario/video-converter/pkg/runner"
)

type MediaDetails struct {
	Path    string
	Streams []Stream
}

type Stream struct {
	Index         int
	CodecType     string
	CodecName     string
	CodecLongName string

	// in some media the cover image is a stream
	IsAttachment bool
}

// executes ffprobe getting information about file streams
func Info(path string, runner runner.Runner) (MediaDetails, error) {
	args := []string{path, "-show_streams", "-loglevel", "error", "-print_format", "json"}
	buf, err := runner.FFprobe(args)
	if err != nil {
		return MediaDetails{}, err
	}

	details, err := decode(buf)
	if err != nil {
		return MediaDetails{}, err
	}

	media := details.IntoMediaDetails()
	media.Path = path

	return media, nil
}

type detailsJson struct {
	Streams []streamJson
}

type streamJson struct {
	Index         int
	CodecType     string `json:"codec_type"`
	CodecName     string `json:"codec_name"`
	CodecLongName string `json:"codec_long_name"`
	Disposition   dispositionJson
}

type dispositionJson struct {
	AttachmentPic int `json:"attached_pic"`
}

func (d detailsJson) IntoMediaDetails() MediaDetails {
	var streams []Stream
	for _, s := range d.Streams {
		isAttachment := s.Disposition.AttachmentPic == 1

		streams = append(streams, Stream{s.Index, s.CodecType, s.CodecName, s.CodecLongName, isAttachment})
	}

	return MediaDetails{
		Streams: streams,
	}
}

func decode(buf []byte) (detailsJson, error) {
	var data detailsJson
	err := json.Unmarshal(buf, &data)
	if err != nil {
		return detailsJson{}, err
	}

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
