//go:build unit

package accessmanagement

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestPermissionsGet(t *testing.T) {
	assert.Equal(t, "get", cmdPermissionsGet.Name())
	testcobra.CheckCobraMetadata(t, cmdPermissionsGet)
	testcobra.CheckCobraRequiredFlags(t, cmdPermissionsGet, []string{})
}
