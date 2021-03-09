package execution

import (
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// MockStatusReporter is a mock implementation of the ExecutionStatusReporter
// interface that provides method spies for testing scenarios.
type MockStatusReporter struct {
	RecipeAvailableErr         error
	RecipesAvailableErr        error
	RecipesSelectedErr         error
	RecipeFailedErr            error
	RecipeInstalledErr         error
	RecipeInstallingErr        error
	RecipeRecommendedErr       error
	RecipeSkippedErr           error
	InstallCompleteErr         error
	InstallCanceledErr         error
	DiscoveryCompleteErr       error
	RecipeAvailableCallCount   int
	RecipesAvailableCallCount  int
	RecipesSelectedCallCount   int
	RecipeFailedCallCount      int
	RecipeInstalledCallCount   int
	RecipeInstallingCallCount  int
	RecipeRecommendedCallCount int
	RecipeSkippedCallCount     int
	InstallCompleteCallCount   int
	InstallCanceledCallCount   int
	DiscoveryCompleteCallCount int

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

func (r *MockStatusReporter) RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error {
	r.RecipeFailedCallCount++
	if len(r.ReportFailed) == 0 {
		r.ReportFailed = make(map[string]int)
	}
	r.ReportFailed[event.Recipe.Name]++
	return r.RecipeFailedErr
}

func (r *MockStatusReporter) RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error {
	r.RecipeInstalledCallCount++
	if len(r.ReportInstalled) == 0 {
		r.ReportInstalled = make(map[string]int)
	}
	r.ReportInstalled[event.Recipe.Name]++
	return r.RecipeInstalledErr
}

func (r *MockStatusReporter) RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error {
	r.RecipeInstallingCallCount++
	if len(r.ReportInstalling) == 0 {
		r.ReportInstalling = make(map[string]int)
	}
	r.ReportInstalling[event.Recipe.Name]++
	return r.RecipeInstallingErr
}

func (r *MockStatusReporter) RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error {
	r.RecipeRecommendedCallCount++
	if len(r.ReportRecommended) == 0 {
		r.ReportRecommended = make(map[string]int)
	}
	r.ReportRecommended[event.Recipe.Name]++
	return r.RecipeRecommendedErr
}

func (r *MockStatusReporter) RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error {
	r.RecipeSkippedCallCount++
	if len(r.ReportSkipped) == 0 {
		r.ReportSkipped = make(map[string]int)
	}
	r.ReportSkipped[event.Recipe.Name]++
	return r.RecipeSkippedErr
}

func (r *MockStatusReporter) RecipeAvailable(status *InstallStatus, recipe types.Recipe) error {
	r.RecipeAvailableCallCount++
	if len(r.ReportAvailable) == 0 {
		r.ReportAvailable = make(map[string]int)
	}
	r.ReportAvailable[recipe.Name]++
	return r.RecipeAvailableErr
}

func (r *MockStatusReporter) RecipesAvailable(status *InstallStatus, recipes []types.Recipe) error {
	r.RecipesAvailableCallCount++
	return r.RecipesAvailableErr
}

func (r *MockStatusReporter) RecipesSelected(status *InstallStatus, recipes []types.Recipe) error {
	r.RecipesSelectedCallCount++
	return r.RecipesSelectedErr
}

func (r *MockStatusReporter) InstallComplete(status *InstallStatus) error {
	r.InstallCompleteCallCount++
	return r.InstallCompleteErr
}

func (r *MockStatusReporter) InstallCanceled(status *InstallStatus) error {
	r.InstallCanceledCallCount++
	return r.InstallCanceledErr
}

func (r *MockStatusReporter) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
	r.DiscoveryCompleteCallCount++
	return r.DiscoveryCompleteErr
}
