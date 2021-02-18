// +build unit

package execution

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockStatusReporter_interface(t *testing.T) {
	var r StatusSubscriber = NewMockStatusReporter()
	require.NotNil(t, r)
}
