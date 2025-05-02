package runner

//go:generate mockgen -destination runner_mock.go -package runner . Runner

import (
	"bufio"
	"bytes"
	"context"
	"os/exec"
	"strings"
)

type Error struct {
	Args string
	Msg  string
	Err  error
}

func (e *Error) Error() string {
	return e.Err.Error() + "\nArgs: " + e.Args + "\nFFmpeg msg: " + e.Msg
}

func (e *Error) Unwrap() error { return e.Err }

type Runner interface {
	FFmpeg(ctx context.Context, args []string, stdout chan []byte, error chan error)
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
		return []byte{}, &Error{Args: strings.Join(args, ", "), Msg: stderr.String(), Err: err}
	}

	return stdout.Bytes(), nil
}

// FFmpeg runs ffmpeg with arguments provided, and send to channel stdout line by line of
// ffmpeg stdout process or any error occurred to channel errors.
func (r *DefaultRunner) FFmpeg(ctx context.Context, args []string, stdout chan []byte, errors chan error) {
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

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
		errors <- &Error{Args: strings.Join(args, ", "), Msg: stderr.String(), Err: err}
		return
	}

	errors <- nil
}
