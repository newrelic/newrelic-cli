package execution

import (
	log "github.com/sirupsen/logrus"
	"io"
)

type LineCaptureBuffer struct {
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
	for _, b := range p {
		if b == '\n' {
			s := string(c.current)
			c.fullRecipeOutput = append(c.fullRecipeOutput, s)

			if s != "" {
				log.Debugf(s)
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
	return string(c.current)
}

func (c *LineCaptureBuffer) GetFullRecipeOutput() []string {
	return c.fullRecipeOutput
}
