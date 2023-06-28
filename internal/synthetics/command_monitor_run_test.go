//go:build unit

package synthetics

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCmdSyntheticsMonitorRun(t *testing.T) {
	cmd := &cobra.Command{
		Use:  cmdMonitorRun.Use,
		RunE: execCmdMonitorRunE,
	}

	err := cmd.Execute()

	assert.NoError(t, err)
}
