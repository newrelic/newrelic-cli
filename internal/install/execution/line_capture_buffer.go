package execution

import (
	"container/list"
	log "github.com/sirupsen/logrus"
	"io"
)

type LineCaptureBuffer struct {
	LastFullLine string
	current      []byte
	multiple     list.List
	writer       io.Writer
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

			if s != "" {
				log.Debugf(s)
				c.multiple.PushFront(s)
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
