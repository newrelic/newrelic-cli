//go:build unit
// +build unit

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestSemver(t *testing.T) {
	assert.Equal(t, "semver", cmdSemver.Name())

	testcobra.CheckCobraMetadata(t, cmdSemver)
}

func TestSemverCheck(t *testing.T) {
	assert.Equal(t, "check", cmdSemverCheck.Name())

	testcobra.CheckCobraMetadata(t, cmdSemverCheck)
	testcobra.CheckCobraRequiredFlags(t, cmdSemverCheck, []string{"constraint", "version"})
}
