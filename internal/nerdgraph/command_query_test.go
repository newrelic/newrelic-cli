// +build unit

package nerdgraph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	command := queryCmd

	assert.Equal(t, "query", command.Name())
	assert.True(t, command.HasFlags())
}
