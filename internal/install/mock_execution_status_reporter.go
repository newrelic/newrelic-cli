package install

type mockExecutionStatusReporter struct {
	err error
}

func newMockExecutionStatusReporter() *mockExecutionStatusReporter {
	return &mockExecutionStatusReporter{}
}

func (r *mockExecutionStatusReporter) reportRecipeFailed(event recipeStatusEvent) error {
	return r.err
}

func (r *mockExecutionStatusReporter) reportRecipeInstalled(event recipeStatusEvent) error {
	return r.err
}

func (r *mockExecutionStatusReporter) reportRecipesAvailable(recipes []recipe) error {
	return r.err
}
