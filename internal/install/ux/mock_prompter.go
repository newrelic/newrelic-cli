package ux

type MockPrompter struct {
	PromptYesNoVal             bool
	PromptYesNoErr             error
	PromptYesNoCallCount       int
	PromptMultiSelectVal       []string
	PromptMultiSelectErr       error
	PromptMultiSelectCallCount int
	PromptMultiSelectAll       bool
}

func NewMockPrompter() *MockPrompter {
	return &MockPrompter{}
}

func (p *MockPrompter) PromptYesNo(msg string) (bool, error) {
	p.PromptYesNoCallCount++
	return p.PromptYesNoVal, p.PromptYesNoErr
}

func (p *MockPrompter) MultiSelect(msg string, options []string) ([]string, error) {
	p.PromptMultiSelectCallCount++

	if p.PromptMultiSelectAll {
		return options, nil
	}

	return p.PromptMultiSelectVal, p.PromptMultiSelectErr
}
