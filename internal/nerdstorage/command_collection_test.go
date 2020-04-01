// +build unit

package nerdstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestCollectionGet(t *testing.T) {
	assert.Equal(t, "get", cmdCollectionGet.Name())

	testcobra.CheckCobraMetadata(t, cmdCollectionGet)
	testcobra.CheckCobraRequiredFlags(t, cmdCollectionGet, []string{"packageId", "scope", "collection"})
}

func TestCollectionDelete(t *testing.T) {
	assert.Equal(t, "delete", cmdCollectionDelete.Name())

	testcobra.CheckCobraMetadata(t, cmdCollectionDelete)
	testcobra.CheckCobraRequiredFlags(t, cmdCollectionDelete, []string{"packageId", "scope", "collection"})
}
