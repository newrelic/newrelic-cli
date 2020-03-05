// +build unit

package apm

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

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
