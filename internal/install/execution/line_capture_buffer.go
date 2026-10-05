package execution

import (
	"io"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

// LineCaptureBuffer is handed to exec.Cmd as Stdout and/or Stderr. When a
// single instance is shared between both (see GoTaskRecipeExecutor.Execute,
// non-silent path), exec.Cmd reads a child process's stdout and stderr pipes
// from two separate goroutines, so Write must be safe for concurrent callers.
type LineCaptureBuffer struct {
	mu               sync.Mutex
	LastFullLine     string
	fullRecipeOutput []string
	current          []byte
	writer           io.Writer
}

func NewLineCaptureBuffer(w io.Writer) *LineCaptureBuffer {
	b := &LineCaptureBuffer{
		writer:           w,
		fullRecipeOutput: []string{},
	}

	return b
}

func (c *LineCaptureBuffer) Write(p []byte) (n int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, b := range p {
		if b == '\n' {
			s := string(c.current)
			c.fullRecipeOutput = append(c.fullRecipeOutput, s)

			if strings.TrimSpace(s) != "" {
				log.Debugf("%s", s)
				c.LastFullLine = s
			}

			c.current = []byte{}
		} else {
			c.current = append(c.current, b)
		}
	}

	if c.writer == nil {
		return 0, nil
	}

	return c.writer.Write(p)
}

func (c *LineCaptureBuffer) Current() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return string(c.current)
}

func (c *LineCaptureBuffer) GetFullRecipeOutput() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.fullRecipeOutput
}
