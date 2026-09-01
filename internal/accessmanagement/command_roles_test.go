//go:build unit

package accessmanagement

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestRolesGet(t *testing.T) {
	assert.Equal(t, "get", cmdRolesGet.Name())
	testcobra.CheckCobraMetadata(t, cmdRolesGet)
	testcobra.CheckCobraRequiredFlags(t, cmdRolesGet, []string{})
}
