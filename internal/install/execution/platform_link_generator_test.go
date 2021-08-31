//go:build unit
// +build unit

package execution

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestGenerateExplorerLink(t *testing.T) {
	g := NewPlatformLinkGenerator()

	expectedEncodedQueryParamSubstring := utils.Base64Encode(cliURLReferrerParam)

	installStatus := InstallStatus{}

	result := g.GenerateExplorerLink(installStatus)

	require.Contains(t, result, expectedEncodedQueryParamSubstring)
}
