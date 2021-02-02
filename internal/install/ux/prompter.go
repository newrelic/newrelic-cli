package ux

type Prompter interface {
	PromptYesNo(msg string) (bool, error)
	MultiSelect(msg string, options []string) ([]string, error)
}
