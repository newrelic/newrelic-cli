//+ build unit
package install

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewExecutionStatus(t *testing.T) {
	s := newExecutionStatus()
	require.NotEqual(t, 0, s.Timestamp)
}
