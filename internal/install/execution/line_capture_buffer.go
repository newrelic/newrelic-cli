package execution

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type LineCaptureBuffer struct {
	LastFullLine     string
	current          []byte
	cliOutput        *os.File
	writer           io.Writer
	fullRecipeOutput []string
}

func NewLineCaptureToFileBuffer(w io.Writer, outputFile *os.File) *LineCaptureBuffer {
	b := &LineCaptureBuffer{
		cliOutput:        outputFile,
		writer:           w,
		fullRecipeOutput: []string{},
	}

	return b
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

			if s != "" {
				log.Debugf(s)
				c.LastFullLine = s
			}

			c.fullRecipeOutput = append(c.fullRecipeOutput, s)
			if nil != c.cliOutput {
				_, err := c.cliOutput.WriteString(s + "\n")
				if nil != err {
					log.Debugf("Couldn't write to cli output file: %e", err)
				}
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
