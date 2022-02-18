package ux

import (
	"fmt"
	"strings"
	"time"

	spinnerLib "github.com/briandowns/spinner"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
)

type SpinnerProgressIndicator struct {
	*spinnerLib.Spinner
	AssumeYes bool
}

func NewSpinnerProgressIndicator() *SpinnerProgressIndicator {
	s := &SpinnerProgressIndicator{}
	s.Spinner = spinnerLib.New(spinnerLib.CharSets[4], 750*time.Millisecond)
	_ = s.Spinner.Color("green")
	s.Spinner.HideCursor = true
	s.AssumeYes = true
	return s
}

func (s *SpinnerProgressIndicator) Start(msg string) {
	// Suppress spinner output when logging at debug or trace level.
	// Output is garbled when verbose log messages are sent during an active spinner.
	if !config.Logger.IsLevelEnabled(log.DebugLevel) && s.AssumeYes {

		dots := ""
		s.Spinner.PostUpdate = func(s *spinnerLib.Spinner) {
			if dots == ".." {
				dots = ""
			}
			dots += "."
			s.Suffix = fmt.Sprintf(" %s%s", msg, dots)
		}

		s.Spinner.Start() // Start the spinner
	} else {
		c := color.New(color.FgCyan)
		c.Printf("==>")
		x := color.New(color.Bold)
		x.Printf(" %s", msg)
		fmt.Println()
	}
	log.Debug(msg)
}

func (s *SpinnerProgressIndicator) Stop() {
	if !config.Logger.IsLevelEnabled(log.DebugLevel) && s.AssumeYes {
		s.Suffix = ""
		s.Spinner.Stop()
	}
}

func (s *SpinnerProgressIndicator) Fail(msg string) {

	msg = fmt.Sprintf("%v %s\n", IconError, msg)
	if !config.Logger.IsLevelEnabled(log.DebugLevel) && s.AssumeYes {
		s.Suffix = ""
		s.FinalMSG = msg
		s.Spinner.Stop()
	} else {
		fmt.Print(msg)
	}

	printInstallFinalMessage(msg, "Failed", color.BgMagenta)
}

func (s *SpinnerProgressIndicator) Success(msg string) {

	msg = fmt.Sprintf("%v %s\n", IconSuccess, msg)
	if !config.Logger.IsLevelEnabled(log.DebugLevel) && s.AssumeYes {
		s.Suffix = ""
		s.FinalMSG = msg
		// s.FinalMSG = IconSuccess + " Connected to New Relic Platform.\n"
		s.Spinner.Stop()

	} else {
		fmt.Print(msg)
	}

	printInstallFinalMessage(msg, "Installed", color.BgGreen)
}

func printInstallFinalMessage(msg string, printText string, bgColor color.Attribute) {

	if strings.Contains(msg, "Installing") {
		white := color.New(color.FgWhite)
		boldWhite := white.Add(color.Bold)
		background := boldWhite.Add(bgColor)
		fmt.Print("  ")
		background.Print(printText)
		fmt.Println()
	}
}

func (s *SpinnerProgressIndicator) Canceled(msg string) {

	msg = fmt.Sprintf("%v %s\n", IconExclamation, msg)
	if !config.Logger.IsLevelEnabled(log.DebugLevel) && s.AssumeYes {
		s.Suffix = ""
		s.FinalMSG = msg
		// s.FinalMSG = IconSuccess + " Connected to New Relic Platform.\n"
		s.Spinner.Stop()
	} else {
		fmt.Print(msg)
	}
	printInstallFinalMessage(msg, "Cancelled", color.BgBlue)
}
