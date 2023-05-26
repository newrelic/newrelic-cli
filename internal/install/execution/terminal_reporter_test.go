//go:build unit
// +build unit

package execution

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestTerminalStatusReporter_interface(t *testing.T) {
	var r StatusSubscriber = NewTerminalStatusReporter()
	require.NotNil(t, r)
}

func Test_ShouldGenerateEntityLink(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockPlatformLinkGenerator()
	status := &InstallStatus{
		PlatformLinkGenerator: g,
	}
	recipeStatus := &RecipeStatus{
		Status: RecipeStatusTypes.INSTALLED,
	}
	status.Statuses = append(status.Statuses, recipeStatus)

	err := r.InstallComplete(status)
	require.NoError(t, err)
	require.Equal(t, 1, g.GenerateEntityLinkCallCount)
	require.Equal(t, 0, g.GenerateExplorerLinkCallCount)
}

func Test_ShouldNotGenerateEntityLink(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockPlatformLinkGenerator()

	status := &InstallStatus{
		PlatformLinkGenerator: g,
	}
	recipeStatus := &RecipeStatus{
		Status: RecipeStatusTypes.FAILED,
	}
	status.Statuses = append(status.Statuses, recipeStatus)

	err := r.InstallComplete(status)
	require.NoError(t, err)
	require.Equal(t, 0, g.GenerateEntityLinkCallCount)
	require.Equal(t, 0, g.GenerateExplorerLinkCallCount)
}

func Test_ShouldNotGenerateEntityLinkWhenNoRecipes(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockPlatformLinkGenerator()

	status := &InstallStatus{
		PlatformLinkGenerator: g,
	}

	err := r.InstallComplete(status)
	require.NoError(t, err)
	require.Equal(t, 0, g.GenerateEntityLinkCallCount)
	require.Equal(t, 0, g.GenerateExplorerLinkCallCount)
}

func Test_ShouldGenerateExplorerLink(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockPlatformLinkGenerator()

	status := &InstallStatus{
		PlatformLinkGenerator: g,
	}
	recipeStatus := &RecipeStatus{
		Status: RecipeStatusTypes.INSTALLED,
	}
	status.Statuses = append(status.Statuses, recipeStatus)
	status.successLinkConfig = types.OpenInstallationSuccessLinkConfig{
		Type:   "explorer",
		Filter: "\"`tags.language` = 'java'\"",
	}

	err := r.InstallComplete(status)
	require.NoError(t, err)
	require.Equal(t, 0, g.GenerateEntityLinkCallCount)
	require.Equal(t, 1, g.GenerateExplorerLinkCallCount)
}

func Test_ShouldGenerateLoggingLink(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockPlatformLinkGenerator()

	status := &InstallStatus{
		PlatformLinkGenerator: g,
	}

	loggingRecipeStatus := &RecipeStatus{
		DisplayName: "Logs integration",
		Name:        types.LoggingRecipeName,
		Status:      RecipeStatusTypes.INSTALLED,
	}

	status.Statuses = append(status.Statuses, loggingRecipeStatus)

	err := r.InstallComplete(status)
	require.NoError(t, err)
	require.Equal(t, 1, g.GenerateEntityLinkCallCount)
	require.Equal(t, 1, g.GenerateLoggingLinkCallCount)
}

func Test_ShouldNotGenerateExplorerLink(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockPlatformLinkGenerator()

	status := &InstallStatus{
		PlatformLinkGenerator: g,
	}
	recipeStatus := &RecipeStatus{
		Status: RecipeStatusTypes.FAILED,
	}
	status.Statuses = append(status.Statuses, recipeStatus)
	status.successLinkConfig = types.OpenInstallationSuccessLinkConfig{
		Type:   "explorer",
		Filter: "\"`tags.language` = 'java'\"",
	}

	err := r.InstallComplete(status)
	require.NoError(t, err)
	require.Equal(t, 0, g.GenerateEntityLinkCallCount)
	require.Equal(t, 0, g.GenerateExplorerLinkCallCount)
}

func Test_ShouldNotGenerateExplorerLinkWhenNoRecipes(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockPlatformLinkGenerator()

	status := &InstallStatus{
		PlatformLinkGenerator: g,
	}
	status.successLinkConfig = types.OpenInstallationSuccessLinkConfig{
		Type:   "explorer",
		Filter: "\"`tags.language` = 'java'\"",
	}

	err := r.InstallComplete(status)
	require.NoError(t, err)
	require.Equal(t, 0, g.GenerateEntityLinkCallCount)
	require.Equal(t, 0, g.GenerateExplorerLinkCallCount)
}

func TestTerminalStatusReporter_ShouldNotIncludeDetectedRecipeInSummary(t *testing.T) {
	r := NewTerminalStatusReporter()

	status := &InstallStatus{}
	recipeInstalled := &RecipeStatus{
		Name:        "test-recipe-installed",
		DisplayName: "Test Recipe Installed",
		Status:      RecipeStatusTypes.INSTALLED,
	}
	recipeDetected := &RecipeStatus{
		Name:        "test-recipe-detected",
		DisplayName: "Test Recipe Detected",
		Status:      RecipeStatusTypes.DETECTED,
	}
	recipeCanceled := &RecipeStatus{
		Name:        "test-recipe-canceld",
		DisplayName: "Test Recipe Canceled",
		Status:      RecipeStatusTypes.CANCELED,
	}

	status.Statuses = []*RecipeStatus{
		recipeInstalled,
		recipeDetected,
		recipeCanceled,
	}

	expected := []*RecipeStatus{
		recipeInstalled,
		recipeCanceled,
	}

	recipesToSummarize := r.getRecipesStatusesForInstallationSummary(status)

	require.Equal(t, len(expected), len(recipesToSummarize))
	require.Equal(t, expected[0].Name, recipesToSummarize[0].Name)
	require.Equal(t, expected[1].Name, recipesToSummarize[1].Name)
}
