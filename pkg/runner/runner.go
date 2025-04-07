package runner

import (
	"bytes"
	"os/exec"
)

type Runner struct{}

type FFprobe struct{}

type FFmpeg struct{}

func NewRunner() (Runner, error) {
	return Runner{}, nil
}

func (r *Runner) FFprobe() *FFprobe {
	return &FFprobe{}
}

func (r *Runner) FFmpeg() *FFmpeg {
	return &FFmpeg{}
}

func (r *FFprobe) Exec(args []string) ([]byte, error) {
	cmd := exec.Command("ffprobe", args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return []byte{}, err
	}

	return stdout.Bytes(), nil
}

func (r *FFmpeg) Exec(args []string, progress chan float32) error {
	return nil
}
