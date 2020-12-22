package install

type mockPrompter struct {
	promptYesNoVal       bool
	promptYesNoErr       error
	promptYesNoCallCount int
}

func (p *mockPrompter) promptYesNo(msg string) (bool, error) {
	p.promptYesNoCallCount++
	return p.promptYesNoVal, p.promptYesNoErr
}
