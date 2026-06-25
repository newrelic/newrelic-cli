//go:build unit

package usermanagement

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestUsersGet(t *testing.T) {
	assert.Equal(t, "get", cmdUsersGet.Name())
	testcobra.CheckCobraMetadata(t, cmdUsersGet)
	testcobra.CheckCobraRequiredFlags(t, cmdUsersGet, []string{"authDomainId"})
}

func TestUsersCreate(t *testing.T) {
	assert.Equal(t, "create", cmdUsersCreate.Name())
	testcobra.CheckCobraMetadata(t, cmdUsersCreate)
	testcobra.CheckCobraRequiredFlags(t, cmdUsersCreate, []string{"authDomainId", "email", "name"})
}

func TestUsersUpdate(t *testing.T) {
	assert.Equal(t, "update", cmdUsersUpdate.Name())
	testcobra.CheckCobraMetadata(t, cmdUsersUpdate)
	testcobra.CheckCobraRequiredFlags(t, cmdUsersUpdate, []string{"id"})
}

func TestUsersDelete(t *testing.T) {
	assert.Equal(t, "delete", cmdUsersDelete.Name())
	testcobra.CheckCobraMetadata(t, cmdUsersDelete)
	testcobra.CheckCobraRequiredFlags(t, cmdUsersDelete, []string{"id"})
}
