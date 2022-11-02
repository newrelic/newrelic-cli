package execution

import "io"

type MockRecipeLogForwarder struct {
	optIn bool
}

func NewMockRecipeLogForwarder() *MockRecipeLogForwarder {
	return &MockRecipeLogForwarder{}
}

func (rlf *MockRecipeLogForwarder) PromptUserToSendLogs(reader io.Reader) bool {
	return false
}

func (rlf *MockRecipeLogForwarder) SendLogsToNewRelic(recipeName string, recipeOutput []string) {

}

func (rlf *MockRecipeLogForwarder) HasUserOptedIn() bool {
	return rlf.optIn
}

func (rlf *MockRecipeLogForwarder) SetUserOptedIn(val bool) {
	rlf.optIn = val
}
