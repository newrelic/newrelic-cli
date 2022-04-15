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
	showSpinner bool
}

func NewSpinnerProgressIndicator() *SpinnerProgressIndicator {
	s := &SpinnerProgressIndicator{}
	s.Spinner = spinnerLib.New(spinnerLib.CharSets[4], 750*time.Millisecond)
	_ = s.Spinner.Color("green")
	s.Spinner.HideCursor = true
	s.showSpinner = true
	return s
}

func (s *SpinnerProgressIndicator) ShowSpinner(ss bool) {
	s.showSpinner = ss
}

func (s *SpinnerProgressIndicator) Start(msg string) {
	// Suppress spinner output when logging at debug or trace level.
	// Output is garbled when verbose log messages are sent during an active spinner.
	if !config.Logger.IsLevelEnabled(log.DebugLevel) && s.showSpinner {

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
}

func (s *SpinnerProgressIndicator) Stop() {
	if !config.Logger.IsLevelEnabled(log.DebugLevel) && s.showSpinner {
		s.Suffix = ""
		s.Spinner.Stop()
	}
}

func (s *SpinnerProgressIndicator) Fail(msg string) {

	msg = fmt.Sprintf("%v %s\n", IconError, msg)
	if !config.Logger.IsLevelEnabled(log.DebugLevel) && s.showSpinner {
		s.Suffix = ""
		s.FinalMSG = msg
		s.Spinner.Stop()
	} else {
		fmt.Print(msg)
	}

	printInstallFinalMessage("Failed", color.BgMagenta)
}

func (s *SpinnerProgressIndicator) Success(msg string) {

	msg = fmt.Sprintf("%v %s\n", IconSuccess, msg)
	if !config.Logger.IsLevelEnabled(log.DebugLevel) && s.showSpinner {
		s.Suffix = ""
		s.FinalMSG = msg
		s.Spinner.Stop()

	} else {
		fmt.Print(msg)
	}

	if strings.Contains(msg, "Installing") {
		printInstallFinalMessage("Installed", color.BgGreen)
	} else {
		printInstallFinalMessage("Connected", color.BgGreen)
	}
}

func printInstallFinalMessage(printText string, bgColor color.Attribute) {

	white := color.New(color.FgWhite)
	boldWhite := white.Add(color.Bold)
	background := boldWhite.Add(bgColor)
	fmt.Print("  ")
	background.Print(fmt.Sprintf(" %s ", printText))
	fmt.Println()
}

func (s *SpinnerProgressIndicator) Canceled(msg string) {

	msg = fmt.Sprintf("%v %s\n", IconExclamation, msg)
	if !config.Logger.IsLevelEnabled(log.DebugLevel) && s.showSpinner {
		s.Suffix = ""
		s.FinalMSG = msg
		s.Spinner.Stop()
	} else {
		fmt.Print(msg)
	}
	printInstallFinalMessage("Cancelled", color.BgBlue)
}
