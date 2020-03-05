// +build unit

package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntitiesSearch(t *testing.T) {
	command := entitiesSearch

	assert.Equal(t, "search", command.Name())
	assert.True(t, command.HasFlags())
}
