package entities

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestEntitiesDescribeTags(t *testing.T) {
	r := entitiesDescribeTags.Flag("guid")
	assert.Equal(t, []string{"true"}, r.Annotations[cobra.BashCompOneRequiredFlag])
}
