package aimonitoring

import (
	"errors"
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

var promptConfirm = func(msg string) (bool, error) {
	confirmed := false
	err := survey.AskOne(&survey.Confirm{Message: msg, Default: false}, &confirmed)
	return confirmed, err
}

func confirmAction(assumeYes bool, msg string) (bool, error) {
	if assumeYes {
		return true, nil
	}

	confirmed, err := promptConfirm(msg)
	if errors.Is(err, terminal.InterruptErr) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("confirmation required: re-run with --yes to skip this prompt: %w", err)
	}

	return confirmed, nil
}
