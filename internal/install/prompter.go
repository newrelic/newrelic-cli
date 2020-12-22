package install

type prompter interface {
	promptYesNo(msg string) (bool, error)
}
