//go:build unit
// +build unit

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestTerraform(t *testing.T) {
	assert.Equal(t, "terraform", cmdTerraform.Name())

	testcobra.CheckCobraMetadata(t, cmdSemver)
}

func TestTerraformDashboard(t *testing.T) {
	assert.Equal(t, "dashboard", cmdTerraformDashboard.Name())

	testcobra.CheckCobraMetadata(t, cmdTerraformDashboard)
}

func TestLabelContainingNumbers(t *testing.T) {
	label = "label_1"
	assert.NoError(t, cmdTerraformDashboard.Args(nil, nil))

	label = "label-1"
	assert.Error(t, cmdTerraformDashboard.Args(nil, nil))

	label = "1_label"
	assert.Error(t, cmdTerraformDashboard.Args(nil, nil))
}
