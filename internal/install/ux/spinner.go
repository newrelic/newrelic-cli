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
	indentation = "  "
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
	// Suppress spinner output when logging at debug or trace level.
	// Output is garbled when verbose log messages are sent during an active spinner.
	if !config.Logger.IsLevelEnabled(log.DebugLevel) {
		s.Spinner = spinnerLib.New(charSet, interval)
		s.Prefix = indentation
		s.Suffix = fmt.Sprintf(" %s", msg)

		fmt.Println()

		s.Spinner.Start()
	}
	log.Debug(msg)
}

func (s *Spinner) Stop() {
	// Suppress stopping the spinner when logging at debug or trace level.
	// See above.
	if !config.Logger.IsLevelEnabled(log.DebugLevel) {
		s.Spinner.Stop()
		fmt.Println(s.Suffix)
		fmt.Println()
	}
	log.Debug(s.Suffix)
}

func (s *Spinner) Fail(msg string) {
	s.FinalMSG = indentation + crossmark
	s.Suffix = s.Suffix + "failed."
}

func (s *Spinner) Success(msg string) {
	s.FinalMSG = indentation + checkmark
	s.Suffix = s.Suffix + "success."
}

func (s *Spinner) Canceled(msg string) {
	s.FinalMSG = indentation + crossmark
	s.Suffix = s.Suffix + "canceled."
}
