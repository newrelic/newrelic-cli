//go:build unit

package usermanagement

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestAuthDomainsGet(t *testing.T) {
	assert.Equal(t, "get", cmdAuthDomainsGet.Name())
	testcobra.CheckCobraMetadata(t, cmdAuthDomainsGet)
	testcobra.CheckCobraRequiredFlags(t, cmdAuthDomainsGet, []string{})
}
