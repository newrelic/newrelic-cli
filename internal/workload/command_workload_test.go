// +build unit

package workload

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestCreate(t *testing.T) {
	assert.Equal(t, "create", cmdCreate.Name())

	testcobra.CheckCobraMetadata(t, cmdCreate)
	testcobra.CheckCobraRequiredFlags(t, cmdCreate, []string{})
}

func TestUpdate(t *testing.T) {
	assert.Equal(t, "update", cmdUpdate.Name())

	testcobra.CheckCobraMetadata(t, cmdUpdate)
	testcobra.CheckCobraRequiredFlags(t, cmdUpdate, []string{})
}

func TestDuplicate(t *testing.T) {
	assert.Equal(t, "duplicate", cmdDuplicate.Name())

	testcobra.CheckCobraMetadata(t, cmdDuplicate)
	testcobra.CheckCobraRequiredFlags(t, cmdDuplicate, []string{})
}

func TestDelete(t *testing.T) {
	assert.Equal(t, "delete", cmdDelete.Name())

	testcobra.CheckCobraMetadata(t, cmdDelete)
	testcobra.CheckCobraRequiredFlags(t, cmdDelete, []string{})
}
