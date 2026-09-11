//go:build integration

package profile

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

// clearAmbientCredentialEnvVars unsets NEW_RELIC_* env vars so a developer's
// shell can't leak an override into these tests. Must fully unset, not set
// to "" - os.LookupEnv treats present-but-empty as a real override too.
func clearAmbientCredentialEnvVars(t *testing.T) {
	t.Helper()
	for _, name := range []string{"NEW_RELIC_ACCOUNT_ID", "NEW_RELIC_API_KEY", "NEW_RELIC_LICENSE_KEY", "NEW_RELIC_REGION"} {
		name := name
		original, wasSet := os.LookupEnv(name)
		require.NoError(t, os.Unsetenv(name))
		if wasSet {
			t.Cleanup(func() { os.Setenv(name, original) })
		}
	}
}

// resetListFormatFlagState clears stale Changed state on cmdList's merged
// "format" flag. cmdList is a package-level singleton reused across tests,
// and cobra caches a parent's persistent flags into it on first Execute() -
// a later test's fresh root command has its "format" flag silently ignored
// in favor of that cached one, Changed value included.
func resetListFormatFlagState() {
	if f := cmdList.Flags().Lookup("format"); f != nil {
		f.Changed = false
	}
}

// buildTestRootCmd replicates the real root command's persistent --format
// flag (cmd/newrelic/command.go) with the `profile` tree attached beneath
// it, so cmd.Flags().Changed("format") behaves as it does in production -
// a bare *cobra.Command never has --format registered.
func buildTestRootCmd() *cobra.Command {
	var outputFormat string

	root := &cobra.Command{
		Use: "newrelic",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			utils.LogIfError(output.SetFormat(output.ParseFormat(outputFormat)))
		},
	}
	root.PersistentFlags().StringVar(&outputFormat, "format", output.DefaultFormat.String(), "")
	root.PersistentFlags().StringVar(&config.FlagProfileName, "profile", "", "")
	root.AddCommand(Command)

	return root
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	fn()

	require.NoError(t, w.Close())
	os.Stdout = old

	out, err := io.ReadAll(r)
	require.NoError(t, err)

	return string(out)
}

func setupTwoProfilesWithADefault(t *testing.T) {
	t.Helper()
	clearAmbientCredentialEnvVars(t)

	config.Init(t.TempDir())
	t.Cleanup(func() { config.FlagProfileName = "" })

	require.NoError(t, configAPI.SetProfileValue("A", config.APIKey, "key-A"))
	require.NoError(t, configAPI.SetProfileValue("A", config.Region, "US"))
	require.NoError(t, configAPI.SetProfileValue("B", config.APIKey, "key-B"))
	require.NoError(t, configAPI.SetProfileValue("B", config.Region, "EU"))
	require.NoError(t, configAPI.SetDefaultProfile("A"))
}

// Skipping the Account ID prompt must leave the field unset, not written as
// 0 (which used to fail validation and abort `profile add`).
func TestAddIntValueToProfile_SkipsAccountIDWhenZero(t *testing.T) {
	clearAmbientCredentialEnvVars(t)
	config.Init(t.TempDir())

	acceptDefaults = true
	defer func() { acceptDefaults = false }()

	addIntValueToProfile("no-account-id", 0, config.AccountID, "Account ID", nil)

	_, err := config.CredentialsProvider.GetIntWithScope("no-account-id", config.AccountID)
	require.Error(t, err, "accountID must remain unset (no value at all), not written as 0")
}

// Sanity check that the skip-on-zero path didn't disable the normal set path.
func TestAddIntValueToProfile_SetsNonZeroValue(t *testing.T) {
	clearAmbientCredentialEnvVars(t)
	config.Init(t.TempDir())

	addIntValueToProfile("with-account-id", 12345, config.AccountID, "Account ID", nil)

	v, err := config.CredentialsProvider.GetIntWithScope("with-account-id", config.AccountID)
	require.NoError(t, err)
	require.Equal(t, int64(12345), v)
}

// `--profile B` selects which profile's credentials to use, but must not
// make `list` mislabel B as the persisted default - only A may be marked so.
func TestCmdList_MarksTruePersistedDefault_RegardlessOfProfileFlag(t *testing.T) {
	setupTwoProfilesWithADefault(t)
	config.FlagProfileName = "B"

	root := buildTestRootCmd()
	resetListFormatFlagState()
	out := captureStdout(t, func() {
		root.SetArgs([]string{"profile", "list", "--profile", "B", "--format", "json"})
		require.NoError(t, root.Execute())
	})

	var rows []map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &rows))

	byName := map[string]map[string]interface{}{}
	for _, r := range rows {
		byName[r["Name"].(string)] = r
	}

	require.Equal(t, true, byName["A"]["isDefault"], "A is the persisted default and must be marked so even with --profile B active")
	require.Equal(t, false, byName["B"]["isDefault"])
}

// With no --format flag, `list` must still render a table, not fall through
// to output.Print()'s own default of JSON.
func TestCmdList_DefaultsToTableFormat_WhenFormatFlagNotPassed(t *testing.T) {
	setupTwoProfilesWithADefault(t)

	root := buildTestRootCmd()
	resetListFormatFlagState()
	out := captureStdout(t, func() {
		root.SetArgs([]string{"profile", "list"})
		require.NoError(t, root.Execute())
	})

	require.NotEmpty(t, out)
	require.NotEqual(t, byte('['), out[0], "no --format flag was passed; output must not be JSON")
	require.Contains(t, out, "Name")
	require.Contains(t, out, "isDefault")
}

// An explicit --format json must take effect and contain no ANSI codes -
// not in the Name field's "(default)" suffix, nor in obfuscated values.
func TestCmdList_RespectsExplicitJSONFormat_CleanOfANSI(t *testing.T) {
	setupTwoProfilesWithADefault(t)

	root := buildTestRootCmd()
	resetListFormatFlagState()
	out := captureStdout(t, func() {
		root.SetArgs([]string{"profile", "list", "--format", "json"})
		require.NoError(t, root.Execute())
	})

	require.NotContains(t, out, "\x1b", "JSON output must contain no ANSI escape codes")

	var rows []map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &rows))
	require.Len(t, rows, 2)

	byName := map[string]map[string]interface{}{}
	for _, r := range rows {
		byName[r["Name"].(string)] = r
	}

	require.Equal(t, "A", byName["A"]["Name"], "Name must be the clean profile name, with no embedded ANSI suffix")
	require.Equal(t, true, byName["A"]["isDefault"])
}
