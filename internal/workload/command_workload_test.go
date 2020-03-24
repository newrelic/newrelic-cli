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
