package uninstall

import (
	"github.com/newrelic/newrelic-cli/internal/install"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
)

// SurveyConfirmer asks yes/no questions on the terminal, reusing the same prompt library install uses.
type SurveyConfirmer struct {
	prompter install.Prompter
}

// NewSurveyConfirmer returns a terminal-backed Confirmer.
func NewSurveyConfirmer() *SurveyConfirmer {
	return &SurveyConfirmer{prompter: ux.NewPromptUIPrompter()}
}

func (c *SurveyConfirmer) Confirm(prompt string) (bool, error) {
	return c.prompter.PromptYesNo(prompt)
}
