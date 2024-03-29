//go:build unit
// +build unit

package nerdstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestDocumentWrite(t *testing.T) {
	assert.Equal(t, "write", cmdDocumentWrite.Name())

	testcobra.CheckCobraMetadata(t, cmdDocumentWrite)
	testcobra.CheckCobraRequiredFlags(t, cmdDocumentWrite, []string{"packageId", "scope", "collection", "documentId", "document"})
}

func TestDocumentGet(t *testing.T) {
	assert.Equal(t, "get", cmdDocumentGet.Name())

	testcobra.CheckCobraMetadata(t, cmdDocumentGet)
	testcobra.CheckCobraRequiredFlags(t, cmdDocumentGet, []string{"packageId", "scope", "collection", "documentId"})
}

func TestDocumentDelete(t *testing.T) {
	assert.Equal(t, "delete", cmdDocumentDelete.Name())

	testcobra.CheckCobraMetadata(t, cmdDocumentDelete)
	testcobra.CheckCobraRequiredFlags(t, cmdDocumentDelete, []string{"packageId", "scope", "collection", "documentId"})
}
