package runner

//go:generate mockgen -destination runner_mock.go -package runner . Runner

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
)

type Runner interface {
	FFmpeg(args []string, stdout chan []byte, error chan error)
	FFprobe(args []string) ([]byte, error)
}

type DefaultRunner struct{}

func NewRunner() (Runner, error) {
	return &DefaultRunner{}, nil
}

func (r *DefaultRunner) FFprobe(args []string) ([]byte, error) {
	cmd := exec.Command("ffprobe", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return []byte{}, fmt.Errorf("%v \nArguments: %#v \nError: %v", err, args, stderr.String())
	}

	return stdout.Bytes(), nil
}

// FFmpeg runs ffmpeg with arguments provided, and send to channel stdout line by line of
// ffmpeg stdout process or any error occurred to channel errors.
func (r *DefaultRunner) FFmpeg(args []string, stdout chan []byte, errors chan error) {
	cmd := exec.Command("ffmpeg", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.StdoutPipe()
	if err != nil {
		errors <- err
		return
	}

	err = cmd.Start()
	if err != nil {
		errors <- err
		return
	}

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		stdout <- scanner.Bytes()
	}

	if err := scanner.Err(); err != nil {
		errors <- err
		return
	}

	if err := cmd.Wait(); err != nil {
		errors <- fmt.Errorf("%v \nArguments: %#v \nError: %v", err, args, stderr.String())
		return
	}

	errors <- nil
}
