package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/agent/version"
	"github.com/newrelic/newrelic-client-go/v2/pkg/agent"
)

func TestAgentNameTitleCase(t *testing.T) {
	t.Parallel()

	expected := "Android"
	agentName := agent.AgentReleasesFilterTypes.ANDROID
	actual := version.AgentNameTitleCase(agentName)

	assert.Equal(t, expected, actual)

	expected = "Browser"
	agentName = agent.AgentReleasesFilterTypes.BROWSER
	actual = version.AgentNameTitleCase(agentName)

	assert.Equal(t, expected, actual)

	expected = ".NET"
	agentName = agent.AgentReleasesFilterTypes.DOTNET
	actual = version.AgentNameTitleCase(agentName)

	assert.Equal(t, expected, actual)

	expected = "Elixir"
	agentName = agent.AgentReleasesFilterTypes.ELIXIR
	actual = version.AgentNameTitleCase(agentName)

	assert.Equal(t, expected, actual)

	expected = "Go"
	agentName = agent.AgentReleasesFilterTypes.GO
	actual = version.AgentNameTitleCase(agentName)

	assert.Equal(t, expected, actual)

	expected = "Infrastructure"
	agentName = agent.AgentReleasesFilterTypes.INFRASTRUCTURE
	actual = version.AgentNameTitleCase(agentName)

	assert.Equal(t, expected, actual)

	expected = "iOS"
	agentName = agent.AgentReleasesFilterTypes.IOS
	actual = version.AgentNameTitleCase(agentName)

	assert.Equal(t, expected, actual)

	expected = "Java"
	agentName = agent.AgentReleasesFilterTypes.JAVA
	actual = version.AgentNameTitleCase(agentName)

	assert.Equal(t, expected, actual)

	expected = "Node.js"
	agentName = agent.AgentReleasesFilterTypes.NODEJS
	actual = version.AgentNameTitleCase(agentName)

	assert.Equal(t, expected, actual)

	expected = "PHP"
	agentName = agent.AgentReleasesFilterTypes.PHP
	actual = version.AgentNameTitleCase(agentName)

	assert.Equal(t, expected, actual)

	expected = "Python"
	agentName = agent.AgentReleasesFilterTypes.PYTHON
	actual = version.AgentNameTitleCase(agentName)

	assert.Equal(t, expected, actual)

	expected = "Ruby"
	agentName = agent.AgentReleasesFilterTypes.RUBY
	actual = version.AgentNameTitleCase(agentName)

	assert.Equal(t, expected, actual)

	expected = "C SDK"
	agentName = agent.AgentReleasesFilterTypes.SDK
	actual = version.AgentNameTitleCase(agentName)

	assert.Equal(t, expected, actual)
}

func TestIsValidAgentName(t *testing.T) {
	t.Parallel()

	agentName := "PYTHON"
	actual := version.IsValidAgentName(agentName)

	assert.True(t, actual)

	agentName = "ASDF"
	actual = version.IsValidAgentName(agentName)

	assert.False(t, actual)
}
