package ux

import (
	"fmt"
	"time"

	spinnerLib "github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
)

type SpinnerProgressIndicator struct {
	*spinnerLib.Spinner
}

func NewSpinnerProgressIndicator() *SpinnerProgressIndicator {
	s := &SpinnerProgressIndicator{}
	s.Spinner = spinnerLib.New(spinnerLib.CharSets[4], 750*time.Millisecond)
	s.Spinner.Color("green")
	s.Spinner.HideCursor = true
	return s
}

func (s *SpinnerProgressIndicator) Start(msg string) {
	// Suppress spinner output when logging at debug or trace level.
	// Output is garbled when verbose log messages are sent during an active spinner.
	if !config.Logger.IsLevelEnabled(log.DebugLevel) {

		dots := ""
		s.Spinner.PostUpdate = func(s *spinnerLib.Spinner) {
			if dots == ".." {
				dots = ""
			}
			dots += "."
			s.Suffix = fmt.Sprintf(" %s%s", msg, dots)
		}

		s.Spinner.Start() // Start the spinner
	}
	log.Debug(msg)
}

func (s *SpinnerProgressIndicator) Stop() {
	if !config.Logger.IsLevelEnabled(log.DebugLevel) {
		s.Suffix = ""
		s.Spinner.Stop()
	}
}

func (s *SpinnerProgressIndicator) Fail(msg string) {
	s.FinalMSG = indentation + crossmark
	s.Suffix = s.Suffix + "incomplete."

	if !config.Logger.IsLevelEnabled(log.DebugLevel) {
		s.Suffix = ""
		s.FinalMSG = fmt.Sprintf("%v %s\n", IconError, msg)
		// s.FinalMSG = IconError + " Connected to New Relic Platform.\n"
		s.Spinner.Stop()
	}
}

func (s *SpinnerProgressIndicator) Success(msg string) {

	if !config.Logger.IsLevelEnabled(log.DebugLevel) {
		s.Suffix = ""
		s.FinalMSG = fmt.Sprintf("%v %s\n", IconSuccess, msg)
		// s.FinalMSG = IconSuccess + " Connected to New Relic Platform.\n"
		s.Spinner.Stop()
	}
}

func (s *SpinnerProgressIndicator) Canceled(msg string) {

	if !config.Logger.IsLevelEnabled(log.DebugLevel) {
		s.Suffix = ""
		s.FinalMSG = fmt.Sprintf("%v %s\n", IconExclamation, msg)
		// s.FinalMSG = IconSuccess + " Connected to New Relic Platform.\n"
		s.Spinner.Stop()
	}
}
