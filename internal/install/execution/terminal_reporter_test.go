// +build unit

package execution

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTerminalStatusReporter_interface(t *testing.T) {
	var r StatusReporter = NewTerminalStatusReporter()
	require.NotNil(t, r)
}
