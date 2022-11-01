package yq

// Borrowed from gojq - cli package
// ref: https://github.com/itchyny/gojq/blob/main/cli/inputs.go

import (
	"bytes"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/itchyny/gojq"
)

type inputReader struct {
	io.Reader
	file *os.File
	buf  *bytes.Buffer
}

func newInputReader(r io.Reader) *inputReader {
	if r, ok := r.(*os.File); ok {
		if _, err := r.Seek(0, io.SeekCurrent); err == nil {
			return &inputReader{r, r, nil}
		}
	}
	var buf bytes.Buffer // do not use strings.Builder because we need to Reset
	return &inputReader{io.TeeReader(r, &buf), nil, &buf}
}

func (ir *inputReader) getContents(offset *int64, line *int) string {
	if buf := ir.buf; buf != nil {
		return buf.String()
	}
	if current, err := ir.file.Seek(0, io.SeekCurrent); err == nil {
		defer func() {
			if _, err := ir.file.Seek(current, io.SeekStart); err != nil {
				log.Fatalln("error reading yaml input from stdin")
			}
		}()
	}
	if _, err := ir.file.Seek(0, io.SeekStart); err != nil {
		log.Fatalln("error reading yaml input from stdin")
	}
	const bufSize = 16 * 1024
	var buf bytes.Buffer // do not use strings.Builder because we need to Reset
	if offset != nil && *offset > bufSize {
		buf.Grow(bufSize)
		for *offset > bufSize {
			n, err := io.Copy(&buf, io.LimitReader(ir.file, bufSize))
			*offset -= n
			*line += bytes.Count(buf.Bytes(), []byte{'\n'})
			buf.Reset()
			if err != nil || n == 0 {
				break
			}
		}
	}
	var r io.Reader
	if offset == nil {
		r = ir.file
	} else {
		r = io.LimitReader(ir.file, bufSize*2)
	}
	_, err := io.Copy(&buf, r)
	if err != nil {
		log.Fatalln("error copying input into buffer")
	}
	return buf.String()
}

type InputIter interface {
	gojq.Iter
	io.Closer
	Name() string
}

type yamlInputIter struct {
	dec   *yaml.Decoder
	ir    *inputReader
	fname string
	err   error
}

func NewYAMLInputIter(r io.Reader, fname string) InputIter {
	ir := newInputReader(r)
	dec := yaml.NewDecoder(ir)
	return &yamlInputIter{dec: dec, ir: ir, fname: fname}
}

func (i *yamlInputIter) Next() (interface{}, bool) {
	if i.err != nil {
		return nil, false
	}
	var v interface{}
	if err := i.dec.Decode(&v); err != nil {
		if err == io.EOF {
			i.err = err
			return nil, false
		}
		i.err = &yamlParseError{i.fname, i.ir.getContents(nil, nil), err}
		return i.err, true
	}
	return normalizeYAML(v), true
}

func (i *yamlInputIter) Close() error {
	i.err = io.EOF
	return nil
}

func (i *yamlInputIter) Name() string {
	return i.fname
}
