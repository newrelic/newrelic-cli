package extensions

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"time"
)

// Manifest describes an extension.
type Manifest struct {
	Command string
}

// Extension represents an extension runner.
type Extension struct {
	CancelFunc context.CancelFunc
	DoneChan   <-chan struct{}
	StdinPipe  io.WriteCloser
	StdoutPipe io.ReadCloser
	StderrPipe io.ReadCloser
	ctx        context.Context
	cmd        *exec.Cmd
	args       []string
	waitErr    error
}

// ConfigOption is a function type for configuring the extension runner.
type ConfigOption func(*Extension) error

// WithTimeout sets a timeout for the command to be run.
func WithTimeout(duration time.Duration) ConfigOption {
	return func(e *Extension) error {
		ctx, cancelFunc := context.WithTimeout(context.Background(), duration)
		e.ctx = ctx
		e.CancelFunc = func() {
			cancelFunc()
		}

		return nil
	}
}

// WithArgs sets the arguments for the command to be run.
func WithArgs(args ...string) ConfigOption {
	return func(e *Extension) error {
		e.args = args
		return nil
	}
}

// New creates a new extension runner.
func New(m *Manifest, opts ...ConfigOption) (*Extension, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	e := &Extension{
		ctx:        ctx,
		CancelFunc: cancelFunc,
	}

	for _, option := range opts {
		if err := option(e); err != nil {
			return nil, err
		}
	}

	e.cmd = exec.CommandContext(e.ctx, m.Command, e.args...)

	stdinPipe, err := e.cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	e.StdinPipe = stdinPipe

	stdoutPipe, err := e.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	e.StdoutPipe = stdoutPipe

	stderrPipe, err := e.cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	e.StderrPipe = stderrPipe

	e.DoneChan = e.ctx.Done()

	return e, nil
}

// Start starts the extension asynchronously.
func (e *Extension) Start() error {
	if err := e.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start cmd: %v", err)
	}

	go func() {
		defer e.CancelFunc()
		e.waitErr = e.cmd.Wait()
	}()

	return nil
}

// Err returns an error if one was encountered.
func (e *Extension) Err() error {
	if e.waitErr != nil {
		return e.waitErr
	}
	if errors.Is(e.ctx.Err(), context.DeadlineExceeded) {
		return e.ctx.Err()
	}

	return nil
}
