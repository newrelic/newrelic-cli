package execution

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type LineCaptureBuffer struct {
	LastFullLine     string
	fullRecipeOutput []string
	current          []byte
	cliOutput        *os.File
	writer           io.Writer
}

func NewLineCaptureBufferMultiWriters(out io.Writer, err io.Writer) *LineCaptureBuffer {
	b := &LineCaptureBuffer{
		writer: io.MultiWriter(out, err),
	}

	return b
}

func NewLineCaptureBuffer(w io.Writer) *LineCaptureBuffer {
	b := &LineCaptureBuffer{
		writer: w,
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
