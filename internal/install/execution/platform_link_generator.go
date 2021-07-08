package execution

import (
	"fmt"
	"strings"

	"github.com/newrelic/newrelic-cli/internal/configuration"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/pkg/region"
)

type PlatformLinkGenerator struct{}

var nrPlatformHostnames = struct {
	Staging string
	US      string
	EU      string
}{
	Staging: "staging-one.newrelic.com",
	US:      "one.newrelic.com",
	EU:      "one.eu.newrelic.com",
}

func NewPlatformLinkGenerator() *PlatformLinkGenerator {
	return &PlatformLinkGenerator{}
}

func (g *PlatformLinkGenerator) GenerateExplorerLink(filter string) string {
	return generateExplorerLink(filter)
}

func (g *PlatformLinkGenerator) GenerateEntityLink(entityGUID string) string {
	return generateEntityLink(entityGUID)
}

// GenerateRedirectURL creates a URL for the user to navigate to after running
// through an installation. The URL is displayed in the CLI out as well and is
// also provided in the nerdstorage document. This provides the user two options
// to see their data - click from the CLI output or from the frontend.
func (g *PlatformLinkGenerator) GenerateRedirectURL(status InstallStatus) string {
	if status.hasAnyRecipeStatus(RecipeStatusTypes.INSTALLED) {
		switch t := status.successLinkConfig.Type; {
		case strings.EqualFold(string(t), "explorer"):
			return g.GenerateExplorerLink(status.successLinkConfig.Filter)
		default:
			return g.GenerateEntityLink(status.HostEntityGUID())
		}
	}

	return ""
}

func generateExplorerLink(filter string) string {
	return fmt.Sprintf("https://%s/launcher/nr1-core.explorer?platform[filters]=%s&platform[accountId]=%d",
		nrPlatformHostname(),
		utils.Base64Encode(filter),
		configuration.GetActiveProfileAccountIDWithFlagOverride(),
	)
}

func generateEntityLink(entityGUID string) string {
	return fmt.Sprintf("https://%s/redirect/entity/%s", nrPlatformHostname(), entityGUID)
}

// nrPlatformHostname returns the host for the platform based on the region set.
func nrPlatformHostname() string {
	r := configuration.GetActiveProfileString(configuration.Region)
	if strings.EqualFold(r, region.Staging.String()) {
		return nrPlatformHostnames.Staging
	}

	if strings.EqualFold(r, region.EU.String()) {
		return nrPlatformHostnames.EU
	}

	return nrPlatformHostnames.US
}
