package execution

import "io"

type MockRecipeLogForwarder struct {
}

func NewMockRecipeLogForwarder() *MockRecipeLogForwarder {
	return &MockRecipeLogForwarder{}
}

func (rlf *MockRecipeLogForwarder) PromptUserToSendLogs(reader io.Reader) bool {
	return false
}

func (rlf *MockRecipeLogForwarder) SendLogsToNewRelic(recipeName string, recipeOutput []string) {

}
