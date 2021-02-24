package ux

import (
	"fmt"
	"time"

	spinnerLib "github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"
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
	// Only start the spinner of the log level is info or below.
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug(msg)
	} else {
		s.Spinner = spinnerLib.New(charSet, interval)
		s.Prefix = indentation
		s.Suffix = fmt.Sprintf(" %s", msg)
		s.Spinner.Start()
	}
}

func (s *Spinner) Stop() {
	// Only stop the spinner of the log level is info or below.
	if log.IsLevelEnabled(log.InfoLevel) {
		s.Spinner.Stop()
		fmt.Println(s.Suffix)
	}
}

func (s *Spinner) Fail(msg string) {
	s.FinalMSG = indentation + crossmark
	s.Suffix = s.Suffix + "failed."
}

func (s *Spinner) Success(msg string) {
	s.FinalMSG = indentation + checkmark
	s.Suffix = s.Suffix + "success."
}
