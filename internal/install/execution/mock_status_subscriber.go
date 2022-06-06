package execution

import (
	"errors"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// MockStatusSubscriber is a mock implementation of the ExecutionStatusReporter
// interface that provides method spies for testing scenarios.
type MockStatusSubscriber struct {
	RecipeAvailableErr         error
	RecipesSelectedErr         error
	RecipeFailedErr            error
	RecipeInstalledErr         error
	RecipeInstallingErr        error
	RecipeRecommendedErr       error
	RecipeSkippedErr           error
	RecipeUnsupportedErr       error
	InstallStartedErr          error
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
	InstallStartedCallCount    int
	InstallCompleteCallCount   int
	InstallCanceledCallCount   int
	DiscoveryCompleteCallCount int
	RecipeUnsupportedCallCount int
	RecipeDetectedCallCount    int
	RecipeCanceledCallCount    int

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
func NewMockStatusReporter() *MockStatusSubscriber {
	return &MockStatusSubscriber{}
}

func (r *MockStatusSubscriber) RecipeDetected(status *InstallStatus, event RecipeStatusEvent) error {
	r.RecipeDetectedCallCount++
	return nil
}

func (r *MockStatusSubscriber) RecipeCanceled(status *InstallStatus, event RecipeStatusEvent) error {
	r.RecipeCanceledCallCount++
	return nil
}

func (r *MockStatusSubscriber) UpdateRequired(status *InstallStatus) error {
	return nil
}

func (r *MockStatusSubscriber) RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error {
	r.RecipeFailedCallCount++
	if len(r.ReportFailed) == 0 {
		r.ReportFailed = make(map[string]int)
	}
	r.ReportFailed[event.Recipe.Name]++
	return r.RecipeFailedErr
}

func (r *MockStatusSubscriber) RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error {
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

		if s.ValidationDurationMs > 0 {
			r.Durations = append(r.Durations, s.ValidationDurationMs)
		}
	}

	return r.RecipeInstalledErr
}

func (r *MockStatusSubscriber) RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error {
	r.RecipeInstallingCallCount++
	if len(r.ReportInstalling) == 0 {
		r.ReportInstalling = make(map[string]int)
	}
	r.ReportInstalling[event.Recipe.Name]++
	return r.RecipeInstallingErr
}

func (r *MockStatusSubscriber) RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error {
	r.RecipeRecommendedCallCount++
	if len(r.ReportRecommended) == 0 {
		r.ReportRecommended = make(map[string]int)
	}
	r.ReportRecommended[event.Recipe.Name]++
	return r.RecipeRecommendedErr
}

func (r *MockStatusSubscriber) RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error {
	r.RecipeSkippedCallCount++
	if len(r.ReportSkipped) == 0 {
		r.ReportSkipped = make(map[string]int)
	}
	r.ReportSkipped[event.Recipe.Name]++
	return r.RecipeSkippedErr
}

func (r *MockStatusSubscriber) RecipeAvailable(status *InstallStatus, event RecipeStatusEvent) error {
	r.RecipeAvailableCallCount++
	if len(r.ReportAvailable) == 0 {
		r.ReportAvailable = make(map[string]int)
	}
	r.ReportAvailable[event.Recipe.Name]++
	return r.RecipeAvailableErr
}

func (r *MockStatusSubscriber) RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error {
	r.RecipesSelectedCallCount++
	return r.RecipesSelectedErr
}

func (r *MockStatusSubscriber) InstallStarted(status *InstallStatus) error {
	r.InstallStartedCallCount++
	return r.InstallStartedErr
}

func (r *MockStatusSubscriber) InstallComplete(status *InstallStatus) error {
	r.InstallCompleteErr = errors.New(status.Error.Message)
	r.InstallCompleteCallCount++
	return r.InstallCompleteErr
}

func (r *MockStatusSubscriber) InstallCanceled(status *InstallStatus) error {
	r.InstallCanceledCallCount++
	return r.InstallCanceledErr
}

func (r *MockStatusSubscriber) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
	r.DiscoveryCompleteCallCount++
	return r.DiscoveryCompleteErr
}

func (r *MockStatusSubscriber) RecipeUnsupported(status *InstallStatus, event RecipeStatusEvent) error {
	r.RecipeUnsupportedCallCount++
	return r.RecipeUnsupportedErr
}
