package install

type mockExecutionStatusReporter struct {
	reportRecipesAvailableErr       error
	reportRecipeFailedErr           error
	reportRecipeInstalledErr        error
	reportRecipesAvailableCallCount int
	reportRecipeFailedCallCount     int
	reportRecipeInstalledCallCount  int
}

func newMockExecutionStatusReporter() *mockExecutionStatusReporter {
	return &mockExecutionStatusReporter{}
}

func (r *mockExecutionStatusReporter) reportRecipeFailed(event recipeStatusEvent) error {
	r.reportRecipeFailedCallCount++
	return r.reportRecipeFailedErr
}

func (r *mockExecutionStatusReporter) reportRecipeInstalled(event recipeStatusEvent) error {
	r.reportRecipeInstalledCallCount++
	return r.reportRecipeInstalledErr
}

func (r *mockExecutionStatusReporter) reportRecipesAvailable(recipes []recipe) error {
	r.reportRecipesAvailableCallCount++
	return r.reportRecipesAvailableErr
}
