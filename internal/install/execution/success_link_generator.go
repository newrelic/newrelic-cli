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

func NewConcreteSuccessLinkGenerator() *ConcreteSuccessLinkGenerator {
	return &ConcreteSuccessLinkGenerator{}
}

func (g *ConcreteSuccessLinkGenerator) GenerateExplorerLink(filter string) string {
	return fmt.Sprintf("https://%s/launcher/nr1-core.explorer?platform[filters]=%s&platform[accountId]=%d",
		nrPlatformHostname(),
		utils.Base64Encode(filter),
		credentials.DefaultProfile().AccountID)
}

func (g *ConcreteSuccessLinkGenerator) GenerateEntityLink(entityGUID string) string {
	return fmt.Sprintf("https://one.newrelic.com/redirect/entity/%s", entityGUID)
}

// nrPlatformHostname returns the host for the platform based on the region set.
func nrPlatformHostname() string {
	switch defaultProfile := credentials.DefaultProfile(); {
	case strings.EqualFold(defaultProfile.Region, region.Staging.String()):
		return "staging-one.newrelic.com"
	case strings.EqualFold(defaultProfile.Region, region.EU.String()):
		return "one.eu.newrelic.com"
	default:
		return "one.newrelic.com"
	}
}
