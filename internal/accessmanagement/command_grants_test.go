//go:build unit

package accessmanagement

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestGrantsGet(t *testing.T) {
	assert.Equal(t, "get", cmdGrantsGet.Name())
	testcobra.CheckCobraMetadata(t, cmdGrantsGet)
	testcobra.CheckCobraRequiredFlags(t, cmdGrantsGet, []string{})
}

func TestGrantsCreate(t *testing.T) {
	assert.Equal(t, "create", cmdGrantsCreate.Name())
	testcobra.CheckCobraMetadata(t, cmdGrantsCreate)
	testcobra.CheckCobraRequiredFlags(t, cmdGrantsCreate, []string{"groupId", "roleId", "scope"})
}

func TestGrantsRevoke(t *testing.T) {
	assert.Equal(t, "revoke", cmdGrantsRevoke.Name())
	testcobra.CheckCobraMetadata(t, cmdGrantsRevoke)
	testcobra.CheckCobraRequiredFlags(t, cmdGrantsRevoke, []string{"groupId", "roleId", "scope"})
}
