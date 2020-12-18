package install

type mockExecutionStatusReporter struct {
	reportRecipesAvailableErr       error
	reportRecipeFailedErr           error
	reportRecipeInstalledErr        error
	reportCompleteErr               error
	reportRecipesAvailableCallCount int
	reportRecipeFailedCallCount     int
	reportRecipeInstalledCallCount  int
	reportCompleteCallCount         int
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

func (r *mockExecutionStatusReporter) reportComplete() error {
	r.reportCompleteCallCount++
	return r.reportCompleteErr
}
