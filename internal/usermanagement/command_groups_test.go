//go:build unit

package usermanagement

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestGroupsGet(t *testing.T) {
	assert.Equal(t, "get", cmdGroupsGet.Name())
	testcobra.CheckCobraMetadata(t, cmdGroupsGet)
	testcobra.CheckCobraRequiredFlags(t, cmdGroupsGet, []string{"authDomainId"})
}

func TestGroupsCreate(t *testing.T) {
	assert.Equal(t, "create", cmdGroupsCreate.Name())
	testcobra.CheckCobraMetadata(t, cmdGroupsCreate)
	testcobra.CheckCobraRequiredFlags(t, cmdGroupsCreate, []string{"authDomainId", "name"})
}

func TestGroupsUpdate(t *testing.T) {
	assert.Equal(t, "update", cmdGroupsUpdate.Name())
	testcobra.CheckCobraMetadata(t, cmdGroupsUpdate)
	testcobra.CheckCobraRequiredFlags(t, cmdGroupsUpdate, []string{"id", "name"})
}

func TestGroupsDelete(t *testing.T) {
	assert.Equal(t, "delete", cmdGroupsDelete.Name())
	testcobra.CheckCobraMetadata(t, cmdGroupsDelete)
	testcobra.CheckCobraRequiredFlags(t, cmdGroupsDelete, []string{"id"})
}

func TestGroupsMembersAdd(t *testing.T) {
	assert.Equal(t, "add", cmdGroupsMembersAdd.Name())
	testcobra.CheckCobraMetadata(t, cmdGroupsMembersAdd)
	testcobra.CheckCobraRequiredFlags(t, cmdGroupsMembersAdd, []string{"groupId", "userId"})
}

func TestGroupsMembersRemove(t *testing.T) {
	assert.Equal(t, "remove", cmdGroupsMembersRemove.Name())
	testcobra.CheckCobraMetadata(t, cmdGroupsMembersRemove)
	testcobra.CheckCobraRequiredFlags(t, cmdGroupsMembersRemove, []string{"groupId", "userId"})
}
