package execution

import (
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// MockStatusReporter is a mock implementation of the ExecutionStatusReporter
// interface that provides method spies for testing scenarios.
type MockStatusReporter struct {
	RecipeAvailableErr         error
	RecipesSelectedErr         error
	RecipeFailedErr            error
	RecipeInstalledErr         error
	RecipeInstallingErr        error
	RecipeRecommendedErr       error
	RecipeSkippedErr           error
	RecipeUnsupportedErr       error
	InstallCompleteErr         error
	InstallCanceledErr         error
	DiscoveryCompleteErr       error
	RecipeAvailableCallCount   int
	RecipesSelectedCallCount   int
	RecipeFailedCallCount      int
	RecipeInstalledCallCount   int
	RecipeInstallingCallCount  int
	RecipeRecommendedCallCount int
	RecipeSkippedCallCount     int
	InstallCompleteCallCount   int
	InstallCanceledCallCount   int
	DiscoveryCompleteCallCount int
	RecipeUnsupportedCallCount int

	ReportSkipped     map[string]int
	ReportInstalled   map[string]int
	ReportInstalling  map[string]int
	ReportRecommended map[string]int
	ReportFailed      map[string]int
	ReportAvailable   map[string]int

	GUIDs      []string
	Durations  []int64
	RecipeGUID map[string]string
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

	r.GUIDs = status.EntityGUIDs

	if len(r.RecipeGUID) == 0 {
		r.RecipeGUID = make(map[string]string)
	}

	for _, s := range status.Statuses {
		r.RecipeGUID[s.Name] = s.EntityGUID

		if s.ValidationDurationMilliseconds > 0 {
			r.Durations = append(r.Durations, s.ValidationDurationMilliseconds)
		}
	}

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

func (r *MockStatusReporter) RecipeAvailable(status *InstallStatus, recipe types.OpenInstallationRecipe) error {
	r.RecipeAvailableCallCount++
	if len(r.ReportAvailable) == 0 {
		r.ReportAvailable = make(map[string]int)
	}
	r.ReportAvailable[recipe.Name]++
	return r.RecipeAvailableErr
}

func (r *MockStatusReporter) RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error {
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

func (r *MockStatusReporter) RecipeUnsupported(status *InstallStatus, event RecipeStatusEvent) error {
	r.RecipeUnsupportedCallCount++
	return r.RecipeUnsupportedErr
}

func (r *MockStatusReporter) ObservabilityPacksFetchPending(status *InstallStatus) error {
	return nil
}

func (r *MockStatusReporter) ObservabilityPacksFetchSuccess(status *InstallStatus) error {
	return nil
}

func (r *MockStatusReporter) ObservabilityPacksFetchFailed(status *InstallStatus) error {
	return nil
}

func (r *MockStatusReporter) ObservabilityPacksInstallPending(status *InstallStatus) error {
	return nil
}

func (r *MockStatusReporter) ObservabilityPacksInstallSuccess(status *InstallStatus) error {
	return nil
}

func (r *MockStatusReporter) ObservabilityPacksInstallFailed(status *InstallStatus) error {
	return nil
}
