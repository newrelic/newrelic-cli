package execution

import (
	"container/list"
	"io"

	log "github.com/sirupsen/logrus"
)

type ErrorCaptureBuffer struct {
	LastFullLine string
	current      []byte
	AllLines     list.List
	writer       io.Writer
}

func NewErrorCaptureBuffer(w io.Writer) *ErrorCaptureBuffer {
	b := &ErrorCaptureBuffer{
		writer: w,
	}

	return b
}

func (c *ErrorCaptureBuffer) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if b == '\n' {
			s := string(c.current)

			if s != "" {
				log.Debugf(s)
				c.AllLines.PushFront(s)
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

func (c *ErrorCaptureBuffer) Current() string {
	return string(c.current)
}
