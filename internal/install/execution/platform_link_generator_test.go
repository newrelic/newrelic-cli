//go:build unit
// +build unit

package execution

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

func TestGenerateRedirectURL_InstallSuccess(t *testing.T) {
	t.Parallel()

	g := NewPlatformLinkGenerator()

	recipeName := "infrastructure-agent-installer"
	entityGUID := "ABC123"
	recipe := types.OpenInstallationRecipe{
		Name:        recipeName,
		DisplayName: "Infrastructure Agent",
	}
	recipeStatus := &RecipeStatus{
		DisplayName: "Infrastructure Agent",
		Name:        recipeName,
		Status:      RecipeStatusTypes.INSTALLED,
		EntityGUID:  entityGUID,
	}
	installStatus := InstallStatus{
		recipesSelected: []types.OpenInstallationRecipe{recipe},
		Installed:       []*RecipeStatus{recipeStatus},
		EntityGUIDs:     []string{entityGUID},
	}

	expectedURL := fmt.Sprintf("https://%s/redirect/entity/%s", nrPlatformHostname(), entityGUID)
	result := g.GenerateRedirectURL(installStatus)
	require.Equal(t, expectedURL, result)
}

func TestGenerateRedirectURL_InstallPartialSuccess(t *testing.T) {
	t.Parallel()

	g := NewPlatformLinkGenerator()

	infraEntityGUID := "ABC123"
	infraRecipe := types.OpenInstallationRecipe{
		Name:        "infrastructure-agent-installer",
		DisplayName: "Infrastructure Agent",
	}
	logsRecipe := types.OpenInstallationRecipe{
		Name:        "logs-integration",
		DisplayName: "Logs integration",
	}
	installedRecipeStatus := &RecipeStatus{
		DisplayName: "Infrastructure Agent",
		Name:        "infrastructure-agent-installer",
		Status:      RecipeStatusTypes.INSTALLED,
		EntityGUID:  infraEntityGUID,
	}
	failedRecipeStatus := &RecipeStatus{
		DisplayName: "Logs integration",
		Name:        "logs-integration",
		Status:      RecipeStatusTypes.FAILED,
	}
	installStatus := InstallStatus{
		recipesSelected: []types.OpenInstallationRecipe{infraRecipe, logsRecipe},
		Installed:       []*RecipeStatus{installedRecipeStatus},
		Failed:          []*RecipeStatus{failedRecipeStatus},
		EntityGUIDs:     []string{infraEntityGUID},
		Statuses:        []*RecipeStatus{installedRecipeStatus, failedRecipeStatus},
	}
	expectedEncodedQueryParamSubstring := utils.Base64Encode(generateReferrerParam(infraEntityGUID))

	result := g.GenerateRedirectURL(installStatus)
	require.Contains(t, result, expectedEncodedQueryParamSubstring)
}

func TestGenerateRedirectURL_InstallFailed(t *testing.T) {
	t.Parallel()

	g := NewPlatformLinkGenerator()

	infraRecipe := types.OpenInstallationRecipe{
		Name:        "infrastructure-agent-installer",
		DisplayName: "Infrastructure Agent",
	}
	failedRecipeStatus := &RecipeStatus{
		DisplayName: "Infrastructure Agent",
		Name:        "infrastructure-agent-installer",
		Status:      RecipeStatusTypes.FAILED,
	}
	installStatus := InstallStatus{
		recipesSelected: []types.OpenInstallationRecipe{infraRecipe},
		Failed:          []*RecipeStatus{failedRecipeStatus},
		Statuses:        []*RecipeStatus{failedRecipeStatus},
	}
	expectedEncodedQueryParamSubstring := "eyJuZXJkbGV0SWQiOiJucjEtaW5zdGFsbC1uZXdyZWxpYy5pbnN0YWxsYXRpb24tcGxhbiIsInJlZmVycmVyIjoibmV3cmVsaWMtY2xpIn0="

	result := g.GenerateRedirectURL(installStatus)
	require.Contains(t, result, expectedEncodedQueryParamSubstring)
}

func TestGenerateRedirectURL_NoRecipesInstalled(t *testing.T) {
	t.Parallel()

	g := NewPlatformLinkGenerator()

	installStatus := InstallStatus{}
	expectedEncodedQueryParamSubstring := "eyJuZXJkbGV0SWQiOiJucjEtaW5zdGFsbC1uZXdyZWxpYy5pbnN0YWxsYXRpb24tcGxhbiIsInJlZmVycmVyIjoibmV3cmVsaWMtY2xpIn0="

	result := g.GenerateRedirectURL(installStatus)
	require.Contains(t, result, expectedEncodedQueryParamSubstring)
}
