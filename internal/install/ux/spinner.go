package ux

import (
	"fmt"
	"time"

	spinnerLib "github.com/briandowns/spinner"
)

const (
	interval  = 100 * time.Millisecond
	checkmark = "\u2705"
	crossmark = "\u274C"
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
	s.Suffix = fmt.Sprintf(" %s", msg)
	s.Spinner.Start()
}

func (s *Spinner) Stop() {
	s.Spinner.Stop()
	fmt.Println(s.Suffix)
}

func (s *Spinner) Fail() {
	s.FinalMSG = crossmark
}

func (s *Spinner) Success() {
	s.FinalMSG = checkmark
}
