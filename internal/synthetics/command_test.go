//go:build unit
// +build unit

package synthetics

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestSyntheticsCommand(t *testing.T) {
	testcobra.CheckCobraMetadata(t, cmdMon)
}
