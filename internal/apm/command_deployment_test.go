// +build unit

package apm

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

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
