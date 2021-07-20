package ux

import (
	"fmt"
	"time"

	spinnerLib "github.com/briandowns/spinner"
)

const (
	interval    = 100 * time.Millisecond
	checkmark   = "\u2705"
	crossmark   = "\u274C"
	indentation = "    "
)

var (
	charSet = spinnerLib.CharSets[14]
)

type Spinner struct {
	*spinnerLib.Spinner
}

func NewSpinner() *Spinner {
	s := Spinner{}
	s.Spinner = spinnerLib.New(charSet, interval)
	return &s
}

func (s *Spinner) Start(msg string) {
	s.Spinner = spinnerLib.New(charSet, interval)
	s.Prefix = indentation
	s.Suffix = fmt.Sprintf(" %s", msg)
	s.Spinner.Start()
}

func (s *Spinner) Stop() {
	s.Spinner.Stop()
	fmt.Println(s.Suffix)
}

func (s *Spinner) Fail(msg string) {
	s.FinalMSG = indentation + crossmark
	s.Suffix = s.Suffix + "failed."
}

func (s *Spinner) Success(msg string) {
	s.FinalMSG = indentation + checkmark
	s.Suffix = s.Suffix + "success."
}
