package ux

import (
	"fmt"
	"time"

	spinnerLib "github.com/briandowns/spinner"
)

const (
	interval  = 100 * time.Millisecond
	checkmark = "\u2705"
	boom      = "\u1F4A5"
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
	s.Suffix = msg
	s.Spinner.Start()
}

func (s *Spinner) Stop() {
	s.Spinner.Stop()
	fmt.Println(s.Suffix)
}

func (s *Spinner) Fail() {
	s.FinalMSG = boom
}

func (s *Spinner) Success() {
	s.FinalMSG = checkmark
}
