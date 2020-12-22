package ux

type MockPrompter struct {
	PromptYesNoVal       bool
	PromptYesNoErr       error
	PromptYesNoCallCount int
}

func (p *MockPrompter) PromptYesNo(msg string) (bool, error) {
	p.PromptYesNoCallCount++
	return p.PromptYesNoVal, p.PromptYesNoErr
}
