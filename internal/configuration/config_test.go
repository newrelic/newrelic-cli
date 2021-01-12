// +build integration

package configuration

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Must be called first
	testScenario := setupTestScenario(t)
	defer testScenario.teardown()
	testScenario.writeFiles(t)

	err := load()
	require.NoError(t, err)

	require.Equal(t, "info", GetConfigValue(LogLevel))
	require.Equal(t, "testApiKey", GetCredentialValue(APIKey))
	require.Equal(t, "default", defaultProfileValue)
}

func TestSetConfigValues(t *testing.T) {
	// Must be called first
	testScenario := setupTestScenario(t)
	defer testScenario.teardown()
	testScenario.writeFiles(t)

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
	testScenario.writeFiles(t)

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
	rand.Seed(time.Now().UnixNano())
	configFilename = fmt.Sprintf("config%s.json", randNumBytes(8))
	credsFilename = fmt.Sprintf("creds%s.json", randNumBytes(8))
	defaultProfileFilename = fmt.Sprintf("default-profile%s.json", randNumBytes(8))
	configDir = os.TempDir()

	configFilePath := path.Join(configDir, configFilename)
	credsFilePath := path.Join(configDir, credsFilename)
	defaultProfileFilePath := path.Join(configDir, defaultProfileFilename)

	// Must load the config prior to tests
	err := load()
	require.NoError(t, err)
	_, err = os.Stat(configFilePath)
	require.True(t, os.IsNotExist(err))

	_, err = os.Stat(credsFilePath)
	require.True(t, os.IsNotExist(err))

	_, err = os.Stat(defaultProfileFilePath)
	require.True(t, os.IsNotExist(err))

	require.Nil(t, GetCredentialValue(APIKey))
	require.Nil(t, GetConfigValue(LogLevel))
	require.Empty(t, GetDefaultProfileName())

	err = SetLogLevel("debug")
	require.NoError(t, err)
	require.Equal(t, "debug", GetConfigValue(LogLevel))

	_, err = os.Stat(configFilePath)
	require.NoError(t, err)

	err = SetAPIKey("default", "NRAK-abc123")
	require.NoError(t, err)
	require.Equal(t, "NRAK-abc123", GetCredentialValue(APIKey))

	_, err = os.Stat(credsFilePath)
	require.NoError(t, err)

	err = SetDefaultProfileName("default")
	require.NoError(t, err)
	require.Equal(t, "default", GetDefaultProfileName())

	_, err = os.Stat(defaultProfileFilePath)
	require.NoError(t, err)

	os.Remove(configFilePath)
	os.Remove(credsFilePath)
	os.Remove(defaultProfileFilePath)
}

type testScenario struct {
	configFile         *os.File
	configJSON         string
	credsFile          *os.File
	credsJSON          string
	defaultProfileFile *os.File
	defaultProfileJSON string
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
	defaultProfileFile, err := ioutil.TempFile("", "default-profile*.json")
	require.NoError(t, err)

	defaultProfileJSON := `"default"`

	// package-level vars
	configDir = filepath.Dir(configFile.Name())
	configFilename = filepath.Base(configFile.Name())
	credsFilename = filepath.Base(credsFile.Name())
	defaultProfileFilename = filepath.Base(defaultProfileFile.Name())

	s := testScenario{
		configFile:         configFile,
		configJSON:         configJSON,
		credsFile:          credsFile,
		credsJSON:          credsJSON,
		defaultProfileFile: defaultProfileFile,
		defaultProfileJSON: defaultProfileJSON,
	}

	return s
}

func (s testScenario) writeFiles(t *testing.T) {
	_, err := s.configFile.Write([]byte(s.configJSON))
	require.NoError(t, err)

	_, err = s.credsFile.Write(([]byte(s.credsJSON)))
	require.NoError(t, err)

	_, err = s.defaultProfileFile.Write(([]byte(s.defaultProfileJSON)))
	require.NoError(t, err)
}

const numBytes = "0123456789"

func randNumBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = numBytes[rand.Intn(len(numBytes))]
	}
	return string(b)
}
