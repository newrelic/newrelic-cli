package ux

import (
	"fmt"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"
	"gopkg.in/AlecAivazis/survey.v1/terminal"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type PromptUIPrompter struct{}

func NewPromptUIPrompter() *PromptUIPrompter {
	return &PromptUIPrompter{}
}

func (p *PromptUIPrompter) PromptYesNo(msg string) (bool, error) {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . | bold }} [Y/n] ",
		Valid:   "{{ . | bold }} [Y/n] ",
		Invalid: "{{ . | bold }} [Y/n] ",
		Success: "{{ . | faint }} ",
	}

	prompt := promptui.Prompt{
		Default:   "Y",
		AllowEdit: true,
		Label:     msg,
		Templates: templates,
		Validate:  validateYesNo,
	}

	response, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}

		return false, err
	}

	lowerMsg := strings.ToLower(response)
	if strings.HasPrefix(lowerMsg, "n") {
		return false, nil
	}

	return true, nil
}

func validateYesNo(msg string) error {
	lowerMsg := strings.ToLower(msg)

	if strings.HasPrefix(lowerMsg, "y") || strings.HasPrefix(lowerMsg, "n") {
		return nil
	}

	return fmt.Errorf("response must begin with 'y' or 'n'")
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
