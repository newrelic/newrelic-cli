package execution

import "github.com/newrelic/newrelic-cli/internal/install/types"

// MockStatusReporter is a mock implementation of the ExecutionStatusReporter
// interface that provides method spies for testing scenarios.
type MockStatusReporter struct {
	ReportRecipeAvailableErr        error
	ReportRecipesAvailableErr       error
	ReportRecipeFailedErr           error
	ReportRecipeInstalledErr        error
	ReportRecipeInstallingErr       error
	ReportRecipeSkippedErr          error
	ReportCompleteErr               error
	ReportRecipeAvailableCallCount  int
	ReportRecipesAvailableCallCount int
	ReportRecipeFailedCallCount     int
	ReportRecipeInstalledCallCount  int
	ReportRecipeInstallingCallCount int
	ReportRecipeSkippedCallCount    int
	ReportCompleteCallCount         int
}

// NewMockStatusReporter returns a new instance of MockExecutionStatusReporter.
func NewMockStatusReporter() *MockStatusReporter {
	return &MockStatusReporter{}
}

func (r *MockStatusReporter) ReportRecipeFailed(status *StatusRollup, event RecipeStatusEvent) error {
	r.ReportRecipeFailedCallCount++
	return r.ReportRecipeFailedErr
}

func (r *MockStatusReporter) ReportRecipeInstalled(status *StatusRollup, event RecipeStatusEvent) error {
	r.ReportRecipeInstalledCallCount++
	return r.ReportRecipeInstalledErr
}

func (r *MockStatusReporter) ReportRecipeInstalling(status *StatusRollup, event RecipeStatusEvent) error {
	r.ReportRecipeInstallingCallCount++
	return r.ReportRecipeInstallingErr
}

func (r *MockStatusReporter) ReportRecipeSkipped(status *StatusRollup, event RecipeStatusEvent) error {
	r.ReportRecipeSkippedCallCount++
	return r.ReportRecipeSkippedErr
}

func (r *MockStatusReporter) ReportRecipeAvailable(status *StatusRollup, recipe types.Recipe) error {
	r.ReportRecipeAvailableCallCount++
	return r.ReportRecipeAvailableErr
}

func (r *MockStatusReporter) ReportRecipesAvailable(status *StatusRollup, recipes []types.Recipe) error {
	r.ReportRecipesAvailableCallCount++
	return r.ReportRecipesAvailableErr
}

func (r *MockStatusReporter) ReportComplete(status *StatusRollup) error {
	r.ReportCompleteCallCount++
	return r.ReportCompleteErr
}
