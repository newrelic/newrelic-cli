package execution

import (
	"fmt"
	"strings"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/pkg/region"
)

type SuccessLinkGenerator interface {
	GenerateExplorerLink(filter string) string
	GenerateEntityLink(entityGUID string) string
}

type ConcreteSuccessLinkGenerator struct{}

var nrPlatformHostnames = struct {
	Staging string
	US      string
	EU      string
}{
	Staging: "staging-one.newrelic.com",
	US:      "one.newrelic.com",
	EU:      "one.eu.newrelic.com",
}

func NewConcreteSuccessLinkGenerator() *ConcreteSuccessLinkGenerator {
	return &ConcreteSuccessLinkGenerator{}
}

func (g *ConcreteSuccessLinkGenerator) GenerateExplorerLink(filter string) string {
	return generateExplorerLink(filter)
}

func (g *ConcreteSuccessLinkGenerator) GenerateEntityLink(entityGUID string) string {
	return generateEntityLink(entityGUID)
}

func generateSuccessURL(status InstallStatus) string {
	if status.hasAnyRecipeStatus(RecipeStatusTypes.INSTALLED) {
		switch t := status.successLinkConfig.Type; {
		case strings.EqualFold(string(t), "explorer"):
			return generateExplorerLink(status.successLinkConfig.Filter)
		default:
			return generateEntityLink(status.HostEntityGUID())
		}
	}

	return ""
}

func generateExplorerLink(filter string) string {
	return fmt.Sprintf("https://%s/launcher/nr1-core.explorer?platform[filters]=%s&platform[accountId]=%d",
		nrPlatformHostname(),
		utils.Base64Encode(filter),
		credentials.DefaultProfile().AccountID,
	)
}

func generateEntityLink(entityGUID string) string {
	return fmt.Sprintf("https://%s/redirect/entity/%s", nrPlatformHostname(), entityGUID)
}

// nrPlatformHostname returns the host for the platform based on the region set.
func nrPlatformHostname() string {
	defaultProfile := credentials.DefaultProfile()
	if defaultProfile == nil {
		return nrPlatformHostnames.US
	}

	if strings.EqualFold(defaultProfile.Region, region.Staging.String()) {
		return nrPlatformHostnames.Staging
	}

	if strings.EqualFold(defaultProfile.Region, region.EU.String()) {
		return nrPlatformHostnames.EU
	}

	return nrPlatformHostnames.US
}
