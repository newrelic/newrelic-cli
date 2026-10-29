//go:build unit

package execution

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

func TestGenerateRedirectURL_ShouldGenerateEntityLinkWhenOneInstalledWithNoFailure(t *testing.T) {
	t.Parallel()

	recipeName := "infrastructure-agent-installer"
	b := newPlatformLinkGeneratorBuilder()
	b.recipeStatusUpdate(recipeName, "Installed")
	g, s := b.build()

	expectedURL := fmt.Sprintf("https://%s/redirect/entity/%s", nrPlatformHostname(), recipeName)

	result := g.GenerateRedirectURL(*s)
	require.Equal(t, 1, len(s.EntityGUIDs))
	require.Equal(t, 1, len(s.Installed))
	require.Equal(t, expectedURL, result)
}

func TestGenerateRedirectURL_ShoudGenerateStatusLinkWhenMoreThanOneInstalled(t *testing.T) {
	t.Parallel()

	infraName := "infrastructure-agent-installer"
	loggingName := "Log-integration"

	b := newPlatformLinkGeneratorBuilder()
	b.recipeStatusUpdate(infraName, "Installed")
	b.recipeStatusUpdate(loggingName, "Installed")
	g, s := b.build()

	result := g.GenerateRedirectURL(*s)
	require.Contains(t, result, "explorer")
}

func TestGenerateLoggingURL_InstallSuccess(t *testing.T) {
	t.Parallel()

	rName := "MXxBUE18QVBQTElDQVRJT058OTE2NzQxNg"
	loggingName := "Log-integration"

	b := newPlatformLinkGeneratorBuilder()
	b.recipeStatusUpdate(rName, "Installed")
	b.recipeStatusUpdate(loggingName, "Installed")
	g, s := b.build()
	accountID := configAPI.GetActiveProfileAccountID()

	launcherEncodedParams := "eyJxdWVyeSI6IlwiZW50aXR5Lmd1aWQuSU5GUkFcIjpcIk1YeEJVRTE4UVZCUVRFbERRVlJKVDA1OE9URTJOelF4TmdcIiIsInJlZmVycmVyIjoibmV3cmVsaWMtY2xpIn0="
	expectedLoggingLink := fmt.Sprintf("https://%s/launcher/logger.log-launcher?platform[accountId]=%d&launcher=%s", nrPlatformHostname(), accountID, launcherEncodedParams)

	redirectURLResult := g.GenerateRedirectURL(*s)
	loggingLinkResult := g.GenerateLoggingLink(rName)
	require.Contains(t, redirectURLResult, "explorer")
	require.Contains(t, loggingLinkResult, expectedLoggingLink)
}

func TestGenerateRedirectURL_InstallPartialSuccess(t *testing.T) {
	t.Parallel()

	infraName := "infrastructure-agent-installer"
	loggingName := "Log-integration"

	b := newPlatformLinkGeneratorBuilder()
	b.recipeStatusUpdate(infraName, "Installed")
	b.recipeStatusUpdate(loggingName, "Failed")
	g, s := b.build()

	result := g.GenerateRedirectURL(*s)
	expectedEncodedQueryParamSubstring := utils.Base64Encode(g.generateReferrerParam(infraName, b.installStatus.InstallID))

	require.Equal(t, 1, len(s.EntityGUIDs))
	require.Equal(t, 1, len(s.Installed))
	require.Equal(t, 1, len(s.Failed))
	require.Contains(t, result, expectedEncodedQueryParamSubstring)
}

func TestGenerateRedirectURL_InstallFailed(t *testing.T) {
	t.Parallel()

	infraName := "infrastructure-agent-installer"

	b := newPlatformLinkGeneratorBuilder()
	b.recipeStatusUpdate(infraName, "Failed")
	g, s := b.build()

	result := g.GenerateRedirectURL(*s)
	require.Contains(t, result, "explorer")
}

func TestGenerateRedirectURL_InstallCanceled(t *testing.T) {
	t.Parallel()

	infraName := "infrastructure-agent-installer"

	b := newPlatformLinkGeneratorBuilder()
	b.recipeStatusUpdate(infraName, "Canceled")
	g, s := b.build()

	result := g.GenerateRedirectURL(*s)
	require.Contains(t, result, "explorer")
}

func TestGenerateRedirectURL_NoRecipesInstalled(t *testing.T) {
	t.Parallel()

	b := newPlatformLinkGeneratorBuilder()
	g, s := b.build()

	result := g.GenerateRedirectURL(*s)
	require.Contains(t, result, "explorer")
}

func TestGetAccountPlanManagementURL(t *testing.T) {
	t.Parallel()

	result := GetAccountPlanManagementURL()
	require.Contains(t, result, "plan-management/home?account=")
}

func TestGenerateFleetLink_WithFleetGUID(t *testing.T) {
	t.Parallel()

	fleetGUID := "MTIyMTMwNjh8TkdFUHxGTEVFVHwwMTljOGY0OS03ZDY2LTczOGEtYjQ4Ny03NjI5YzE1ZjNiOWI"

	g := NewPlatformLinkGenerator()
	g.apiKey = ""

	result := g.GenerateFleetLink(fleetGUID)

	require.Contains(t, result, "launcher/new-relic-control.launcher")
	require.Contains(t, result, "pane=")
	require.Contains(t, result, "platform[accountId]")
	require.Contains(t, result, nrPlatformHostname())
}

func TestGenerateFleetLink_EmptyFleetGUID(t *testing.T) {
	t.Parallel()

	g := NewPlatformLinkGenerator()
	g.apiKey = ""

	result := g.GenerateFleetLink("")

	expectedURL := fmt.Sprintf("https://%s/fleet", nrPlatformHostname())
	require.Equal(t, expectedURL, result)
}

func TestGenerateFleetLink_WhitespaceOnlyFleetGUID(t *testing.T) {
	t.Parallel()

	g := NewPlatformLinkGenerator()
	g.apiKey = ""

	testCases := []string{
		"  ",
		"\t",
		"\n",
		"  \t\n  ",
	}

	for _, tc := range testCases {
		result := g.GenerateFleetLink(tc)
		expectedURL := fmt.Sprintf("https://%s/fleet", nrPlatformHostname())
		require.Equal(t, expectedURL, result, "Should return base URL for whitespace-only GUID: %q", tc)
	}
}

func TestGenerateFleetLauncherParams_TrimsWhitespace(t *testing.T) {
	t.Parallel()

	g := NewPlatformLinkGenerator()

	testCases := []struct {
		name         string
		inputGUID    string
		expectedGUID string
	}{
		{
			name:         "GUID with leading spaces",
			inputGUID:    "  MTIyMTMwNjh8TkdFUHxGTEVFVHwwMTljOGY0OS03ZDY2LTczOGEtYjQ4Ny03NjI5YzE1ZjNiOWI",
			expectedGUID: "MTIyMTMwNjh8TkdFUHxGTEVFVHwwMTljOGY0OS03ZDY2LTczOGEtYjQ4Ny03NjI5YzE1ZjNiOWI",
		},
		{
			name:         "GUID with trailing spaces",
			inputGUID:    "MTIyMTMwNjh8TkdFUHxGTEVFVHwwMTljOGY0OS03ZDY2LTczOGEtYjQ4Ny03NjI5YzE1ZjNiOWI  ",
			expectedGUID: "MTIyMTMwNjh8TkdFUHxGTEVFVHwwMTljOGY0OS03ZDY2LTczOGEtYjQ4Ny03NjI5YzE1ZjNiOWI",
		},
		{
			name:         "GUID with leading and trailing spaces",
			inputGUID:    "  MTIyMTMwNjh8TkdFUHxGTEVFVHwwMTljOGY0OS03ZDY2LTczOGEtYjQ4Ny03NjI5YzE1ZjNiOWI  ",
			expectedGUID: "MTIyMTMwNjh8TkdFUHxGTEVFVHwwMTljOGY0OS03ZDY2LTczOGEtYjQ4Ny03NjI5YzE1ZjNiOWI",
		},
		{
			name:         "GUID with tabs and newlines",
			inputGUID:    "\t\nMTIyMTMwNjh8TkdFUHxGTEVFVHwwMTljOGY0OS03ZDY2LTczOGEtYjQ4Ny03NjI5YzE1ZjNiOWI\n\t",
			expectedGUID: "MTIyMTMwNjh8TkdFUHxGTEVFVHwwMTljOGY0OS03ZDY2LTczOGEtYjQ4Ny03NjI5YzE1ZjNiOWI",
		},
		{
			name:         "GUID without whitespace",
			inputGUID:    "MTIyMTMwNjh8TkdFUHxGTEVFVHwwMTljOGY0OS03ZDY2LTczOGEtYjQ4Ny03NjI5YzE1ZjNiOWI",
			expectedGUID: "MTIyMTMwNjh8TkdFUHxGTEVFVHwwMTljOGY0OS03ZDY2LTczOGEtYjQ4Ny03NjI5YzE1ZjNiOWI",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := g.generateFleetLauncherParams(tc.inputGUID)

			require.Contains(t, result, tc.expectedGUID, "Fleet GUID should be trimmed of whitespace")

			if tc.inputGUID != tc.expectedGUID {
				require.NotContains(t, result, tc.inputGUID, "Fleet GUID should not contain original whitespace")
			}

			require.Contains(t, result, `"nerdletId":"fleet.detail"`)
			require.Contains(t, result, `"referrer":"newrelic-cli"`)
			require.Contains(t, result, `"fleetGuid":"`+tc.expectedGUID+`"`)
		})
	}
}

type platformLinkGeneratorBuilder struct {
	platformLinkGenerator *PlatformLinkGenerator
	installStatus         *InstallStatus
}

func newPlatformLinkGeneratorBuilder() *platformLinkGeneratorBuilder {
	p := &platformLinkGeneratorBuilder{
		platformLinkGenerator: NewPlatformLinkGenerator(),
	}

	p.installStatus = NewInstallStatus(types.InstallerContext{}, make([]StatusSubscriber, 0), p.platformLinkGenerator)
	// We set an API key in the unit test so we don't make an real HTTP request
	// to the New Relic short URL service (see integration test), and so we can test
	// the query param being added for the fallback installation strategy below.
	p.platformLinkGenerator.apiKey = ""
	return p
}

func (p *platformLinkGeneratorBuilder) recipeStatusUpdate(rn, status string) *platformLinkGeneratorBuilder {

	r := types.OpenInstallationRecipe{
		Name: rn,
	}

	rs := RecipeStatusEvent{
		Recipe: r,
	}

	switch status {
	case "Failed":
		p.installStatus.RecipeFailed(rs)
	case "Canceled":
		p.installStatus.RecipeCanceled(rs)
	case "Installed":
		// just for testing, assume name is the same as entity id
		rs.EntityGUID = rn
		p.installStatus.RecipeInstalled(rs)
	}

	return p
}

func (p *platformLinkGeneratorBuilder) build() (*PlatformLinkGenerator, *InstallStatus) {
	p.installStatus.completed(errors.New(""))
	return p.platformLinkGenerator, p.installStatus
}
