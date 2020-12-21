package execution

import "github.com/newrelic/newrelic-cli/internal/install/types"

// MockStatusReporter is a mock implementation of the ExecutionStatusReporter
// interface that provides method spies for testing scenarios.
type MockStatusReporter struct {
	ReportRecipesAvailableErr       error
	ReportRecipeFailedErr           error
	ReportRecipeInstalledErr        error
	ReportCompleteErr               error
	ReportRecipesAvailableCallCount int
	ReportRecipeFailedCallCount     int
	ReportRecipeInstalledCallCount  int
	ReportCompleteCallCount         int
}

// NewMockStatusReporter returns a new instance of MockExecutionStatusReporter.
func NewMockStatusReporter() *MockStatusReporter {
	return &MockStatusReporter{}
}

func (r *MockStatusReporter) ReportRecipeFailed(event RecipeStatusEvent) error {
	r.ReportRecipeFailedCallCount++
	return r.ReportRecipeFailedErr
}

func (r *MockStatusReporter) ReportRecipeInstalled(event RecipeStatusEvent) error {
	r.ReportRecipeInstalledCallCount++
	return r.ReportRecipeInstalledErr
}

func (r *MockStatusReporter) ReportRecipesAvailable(recipes []types.Recipe) error {
	r.ReportRecipesAvailableCallCount++
	return r.ReportRecipesAvailableErr
}

func (r *MockStatusReporter) ReportComplete() error {
	r.ReportCompleteCallCount++
	return r.ReportCompleteErr
}
