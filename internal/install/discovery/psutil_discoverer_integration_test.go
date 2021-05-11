// +build integration

package discovery

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestDiscovery(t *testing.T) {
	cmd := exec.Command("java", "-classpath", "internal/install/mockProcesses", "JavaDaemonTest")
	if err := cmd.Start(); err != nil {
		t.Fatalf("error starting java process")
	}

	mockRecipeFetcher := recipes.NewMockRecipeFetcher()
	mockRecipeFetcher.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			ID:           "test",
			Name:         "java",
			ProcessMatch: []string{"java"},
		},
	}

	pf := NewRegexProcessFilterer(mockRecipeFetcher)
	pd := NewPSUtilDiscoverer(pf)

	manifest, err := pd.Discover(context.Background())

	require.NoError(t, err)
	require.NotNil(t, manifest)
	require.GreaterOrEqual(t, len(manifest.Processes), 1)

	err = cmd.Process.Signal(os.Interrupt)
	if err != nil {
		t.Fatalf("error sending interrupt to java process: %s", err)
	}

	_ = cmd.Wait()
}
