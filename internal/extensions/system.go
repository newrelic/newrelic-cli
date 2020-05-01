package extensions

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type Manifest struct {
	Command string
}

type Extension struct {
	cmd *exec.Cmd
	args []string
	ctx context.Context
	cancelFunc context.CancelFunc
}

type ConfigOption func(*Extension) error

func WithTimeout(duration time.Duration) ConfigOption {
	return func(e *Extension) error {
		ctx, cancelFunc := context.WithTimeout(e.ctx, duration)
		e.ctx = ctx
		e.cancelFunc = cancelFunc
		return nil
	}
}

func WithArgs(args ...string) ConfigOption {
	return func(e *Extension) error {
		e.args = args
		return nil
	}
}

func New(m *Manifest, opts ...ConfigOption) *Extension {
	ctx, cancelFunc := context.WithCancel(context.Background())
	e := &Extension{
		ctx: ctx,
		cancelFunc: cancelFunc,
	}

	for _, option := range opts {
		option(e)
	}
	
	e.cmd = exec.CommandContext(e.ctx, m.Command, e.args...)

	return e
}

func (e *Extension) Start() {
	go start(e)
}

func start(e *Extension) error {
	defer e.cancelFunc()

	if err := e.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start cmd: %v", err)
	}

	if err := e.cmd.Wait(); err != nil {
		return fmt.Errorf("cmd returned error: %v", err)
	}

	return nil
}