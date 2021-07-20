package ux

import (
	"fmt"
	"time"

	spinnerLib "github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
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
	isStarted bool
}

func NewSpinner() *Spinner {
	s := Spinner{}
	s.Spinner = spinnerLib.New(charSet, interval)
	s.isStarted = false
	return &s
}

func (s *Spinner) Start(msg string) {
	// Only start the spinner of the log level is info or below.
	if config.Logger.IsLevelEnabled(log.DebugLevel) {
		log.Debug(msg)
	} else {
		s.Spinner = spinnerLib.New(charSet, interval)
		s.Prefix = indentation
		s.Suffix = fmt.Sprintf(" %s", msg)
		s.Spinner.Start()
		s.isStarted = true
	}
}

func (s *Spinner) Stop() {
	// Only stop the spinner of the log level is info or below.
	if s.isStarted {
		s.Spinner.Stop()
		fmt.Println(s.Suffix)
		s.isStarted = false
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
