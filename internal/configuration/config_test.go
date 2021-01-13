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

func TestSetConfigValue_Basic(t *testing.T) {
	mockConfigFiles := createMockConfigFiles(t)
	defer mockConfigFiles.teardown()

	err := SetLogLevel("debug")
	require.NoError(t, err)

	configValue, err := GetConfigValue(LogLevel)
	require.NoError(t, err)
	require.Equal(t, "debug", configValue)

	err = SetPluginDirectory("/tmp")
	require.NoError(t, err)

	configValue, err = GetConfigValue(PluginDir)
	require.NoError(t, err)
	require.Equal(t, "/tmp", configValue)

	err = SetPreleaseFeatures("ALLOW")
	require.NoError(t, err)

	configValue, err = GetConfigValue(PrereleaseMode)
	require.NoError(t, err)
	require.Equal(t, "ALLOW", configValue)

	err = SetSendUsageData("DISALLOW")
	require.NoError(t, err)

	configValue, err = GetConfigValue(SendUsageData)
	require.NoError(t, err)
	require.Equal(t, "DISALLOW", configValue)
}

func TestSetCredentialValue_Basic(t *testing.T) {
	mockConfigFiles := createMockConfigFiles(t)
	defer mockConfigFiles.teardown()

	err := SetAPIKey("default", "NRAK-abc123")
	require.NoError(t, err)
	require.Equal(t, "NRAK-abc123", GetCredentialValue(APIKey))

	err = SetRegion("default", "US")
	require.NoError(t, err)
	require.Equal(t, "US", GetCredentialValue(Region))

	err = SetAccountID("default", "123456789")
	require.NoError(t, err)
	require.Equal(t, "123456789", GetCredentialValue(AccountID))
}

func TestSetConfigValue_FileNotExists(t *testing.T) {
	setupBlankSlateScenario(t)

	configFilePath := path.Join(configDir, configFilename)
	credsFilePath := path.Join(configDir, credsFilename)
	defaultProfileFilePath := path.Join(configDir, defaultProfileFilename)

	_, err := os.Stat(configFilePath)
	require.True(t, os.IsNotExist(err))

	_, err = os.Stat(credsFilePath)
	require.True(t, os.IsNotExist(err))

	_, err = os.Stat(defaultProfileFilePath)
	require.True(t, os.IsNotExist(err))

	require.Nil(t, GetCredentialValue(APIKey))

	configValue, err := GetConfigValue(LogLevel)
	require.NoError(t, err)
	require.Nil(t, configValue)

	require.Empty(t, GetDefaultProfileName())

	err = SetLogLevel("debug")
	require.NoError(t, err)

	configValue, err = GetConfigValue(LogLevel)
	require.NoError(t, err)
	require.Equal(t, "debug", configValue)

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

func TestGetConfigValue_InvalidKey(t *testing.T) {
	mockConfigFiles := createMockConfigFiles(t)
	defer mockConfigFiles.teardown()

	_, err := GetConfigValue("LOGLEVEL")
	require.NoError(t, err)

	_, err = GetConfigValue("logLevel")
	require.NoError(t, err)

	_, err = GetConfigValue("logLevel")
	require.NoError(t, err)

	_, err = GetConfigValue("invalidKey")
	require.Error(t, err)
}

type mockConfigFiles struct {
	configFile         *os.File
	credsFile          *os.File
	defaultProfileFile *os.File
}

func (s *mockConfigFiles) teardown() {
	os.Remove(s.configFile.Name())
	os.Remove(s.credsFile.Name())
	os.Remove(s.defaultProfileFile.Name())
}

func createMockConfigFiles(t *testing.T) mockConfigFiles {
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
	configDir = filepath.Dir(configFile.Name())
	configFilename = filepath.Base(configFile.Name())
	credsFilename = filepath.Base(credsFile.Name())
	defaultProfileFilename = filepath.Base(defaultProfileFile.Name())

	s := mockConfigFiles{
		configFile:         configFile,
		credsFile:          credsFile,
		defaultProfileFile: defaultProfileFile,
	}

	err = load()
	require.NoError(t, err)

	return s
}

func setupBlankSlateScenario(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	configFilename = fmt.Sprintf("config%s.json", randNumBytes(8))
	credsFilename = fmt.Sprintf("creds%s.json", randNumBytes(8))
	defaultProfileFilename = fmt.Sprintf("default-profile%s.json", randNumBytes(8))
	defaultProfileValue = ""
	configDir = os.TempDir()

	err := load()
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
