// +build unit

package execution

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTerminalStatusReporter_interface(t *testing.T) {
	var r StatusSubscriber = NewTerminalStatusReporter()
	require.NotNil(t, r)
}
