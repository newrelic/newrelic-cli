package ux

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
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
		Default:   "y",
		AllowEdit: true,
		// IsConfirm: true,
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

	return fmt.Errorf("Response must begin with 'y' or 'n'.")
}
