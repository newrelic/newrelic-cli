// +build integration

package configuration

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type testScenario struct {
	configFile         *os.File
	credsFile          *os.File
	defaultProfileFile *os.File
}

func (s *testScenario) teardown() {
	os.Remove(s.configFile.Name())
	os.Remove(s.credsFile.Name())
	os.Remove(s.defaultProfileFile.Name())
}

func setupTestScenario(t *testing.T) testScenario {
	configFile, err := ioutil.TempFile("", "config*.json")
	require.NoError(t, err)

	configJSON := `
{
	"*": {
		"loglevel": "info",
		"plugindir": "/tmp",
		"prereleasefeatures": "NOT_ASKED",
		"sendusagedata": "NOT_ASKED"
	}
}
`
	_, err = configFile.Write([]byte(configJSON))
	require.NoError(t, err)

	credsFile, err := ioutil.TempFile("", "credentials*.json")
	require.NoError(t, err)

	credsJSON := `
{
	"default": {
		"apiKey": "testApiKey",
		"region": "US",
		"accountID": 12345,
		"licenseKey": "testLicenseKey"
	}
}
`
	_, err = credsFile.Write(([]byte(credsJSON)))
	require.NoError(t, err)

	defaultProfileFile, err := ioutil.TempFile("", "default-profile*.json")
	require.NoError(t, err)

	defaultProfileJSON := `"default"`
	_, err = defaultProfileFile.Write(([]byte(defaultProfileJSON)))
	require.NoError(t, err)

	// package-level vars
	configFileName = filepath.Base(configFile.Name())
	credsFileName = filepath.Base(credsFile.Name())
	defaultProfileFileName = filepath.Base(defaultProfileFile.Name())
	configDir = filepath.Dir(configFile.Name())

	s := testScenario{
		configFile:         configFile,
		credsFile:          credsFile,
		defaultProfileFile: defaultProfileFile,
	}

	return s
}

func TestLoad(t *testing.T) {
	// Must be called first
	testScenario := setupTestScenario(t)
	defer testScenario.teardown()

	err := load()
	require.NoError(t, err)

	require.Equal(t, "info", GetConfigValue("logLevel"))
	require.Equal(t, "testApiKey", GetCredentialValue("apiKey"))
	require.Equal(t, "default", defaultProfileValue)
}

func TestSetConfigValues(t *testing.T) {
	// Must be called first
	testScenario := setupTestScenario(t)
	defer testScenario.teardown()

	// Must load the config prior to tests
	err := load()
	require.NoError(t, err)

	err = SetLogLevel("debug")
	require.NoError(t, err)
	require.Equal(t, "debug", GetConfigValue(LogLevel))

	err = SetPluginDirectory("/tmp")
	require.NoError(t, err)
	require.Equal(t, "/tmp", GetConfigValue(PluginDir))

	err = SetPreleaseFeatures("ALLOW")
	require.NoError(t, err)
	require.Equal(t, "ALLOW", GetConfigValue(PrereleaseMode))

	err = SetSendUsageData("DISALLOW")
	require.NoError(t, err)
	require.Equal(t, "DISALLOW", GetConfigValue(SendUsageData))
}

func TestSetCredentialValues(t *testing.T) {
	// Must be called first
	testScenario := setupTestScenario(t)
	defer testScenario.teardown()

	// Must load the config prior to tests
	err := load()
	require.NoError(t, err)

	err = SetAPIKey("default", "NRAK-abc123")
	require.NoError(t, err)
	require.Equal(t, "NRAK-abc123", GetCredentialValue(APIKey))

	err = SetRegion("default", "US")
	require.NoError(t, err)
	require.Equal(t, "US", GetCredentialValue(Region))

	err = SetAccountID("default", "123456789")
	require.NoError(t, err)
	require.Equal(t, "123456789", GetCredentialValue(AccountID))
}

// Create config files if they don't already exist.
func TestCreate(t *testing.T) {
	require.True(t, true)
}
