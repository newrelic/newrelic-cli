package yq

// Borrowed from gojq - cli package
// ref: https://github.com/itchyny/gojq/blob/main/cli/error.go

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	exitCodeDefaultErr = iota
)

type EmptyError struct {
	Err error
}

func (*EmptyError) Error() string {
	return ""
}

func (*EmptyError) IsEmptyError() bool {
	return true
}

func (err *EmptyError) ExitCode() int {
	if err, ok := err.Err.(interface{ ExitCode() int }); ok {
		return err.ExitCode()
	}
	return exitCodeDefaultErr
}

type yamlParseError struct {
	fname, contents string
	err             error
}

func (err *yamlParseError) Error() string {
	var line int
	msg := strings.TrimPrefix(
		strings.TrimPrefix(err.err.Error(), "yaml: "),
		"unmarshal errors:\n  ")
	if fmt.Sscanf(msg, "line %d: ", &line); line == 0 {
		return "invalid yaml: " + err.fname
	}
	msg = msg[strings.Index(msg, ": ")+2:]
	if i := strings.IndexByte(msg, '\n'); i >= 0 {
		msg = msg[:i]
	}
	linestr := getLineByLine(err.contents, line)
	return fmt.Sprintf("invalid yaml: %s:%d\n%s  %s",
		err.fname, line, formatLineInfo(linestr, line, 0), msg)
}

func getLineByLine(str string, line int) (linestr string) {
	ss := &stringScanner{str, 0}
	for {
		str, _, ok := ss.next()
		if !ok {
			break
		}
		if line--; line == 0 {
			linestr = str
			break
		}
	}
	if len(linestr) > 64 {
		linestr = trimLastInvalidRune(linestr[:64])
	}
	return
}

func trimLastInvalidRune(s string) string {
	for i := len(s) - 1; i >= 0 && i > len(s)-utf8.UTFMax; i-- {
		if b := s[i]; b < utf8.RuneSelf {
			return s[:i+1]
		} else if utf8.RuneStart(b) {
			if r, _ := utf8.DecodeRuneInString(s[i:]); r == utf8.RuneError {
				return s[:i]
			}
			break
		}
	}
	return s
}

func formatLineInfo(linestr string, line, column int) string {
	l := strconv.Itoa(line)
	return "    " + l + " | " + linestr + "\n" +
		strings.Repeat(" ", len(l)+column) + "       ^"
}

type stringScanner struct {
	str    string
	offset int
}

func (ss *stringScanner) next() (line string, start int, ok bool) {
	if ss.offset == len(ss.str) {
		return
	}
	start, ok = ss.offset, true
	line = ss.str[start:]
	i := indexNewline(line)
	if i < 0 {
		ss.offset = len(ss.str)
		return
	}
	line = line[:i]
	if strings.HasPrefix(ss.str[start+i:], "\r\n") {
		i++
	}
	ss.offset += i + 1
	return
}

// Faster than strings.IndexAny(str, "\r\n").
func indexNewline(str string) (i int) {
	if i = strings.IndexByte(str, '\n'); i >= 0 {
		str = str[:i]
	}
	if j := strings.IndexByte(str, '\r'); j >= 0 {
		i = j
	}
	return
}
