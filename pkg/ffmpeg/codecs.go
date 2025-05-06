package ffmpeg

import (
	_ "embed"
	"encoding/csv"
	"strings"
)

//go:embed languages.csv
var languageCodecsCSV string

// Languages mapped to iso 639-2 key and value in english name
var LanguageCodecs = make(map[string]string)

// populate LanguageCodecs with base in language-codecs.csv
func init() {
	reader := csv.NewReader(strings.NewReader(languageCodecsCSV))
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, record := range records {
		key := record[0]
		value := record[1]

		LanguageCodecs[key] = value
	}
}

type LanguageCodec string

func (c LanguageCodec) String() string {
	return LanguageCodecs[string(c)]
}

type Codec uint8

const (
	COPY Codec = iota
	// Video codecs
	H264
	H265
	VP9
	// Audio codecs
	MP3
	AAC
	AC3
	FLAC
	OPUS
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
	case MP3:
		return "MP3"
	case AAC:
		return "AAC"
	case AC3:
		return "AC3"
	case FLAC:
		return "FLAC"
	case OPUS:
		return "OPUS"
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
	case MP3:
		return "libmp3lame"
	case AAC:
		return "acc"
	case AC3:
		return "ac3"
	case FLAC:
		return "flac"
	case OPUS:
		return "libopus"
	}

	return "unknown"
}

func ListVideoCodecs() []Codec {
	return []Codec{COPY, H264, H265, VP9}
}

func ListAudioCodecs() []Codec {
	return []Codec{COPY, MP3, AAC, AC3, FLAC, OPUS}
}

func ListSubtitlesCodecs() []Codec {
	return []Codec{COPY}
}
