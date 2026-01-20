package version

import (
	"reflect"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/newrelic/newrelic-client-go/v2/pkg/agent"
)

func AgentNameTitleCase(agentName agent.AgentReleasesFilter) string {
	caser := cases.Title(language.AmericanEnglish)

	switch agentName {
	case agent.AgentReleasesFilterTypes.DOTNET:
		return ".NET"
	case agent.AgentReleasesFilterTypes.IOS:
		return "iOS"
	case agent.AgentReleasesFilterTypes.NODEJS:
		return "Node.js"
	case agent.AgentReleasesFilterTypes.PHP:
		return "PHP"
	case agent.AgentReleasesFilterTypes.SDK:
		return "C SDK"
	default:
		return caser.String(string(agentName))
	}
}

func IsValidAgentName(agentName string) bool {
	a := reflect.ValueOf(agent.AgentReleasesFilterTypes)

	return a.FieldByName(agentName).String() == agentName
}
