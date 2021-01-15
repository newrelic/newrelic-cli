// +build integration

package config

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

func TestGetConfigValue_Basic(t *testing.T) {
	mockConfigFiles := createMockConfigFiles(t)
	defer mockConfigFiles.teardown()

	configValue, err := GetConfigValue(LogLevel)
	require.NoError(t, err)
	require.Equal(t, "info", configValue)

	configValue, err = GetConfigValue(PluginDir)
	require.NoError(t, err)
	require.Equal(t, "/tmp", configValue)

	configValue, err = GetConfigValue(PrereleaseFeatures)
	require.NoError(t, err)
	require.Equal(t, TernaryValues.Unknown.String(), configValue)

	configValue, err = GetConfigValue(SendUsageData)
	require.NoError(t, err)
	require.Equal(t, TernaryValues.Unknown.String(), configValue)
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

func TestGetConfigValue_DefaultValues(t *testing.T) {
	setupBlankSlateScenario(t)

	configValue, err := GetConfigValue(LogLevel)
	require.NoError(t, err)
	require.Equal(t, "info", configValue)

	configValue, err = GetConfigValue(PluginDir)
	require.NoError(t, err)
	require.Equal(t, "", configValue)

	configValue, err = GetConfigValue(PrereleaseFeatures)
	require.NoError(t, err)
	require.Equal(t, TernaryValues.Unknown.String(), configValue)

	configValue, err = GetConfigValue(SendUsageData)
	require.NoError(t, err)
	require.Equal(t, TernaryValues.Unknown.String(), configValue)
}

func TestSetConfigValue_Basic(t *testing.T) {
	mockConfigFiles := createMockConfigFiles(t)
	defer mockConfigFiles.teardown()

	err := SetConfigValue(LogLevel, "debug")
	require.NoError(t, err)

	configValue, err := GetConfigValue(LogLevel)
	require.NoError(t, err)
	require.Equal(t, "debug", configValue)
}

func TestSetConfigValue_FileNotExists(t *testing.T) {
	setupBlankSlateScenario(t)

	configFilePath := path.Join(ConfigDir, configFilename)

	_, err := os.Stat(configFilePath)
	require.True(t, os.IsNotExist(err))

	err = SetConfigValue(LogLevel, "debug")
	require.NoError(t, err)

	configValue, err := GetConfigValue(LogLevel)
	require.NoError(t, err)
	require.Equal(t, "debug", configValue)

	_, err = os.Stat(configFilePath)
	require.NoError(t, err)

	os.Remove(configFilePath)
}

func TestSetConfigValue_InvalidValue(t *testing.T) {
	mockConfigFiles := createMockConfigFiles(t)
	defer mockConfigFiles.teardown()

	var err error

	err = SetConfigValue(LogLevel, "invalid")
	require.Error(t, err)

	err = SetConfigValue(PrereleaseFeatures, "invalid")
	require.Error(t, err)

	err = SetConfigValue(SendUsageData, "invalid")
	require.Error(t, err)

	err = SetConfigValue(PluginDir, "/any/path/is/valid")
	require.NoError(t, err)
}

func TestGetActiveProfileValue_Basic(t *testing.T) {
	EnvVarResolver = &MockEnvResolver{}
	mockConfigFiles := createMockConfigFiles(t)
	defer mockConfigFiles.teardown()

	credsValue, err := GetActiveProfileValue(APIKey)
	require.NoError(t, err)
	require.Equal(t, "testApiKey", credsValue)
}

func TestGetActiveProfileValue_InvalidKey(t *testing.T) {
	EnvVarResolver = &MockEnvResolver{}
	mockConfigFiles := createMockConfigFiles(t)
	defer mockConfigFiles.teardown()

	_, err := GetActiveProfileValue("APIKEY")
	require.NoError(t, err)

	_, err = GetActiveProfileValue("apiKey")
	require.NoError(t, err)

	_, err = GetActiveProfileValue("apikey")
	require.NoError(t, err)

	_, err = GetActiveProfileValue("invalidKey")
	require.Error(t, err)
}

func TestGetActiveProfileValue_EnvVarOverride(t *testing.T) {
	m := &MockEnvResolver{}
	EnvVarResolver = m
	mockConfigFiles := createMockConfigFiles(t)
	defer mockConfigFiles.teardown()

	m.GetenvVal = "newAPIKey"

	credsValue, err := GetActiveProfileValue(APIKey)
	require.NoError(t, err)
	require.Equal(t, "newAPIKey", credsValue)
}

func TestSetProfileValue_Basic(t *testing.T) {
	EnvVarResolver = &MockEnvResolver{}
	mockConfigFiles := createMockConfigFiles(t)
	defer mockConfigFiles.teardown()

	err := SetProfileValue(defaultDefaultProfileName, APIKey, "NRAK-abc123")
	require.NoError(t, err)

	credsValue, err := GetProfileValue(defaultDefaultProfileName, APIKey)
	require.NoError(t, err)
	require.Equal(t, "NRAK-abc123", credsValue)
}

func TestSetProfileValue_FileNotExists(t *testing.T) {
	EnvVarResolver = &MockEnvResolver{}
	setupBlankSlateScenario(t)

	credsFilePath := path.Join(ConfigDir, credsFilename)
	defaultProfileFilePath := path.Join(ConfigDir, defaultProfileFilename)

	_, err := os.Stat(credsFilePath)
	require.True(t, os.IsNotExist(err))

	credsValue, err := GetActiveProfileValue(APIKey)
	require.NoError(t, err)
	require.Nil(t, credsValue)

	err = SetProfileValue("default", APIKey, "NRAK-abc123")
	require.NoError(t, err)

	credsValue, err = GetActiveProfileValue(APIKey)
	require.NoError(t, err)
	require.Equal(t, "NRAK-abc123", credsValue)

	_, err = os.Stat(credsFilePath)
	require.NoError(t, err)

	err = os.Remove(credsFilePath)
	require.NoError(t, err)

	err = os.Remove(defaultProfileFilePath)
	require.NoError(t, err)
}

func TestGetDefaultProfileName_Basic(t *testing.T) {
	mockConfigFiles := createMockConfigFiles(t)
	defer mockConfigFiles.teardown()

	require.Equal(t, "default", GetDefaultProfileName())
}

func TestSetDefaultProfileName_FileNotExists(t *testing.T) {
	EnvVarResolver = &MockEnvResolver{}
	setupBlankSlateScenario(t)

	defaultProfileFilePath := path.Join(ConfigDir, defaultProfileFilename)

	_, err := os.Stat(defaultProfileFilePath)
	require.True(t, os.IsNotExist(err))

	require.Empty(t, GetDefaultProfileName())

	err = SetDefaultProfileName("default")
	require.NoError(t, err)
	require.Equal(t, "default", GetDefaultProfileName())

	_, err = os.Stat(defaultProfileFilePath)
	require.NoError(t, err)

	os.Remove(defaultProfileFilePath)
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
	ConfigDir = filepath.Dir(configFile.Name())
	configFilename = filepath.Base(configFile.Name())
	credsFilename = filepath.Base(credsFile.Name())
	defaultProfileFilename = filepath.Base(defaultProfileFile.Name())

	s := mockConfigFiles{
		configFile:         configFile,
		credsFile:          credsFile,
		defaultProfileFile: defaultProfileFile,
	}

	return s
}

func setupBlankSlateScenario(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	configFilename = fmt.Sprintf("config%s.json", randNumBytes(8))
	credsFilename = fmt.Sprintf("creds%s.json", randNumBytes(8))
	defaultProfileFilename = fmt.Sprintf("default-profile%s.json", randNumBytes(8))
	ConfigDir = os.TempDir()
}

const numBytes = "0123456789"

func randNumBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = numBytes[rand.Intn(len(numBytes))]
	}
	return string(b)
}
