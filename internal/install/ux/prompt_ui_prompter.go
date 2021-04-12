package ux

import (
	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type PromptUIPrompter struct{}

func NewPromptUIPrompter() *PromptUIPrompter {
	return &PromptUIPrompter{}
}

func (p *PromptUIPrompter) PromptYesNo(msg string) (bool, error) {

	yes := false
	prompt := &survey.Confirm{
		Default: true,
		Message: msg,
	}

	err := survey.AskOne(prompt, &yes)
	if err != nil {
		return false, err
	}

	return yes, nil
}

func (p *PromptUIPrompter) MultiSelect(msg string, options []string) ([]string, error) {
	defaults := utils.MakeRange(0, len(options)-1)
	selected := []string{}
	prompt := &survey.MultiSelect{
		Message: msg,
		Options: options,
		Default: defaults,
	}

	err := survey.AskOne(prompt, &selected)
	if err != nil {
		if err == terminal.InterruptErr {
			return nil, types.ErrInterrupt
		}

		if err.Error() == terminal.InterruptErr.Error() {
			return nil, types.ErrInterrupt
		}

		return nil, err
	}

	return selected, nil
}
