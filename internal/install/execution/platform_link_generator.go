package execution

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
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

func (g *PlatformLinkGenerator) GenerateExplorerLink(status InstallStatus) string {
	return generateExplorerLink(status)
}

func (g *PlatformLinkGenerator) GenerateEntityLink(entityGUID string) string {
	return generateEntityLink(entityGUID)
}

// GenerateRedirectURL creates a URL for the user to navigate to after running
// through an installation. The URL is displayed in the CLI out as well and is
// also provided in the nerdstorage document. This provides the user two options
// to see their data - click from the CLI output or from the frontend.
func (g *PlatformLinkGenerator) GenerateRedirectURL(status InstallStatus) string {
	if status.AllSelectedRecipesInstalled() {
		return g.GenerateEntityLink(status.HostEntityGUID())
	}

	if status.hasAnyRecipeStatus(RecipeStatusTypes.INSTALLED) {
		return g.GenerateExplorerLink(status)
	}

	return "" // g.GenerateExplorerLink("")
}

type referrerParamValue struct {
	NerdletID  string `json:"nerdletId,omitempty"`
	Referrer   string `json:"referrer,omitempty"`
	EntityGUID string `json:"entityGuid,omitempty"`
}

// The CLI URL referrer param is a JSON string containing information
// the UI can use to understand how/where the URL was generated. This allows the
// UI to return to its previous state in the case of a user closing the browser
// and then clicking a redirect URL in the CLI's output.
func generateReferrerParam(entityGUID string) string {
	p := referrerParamValue{
		NerdletID: "nr1-install-newrelic.installation-plan",
		Referrer:  "newrelic-cli",
	}

	if entityGUID != "" {
		p.EntityGUID = entityGUID
	}

	stringifiedParam, err := json.Marshal(p)
	if err != nil {
		// TODO: add debug log
		return ""
	}

	return string(stringifiedParam)
}

func generateExplorerLink(status InstallStatus) string {
	return fmt.Sprintf("https://%s/launcher/nr1-core.explorer?platform[filters]=%s&platform[accountId]=%d&cards[0]=%s",
		nrPlatformHostname(),
		utils.Base64Encode(status.successLinkConfig.Filter),
		configAPI.GetActiveProfileAccountID(),
		utils.Base64Encode(generateReferrerParam(status.HostEntityGUID())),
	)
}

func generateEntityLink(entityGUID string) string {
	return fmt.Sprintf("https://%s/redirect/entity/%s", nrPlatformHostname(), entityGUID)
}

// nrPlatformHostname returns the host for the platform based on the region set.
func nrPlatformHostname() string {
	r := configAPI.GetActiveProfileString(config.Region)
	if strings.EqualFold(r, region.Staging.String()) {
		return nrPlatformHostnames.Staging
	}

	if strings.EqualFold(r, region.EU.String()) {
		return nrPlatformHostnames.EU
	}

	return nrPlatformHostnames.US
}
