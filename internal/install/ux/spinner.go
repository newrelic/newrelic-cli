package ux

import (
	"fmt"
	"time"

	spinnerLib "github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
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

func (s *Spinner) Start(recipe types.Recipe) {
	msg := fmt.Sprintf("Installing %s", recipe.Name)

	// Only start the spinner of the log level is info or below.
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug(msg)
	} else {
		s.Spinner = spinnerLib.New(charSet, interval)
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

func (s *Spinner) Fail(recipe types.Recipe) {
	s.FinalMSG = crossmark
}

func (s *Spinner) Success(recipe types.Recipe) {
	s.FinalMSG = checkmark
}
