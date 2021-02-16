package execution

import (
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// MockStatusReporter is a mock implementation of the ExecutionStatusReporter
// interface that provides method spies for testing scenarios.
type MockStatusReporter struct {
	ReportRecipeAvailableErr         error
	ReportRecipesAvailableErr        error
	ReportRecipeFailedErr            error
	ReportRecipeInstalledErr         error
	ReportRecipeInstallingErr        error
	ReportRecipeRecommendedErr       error
	ReportRecipeSkippedErr           error
	ReportCompleteErr                error
	ReportRecipeAvailableCallCount   int
	ReportRecipesAvailableCallCount  int
	ReportRecipeFailedCallCount      int
	ReportRecipeInstalledCallCount   int
	ReportRecipeInstallingCallCount  int
	ReportRecipeRecommendedCallCount int
	ReportRecipeSkippedCallCount     int
	ReportCompleteCallCount          int

	ReportSkipped     map[string]int
	ReportInstalled   map[string]int
	ReportInstalling  map[string]int
	ReportRecommended map[string]int
	ReportFailed      map[string]int
	ReportAvailable   map[string]int
}

// NewMockStatusReporter returns a new instance of MockExecutionStatusReporter.
func NewMockStatusReporter() *MockStatusReporter {
	return &MockStatusReporter{}
}

func (r *MockStatusReporter) ReportRecipeFailed(status *StatusRollup, event RecipeStatusEvent) error {
	r.ReportRecipeFailedCallCount++
	if len(r.ReportFailed) == 0 {
		r.ReportFailed = make(map[string]int)
	}
	r.ReportFailed[event.Recipe.Name]++
	return r.ReportRecipeFailedErr
}

func (r *MockStatusReporter) ReportRecipeInstalled(status *StatusRollup, event RecipeStatusEvent) error {
	r.ReportRecipeInstalledCallCount++
	if len(r.ReportInstalled) == 0 {
		r.ReportInstalled = make(map[string]int)
	}
	r.ReportInstalled[event.Recipe.Name]++
	return r.ReportRecipeInstalledErr
}

func (r *MockStatusReporter) ReportRecipeInstalling(status *StatusRollup, event RecipeStatusEvent) error {
	r.ReportRecipeInstallingCallCount++
	if len(r.ReportInstalling) == 0 {
		r.ReportInstalling = make(map[string]int)
	}
	r.ReportInstalling[event.Recipe.Name]++
	return r.ReportRecipeInstallingErr
}

func (r *MockStatusReporter) ReportRecipeRecommended(status *StatusRollup, event RecipeStatusEvent) error {
	r.ReportRecipeRecommendedCallCount++
	if len(r.ReportRecommended) == 0 {
		r.ReportRecommended = make(map[string]int)
	}
	r.ReportRecommended[event.Recipe.Name]++
	return r.ReportRecipeRecommendedErr
}

func (r *MockStatusReporter) ReportRecipeSkipped(status *StatusRollup, event RecipeStatusEvent) error {
	r.ReportRecipeSkippedCallCount++
	if len(r.ReportSkipped) == 0 {
		r.ReportSkipped = make(map[string]int)
	}
	r.ReportSkipped[event.Recipe.Name]++
	return r.ReportRecipeSkippedErr
}

func (r *MockStatusReporter) ReportRecipeAvailable(status *StatusRollup, recipe types.Recipe) error {
	r.ReportRecipeAvailableCallCount++
	if len(r.ReportAvailable) == 0 {
		r.ReportAvailable = make(map[string]int)
	}
	r.ReportAvailable[recipe.Name]++
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
