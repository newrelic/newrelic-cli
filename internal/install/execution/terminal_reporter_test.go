package execution

import (
	"bytes"
	"fmt"
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
	loggingAgentControlRecipeStatus := &RecipeStatus{
		DisplayName: "Logs integration",
		Name:        types.LoggingAgentControlRecipeName,
		Status:      RecipeStatusTypes.INSTALLED,
	}

	status.Statuses = append(status.Statuses, loggingRecipeStatus)
	status.Statuses = append(status.Statuses, loggingAgentControlRecipeStatus)

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

func TestPrintInstallationSummaryShouldPrint(t *testing.T) {
	r := NewTerminalStatusReporter()
	var output bytes.Buffer

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

	r.printInstallationSummary(&output, status)
	s := output.String()
	fmt.Print(s)

	require.Contains(t, s, "Test Recipe Installed  (installed)")
	require.NotContains(t, s, "Detected")
	require.Contains(t, s, "Test Recipe Canceled  (canceled)")
}

func Test_ShouldNotGenerateRedirectURLForAgentControlOnly(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockPlatformLinkGenerator()
	g.GenerateRedirectURLVal = "https://one.newrelic.com/redirect"

	status := &InstallStatus{
		PlatformLinkGenerator: g,
	}
	agentControlRecipeStatus := &RecipeStatus{
		Name:   types.AgentControlRecipeName,
		Status: RecipeStatusTypes.INSTALLED,
	}
	status.Statuses = append(status.Statuses, agentControlRecipeStatus)

	err := r.InstallComplete(status)
	require.NoError(t, err)

	require.Equal(t, 0, g.GenerateRedirectURLCallCount)
}

func Test_ShouldGenerateRedirectURLForAgentControlPlusOtherRecipes(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockPlatformLinkGenerator()
	g.GenerateRedirectURLVal = "https://one.newrelic.com/redirect"

	status := &InstallStatus{
		PlatformLinkGenerator: g,
	}
	agentControlRecipeStatus := &RecipeStatus{
		Name:   types.AgentControlRecipeName,
		Status: RecipeStatusTypes.INSTALLED,
	}
	otherRecipeStatus := &RecipeStatus{
		Name:   "infrastructure-agent-installer",
		Status: RecipeStatusTypes.INSTALLED,
	}
	status.Statuses = append(status.Statuses, agentControlRecipeStatus, otherRecipeStatus)

	err := r.InstallComplete(status)
	require.NoError(t, err)

	require.Equal(t, 1, g.GenerateRedirectURLCallCount)
}

func Test_ShouldGenerateRedirectURLForNonAgentControlRecipes(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockPlatformLinkGenerator()
	g.GenerateRedirectURLVal = "https://one.newrelic.com/redirect"

	status := &InstallStatus{
		PlatformLinkGenerator: g,
	}
	recipeStatus := &RecipeStatus{
		Name:   "infrastructure-agent-installer",
		Status: RecipeStatusTypes.INSTALLED,
	}
	status.Statuses = append(status.Statuses, recipeStatus)

	err := r.InstallComplete(status)
	require.NoError(t, err)

	require.Equal(t, 1, g.GenerateRedirectURLCallCount)
}

func Test_isOnlyAgentControlInstallation_OnlyAgentControl(t *testing.T) {
	r := NewTerminalStatusReporter()

	status := &InstallStatus{}
	agentControlRecipe := &RecipeStatus{
		Name:   types.AgentControlRecipeName,
		Status: RecipeStatusTypes.INSTALLED,
	}
	loggingAgentControlRecipe := &RecipeStatus{
		Name:   types.LoggingAgentControlRecipeName,
		Status: RecipeStatusTypes.INSTALLED,
	}
	status.Statuses = []*RecipeStatus{agentControlRecipe, loggingAgentControlRecipe}

	result := r.isOnlyAgentControlInstallation(status)
	require.True(t, result, "should return true when only agent-control recipes are installed")
}

func Test_isOnlyAgentControlInstallation_MixedRecipes(t *testing.T) {
	r := NewTerminalStatusReporter()

	status := &InstallStatus{}
	agentControlRecipe := &RecipeStatus{
		Name:   types.AgentControlRecipeName,
		Status: RecipeStatusTypes.INSTALLED,
	}
	otherRecipe := &RecipeStatus{
		Name:   "infrastructure-agent-installer",
		Status: RecipeStatusTypes.INSTALLED,
	}
	status.Statuses = []*RecipeStatus{agentControlRecipe, otherRecipe}

	result := r.isOnlyAgentControlInstallation(status)
	require.False(t, result, "should return false when agent-control is mixed with other recipes")
}

func Test_isOnlyAgentControlInstallation_NoAgentControl(t *testing.T) {
	r := NewTerminalStatusReporter()

	status := &InstallStatus{}
	otherRecipe := &RecipeStatus{
		Name:   "infrastructure-agent-installer",
		Status: RecipeStatusTypes.INSTALLED,
	}
	status.Statuses = []*RecipeStatus{otherRecipe}

	result := r.isOnlyAgentControlInstallation(status)
	require.False(t, result, "should return false when no agent-control recipes are installed")
}

func Test_isOnlyAgentControlInstallation_NoInstalledRecipes(t *testing.T) {
	r := NewTerminalStatusReporter()

	status := &InstallStatus{}
	failedRecipe := &RecipeStatus{
		Name:   types.AgentControlRecipeName,
		Status: RecipeStatusTypes.FAILED,
	}
	status.Statuses = []*RecipeStatus{failedRecipe}

	result := r.isOnlyAgentControlInstallation(status)
	require.False(t, result, "should return false when no recipes are installed (only failed)")
}

func Test_isOnlyAgentControlAttempt_OnlyAgentControlFailed(t *testing.T) {
	r := NewTerminalStatusReporter()

	status := &InstallStatus{}
	failedAgentControlRecipe := &RecipeStatus{
		Name:   types.AgentControlRecipeName,
		Status: RecipeStatusTypes.FAILED,
	}
	status.Statuses = []*RecipeStatus{failedAgentControlRecipe}

	result := r.isOnlyAgentControlAttempt(status)
	require.True(t, result, "should return true when only agent-control recipes were attempted (even if failed)")
}

func Test_isOnlyAgentControlAttempt_MixedRecipesFailed(t *testing.T) {
	r := NewTerminalStatusReporter()

	status := &InstallStatus{}
	failedAgentControlRecipe := &RecipeStatus{
		Name:   types.AgentControlRecipeName,
		Status: RecipeStatusTypes.FAILED,
	}
	failedOtherRecipe := &RecipeStatus{
		Name:   "infrastructure-agent-installer",
		Status: RecipeStatusTypes.FAILED,
	}
	status.Statuses = []*RecipeStatus{failedAgentControlRecipe, failedOtherRecipe}

	result := r.isOnlyAgentControlAttempt(status)
	require.False(t, result, "should return false when mixed recipes were attempted")
}

func Test_isOnlyAgentControlAttempt_OnlyAgentControlInstalled(t *testing.T) {
	r := NewTerminalStatusReporter()

	status := &InstallStatus{}
	installedAgentControlRecipe := &RecipeStatus{
		Name:   types.AgentControlRecipeName,
		Status: RecipeStatusTypes.INSTALLED,
	}
	status.Statuses = []*RecipeStatus{installedAgentControlRecipe}

	result := r.isOnlyAgentControlAttempt(status)
	require.True(t, result, "should return true when only agent-control recipes were attempted and installed")
}

func Test_isOnlyAgentControlAttempt_NoRecipes(t *testing.T) {
	r := NewTerminalStatusReporter()

	status := &InstallStatus{}
	status.Statuses = []*RecipeStatus{}

	result := r.isOnlyAgentControlAttempt(status)
	require.False(t, result, "should return false when no recipes were attempted")
}

func Test_isOnlyAgentControlAttempt_BothAgentControlTypes(t *testing.T) {
	r := NewTerminalStatusReporter()

	status := &InstallStatus{}
	agentControlRecipe := &RecipeStatus{
		Name:   types.AgentControlRecipeName,
		Status: RecipeStatusTypes.FAILED,
	}
	loggingAgentControlRecipe := &RecipeStatus{
		Name:   types.LoggingAgentControlRecipeName,
		Status: RecipeStatusTypes.FAILED,
	}
	status.Statuses = []*RecipeStatus{agentControlRecipe, loggingAgentControlRecipe}

	result := r.isOnlyAgentControlAttempt(status)
	require.True(t, result, "should return true when only agent-control recipes (both types) were attempted")
}

func Test_InstallComplete_IncompleteAgentControlShowsDocLink(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockPlatformLinkGenerator()
	g.GenerateRedirectURLVal = "https://one.newrelic.com/redirect"

	status := &InstallStatus{
		PlatformLinkGenerator: g,
	}
	failedAgentControlRecipe := &RecipeStatus{
		Name:        types.AgentControlRecipeName,
		DisplayName: "Agent Control",
		Status:      RecipeStatusTypes.FAILED,
	}
	status.Statuses = append(status.Statuses, failedAgentControlRecipe)

	err := r.InstallComplete(status)
	require.NoError(t, err)
	require.Equal(t, 0, g.GenerateRedirectURLCallCount)
	require.Equal(t, 1, g.GenerateGuidedInstallDocLinkCallCount)
}

func Test_InstallComplete_IncompleteNonAgentControlShowsRedirectLink(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockPlatformLinkGenerator()
	g.GenerateRedirectURLVal = "https://one.newrelic.com/redirect"

	status := &InstallStatus{
		PlatformLinkGenerator: g,
	}
	failedOtherRecipe := &RecipeStatus{
		Name:        "infrastructure-agent-installer",
		DisplayName: "Infrastructure Agent",
		Status:      RecipeStatusTypes.FAILED,
	}
	status.Statuses = append(status.Statuses, failedOtherRecipe)

	err := r.InstallComplete(status)
	require.NoError(t, err)

	require.Equal(t, 1, g.GenerateRedirectURLCallCount)
}

func Test_InstallComplete_IncompleteAgentControlNoPlatformLinkGenerator(t *testing.T) {
	r := NewTerminalStatusReporter()

	status := &InstallStatus{
		PlatformLinkGenerator: nil,
	}
	failedAgentControlRecipe := &RecipeStatus{
		Name:        types.AgentControlRecipeName,
		DisplayName: "Agent Control",
		Status:      RecipeStatusTypes.FAILED,
	}
	status.Statuses = append(status.Statuses, failedAgentControlRecipe)

	err := r.InstallComplete(status)
	require.NoError(t, err)
}

func Test_InstallComplete_IncompleteNonAgentControlNoLinkGenerator(t *testing.T) {
	r := NewTerminalStatusReporter()

	status := &InstallStatus{
		PlatformLinkGenerator: nil,
	}
	failedOtherRecipe := &RecipeStatus{
		Name:        "infrastructure-agent-installer",
		DisplayName: "Infrastructure Agent",
		Status:      RecipeStatusTypes.FAILED,
	}
	status.Statuses = append(status.Statuses, failedOtherRecipe)

	err := r.InstallComplete(status)
	require.NoError(t, err)
}

func Test_InstallComplete_SuccessfulMixedInstallation(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockPlatformLinkGenerator()
	g.GenerateRedirectURLVal = "https://one.newrelic.com/redirect"

	status := &InstallStatus{
		PlatformLinkGenerator: g,
	}
	agentControlRecipe := &RecipeStatus{
		Name:        types.AgentControlRecipeName,
		DisplayName: "Agent Control",
		Status:      RecipeStatusTypes.INSTALLED,
	}
	infraRecipe := &RecipeStatus{
		Name:        "infrastructure-agent-installer",
		DisplayName: "Infrastructure Agent",
		Status:      RecipeStatusTypes.INSTALLED,
	}
	status.Statuses = append(status.Statuses, agentControlRecipe, infraRecipe)

	err := r.InstallComplete(status)
	require.NoError(t, err)

	require.Equal(t, 1, g.GenerateRedirectURLCallCount)
	require.Equal(t, 0, g.GenerateGuidedInstallDocLinkCallCount)
}
