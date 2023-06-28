//go:build unit

package utils

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCmdDo(t *testing.T) {
	cmd := &cobra.Command{
		Use:  cmdDo.Use,
		RunE: runDoCommandE,
	}

	err := cmd.Execute()

	assert.NoError(t, err)
}
