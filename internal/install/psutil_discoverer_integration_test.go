// +build integration

package install

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiscovery(t *testing.T) {
	cmd := exec.Command("java", "-classpath", "internal/install/mockProcesses", "JavaDaemonTest")
	if err := cmd.Start(); err != nil {
		t.Fatalf("error starting java process")
	}

	mockRecipeFetcher := newMockRecipeFetcher()
	mockRecipeFetcher.fetchRecipesVal = []recipe{
		{
			ID:           "test",
			Name:         "java",
			ProcessMatch: []string{"java"},
		},
	}

	pf := newRegexProcessFilterer(mockRecipeFetcher)
	pd := newPSUtilDiscoverer(pf)

	manifest, err := pd.discover(context.Background())

	require.NoError(t, err)
	require.NotNil(t, manifest)
	require.GreaterOrEqual(t, len(manifest.Processes), 1)

	err = cmd.Process.Signal(os.Interrupt)
	if err != nil {
		t.Fatalf("error sending interrupt to java process: %s", err)
	}

	_ = cmd.Wait()
}
