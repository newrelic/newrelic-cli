//go:build unit
// +build unit

package events

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestPost(t *testing.T) {
	assert.Equal(t, "post", cmdPost.Name())

	testcobra.CheckCobraMetadata(t, cmdPost)
	testcobra.CheckCobraRequiredFlags(t, cmdPost, []string{})
}
