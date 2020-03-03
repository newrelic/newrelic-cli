package apm

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestAPMCommand(t *testing.T) {
	assert.NotEmptyf(t, Command.Use, "Need to set Command.%s on Command %s", "Use", Command.CalledAs())
	assert.NotEmptyf(t, Command.Short, "Need to set Command.%s on Command %s", "Short", Command.CalledAs())

	for _, c := range Command.Commands() {
		assert.NotEmptyf(t, c.Use, "Need to set Command.%s on Command %s", "Use", c.CommandPath())
		assert.NotEmptyf(t, c.Short, "Need to set Command.%s on Command %s", "Short", c.CommandPath())
		assert.NotEmptyf(t, c.Long, "Need to set Command.%s on Command %s", "Long", c.CommandPath())
		assert.NotEmptyf(t, c.Example, "Need to set Command.%s on Command %s", "Example", c.CommandPath())
	}
}

func TestApmDescribeDeployments(t *testing.T) {
	assert.Equal(t, "describe-deployments", apmDescribeDeployments.Name())

	requiredFlags := []string{"applicationId"}

	for _, r := range requiredFlags {
		x := apmDescribeDeployments.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

func TestApmCreateDeployment(t *testing.T) {
	command := apmCreateDeployment
	assert.Equal(t, "create-deployment", command.Name())

	requiredFlags := []string{"applicationId", "revision"}

	for _, r := range requiredFlags {
		x := command.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

func TestApmDeleteDeployment(t *testing.T) {
	command := apmDeleteDeployment
	assert.Equal(t, "delete-deployment", command.Name())

	requiredFlags := []string{"applicationId", "deploymentID"}

	for _, r := range requiredFlags {
		x := command.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

func TestApmApplication(t *testing.T) {
	command := apmApplication
	assert.Equal(t, "application", command.Name())

	requiredFlags := []string{}

	for _, r := range requiredFlags {
		x := command.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

func TestApmGetApplication(t *testing.T) {
	command := apmGetApplication
	assert.Equal(t, "get", command.Name())

	requiredFlags := []string{}

	for _, r := range requiredFlags {
		x := command.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}
