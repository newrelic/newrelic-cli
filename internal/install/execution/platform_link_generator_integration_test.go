//go:build integration
// +build integration

package execution

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateShortNewRelicURL(t *testing.T) {
	t.Parallel()

	skipTestIfMissingEnvVars(t, []string{"NEW_RELIC_API_KEY"})

	g := NewPlatformLinkGenerator()

	longURL := "https://one.newrelic.com/launcher/nr1-core.explorer?pane=eyJuZXJkbGV0SWQiOiJucjEtY29yZS5saXN0aW5nIiwiZmF2b3JpdGVzIjp7InNlbGVjdGVkIjp0cnVlLCJ2aXNpYmxlIjp0cnVlfSwibGFzdFZpZXdlZCI6eyJzZWxlY3RlZCI6ZmFsc2UsInZpc2libGUiOnRydWV9fQ==&sidebars[0]=eyJuZXJkbGV0SWQiOiJucjEtY29yZS5jYXRlZ29yaWVzIiwicm9vdE5lcmRsZXRJZCI6Im5yMS1jb3JlLmxpc3RpbmciLCJmYXZvcml0ZXMiOnsic2VsZWN0ZWQiOnRydWUsInZpc2libGUiOnRydWV9LCJsYXN0Vmlld2VkIjp7InNlbGVjdGVkIjpmYWxzZSwidmlzaWJsZSI6dHJ1ZX19&state=63f7dbc9-1b4f-3610-57ea-45e2c0908345"
	result, err := g.generateShortNewRelicURL(longURL)

	require.NoError(t, err)
	require.Less(t, len(result), len(longURL))
}

func TestGenerateShortNewRelicURL_NoAPIKey(t *testing.T) {
	t.Parallel()

	g := NewPlatformLinkGenerator()
	g.apiKey = "" // unset the API key so an error is returned from the API

	longURL := "https://one.newrelic.com/launcher/nr1-core.explorer?pane=eyJuZXJkbGV0SWQiOiJucjEtY29yZS5saXN0aW5nIiwiZmF2b3JpdGVzIjp7InNlbGVjdGVkIjp0cnVlLCJ2aXNpYmxlIjp0cnVlfSwibGFzdFZpZXdlZCI6eyJzZWxlY3RlZCI6ZmFsc2UsInZpc2libGUiOnRydWV9fQ==&sidebars[0]=eyJuZXJkbGV0SWQiOiJucjEtY29yZS5jYXRlZ29yaWVzIiwicm9vdE5lcmRsZXRJZCI6Im5yMS1jb3JlLmxpc3RpbmciLCJmYXZvcml0ZXMiOnsic2VsZWN0ZWQiOnRydWUsInZpc2libGUiOnRydWV9LCJsYXN0Vmlld2VkIjp7InNlbGVjdGVkIjpmYWxzZSwidmlzaWJsZSI6dHJ1ZX19&state=63f7dbc9-1b4f-3610-57ea-45e2c0908345"
	result, err := g.generateShortNewRelicURL(longURL)

	require.Error(t, err)
	require.Equal(t, result, longURL)
}

func TestGenerateShortNewRelicURL_InvalidAPIKey(t *testing.T) {
	t.Parallel()

	g := NewPlatformLinkGenerator()
	g.apiKey = "abc123" // invalid API key
	longURL := "https://one.newrelic.com/launcher/nr1-core.explorer?pane=eyJuZXJkbGV0SWQiOiJucjEtY29yZS5saXN0aW5nIiwiZmF2b3JpdGVzIjp7InNlbGVjdGVkIjp0cnVlLCJ2aXNpYmxlIjp0cnVlfSwibGFzdFZpZXdlZCI6eyJzZWxlY3RlZCI6ZmFsc2UsInZpc2libGUiOnRydWV9fQ==&sidebars[0]=eyJuZXJkbGV0SWQiOiJucjEtY29yZS5jYXRlZ29yaWVzIiwicm9vdE5lcmRsZXRJZCI6Im5yMS1jb3JlLmxpc3RpbmciLCJmYXZvcml0ZXMiOnsic2VsZWN0ZWQiOnRydWUsInZpc2libGUiOnRydWV9LCJsYXN0Vmlld2VkIjp7InNlbGVjdGVkIjpmYWxzZSwidmlzaWJsZSI6dHJ1ZX19&state=63f7dbc9-1b4f-3610-57ea-45e2c0908345"
	result, err := g.generateShortNewRelicURL(longURL)

	require.Error(t, err)
	require.Equal(t, result, longURL)
}

func TestGenerateShortNewRelicURL_NonNewRelicURL(t *testing.T) {
	t.Parallel()

	g := NewPlatformLinkGenerator()
	longURL := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	result, _ := g.generateShortNewRelicURL(longURL)

	require.Equal(t, result, longURL)
}

func skipTestIfMissingEnvVars(t *testing.T, vars []string) {
	for _, v := range vars {
		if os.Getenv(v) == "" {
			t.Skipf("Skipping test. The following environment variables are required to run this test: %s", strings.Join(vars[:], ", "))
		}
	}
}
