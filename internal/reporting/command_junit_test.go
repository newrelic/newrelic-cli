// +build unit

package reporting

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestJUnit(t *testing.T) {
	assert.Equal(t, "junit", cmdJUnit.Name())

	testcobra.CheckCobraMetadata(t, cmdJUnit)
}
