package ux

type Prompter interface {
	PromptYesNo(msg string) (bool, error)
}
