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

	err := SetLogLevel("debug")
	require.NoError(t, err)

	configValue, err := GetConfigValue(LogLevel)
	require.NoError(t, err)
	require.Equal(t, "debug", configValue)

	err = SetPluginDirectory("/tmp/dir")
	require.NoError(t, err)

	configValue, err = GetConfigValue(PluginDir)
	require.NoError(t, err)
	require.Equal(t, "/tmp/dir", configValue)

	err = SetPreleaseFeatures(TernaryValues.Allow.String())
	require.NoError(t, err)

	configValue, err = GetConfigValue(PrereleaseFeatures)
	require.NoError(t, err)
	require.Equal(t, TernaryValues.Allow.String(), configValue)

	err = SetSendUsageData(TernaryValues.Disallow.String())
	require.NoError(t, err)

	configValue, err = GetConfigValue(SendUsageData)
	require.NoError(t, err)
	require.Equal(t, TernaryValues.Disallow.String(), configValue)
}

func TestSetConfigValue_FileNotExists(t *testing.T) {
	setupBlankSlateScenario(t)

	configFilePath := path.Join(configDir, configFilename)

	_, err := os.Stat(configFilePath)
	require.True(t, os.IsNotExist(err))

	err = SetLogLevel("debug")
	require.NoError(t, err)

	configValue, err := GetConfigValue(LogLevel)
	require.NoError(t, err)
	require.Equal(t, "debug", configValue)

	_, err = os.Stat(configFilePath)
	require.NoError(t, err)

	os.Remove(configFilePath)
}

func TestGetActiveProfileValue_Basic(t *testing.T) {
	envVarResolver = &mockEnvResolver{}
	mockConfigFiles := createMockConfigFiles(t)
	defer mockConfigFiles.teardown()

	credsValue, err := GetActiveProfileValue(APIKey)
	require.NoError(t, err)
	require.Equal(t, "testApiKey", credsValue)

	credsValue, err = GetActiveProfileValue(Region)
	require.NoError(t, err)
	require.Equal(t, "US", credsValue)

	credsValue, err = GetActiveProfileValue(AccountID)
	require.NoError(t, err)
	require.Equal(t, float64(12345), credsValue)

	credsValue, err = GetActiveProfileValue(LicenseKey)
	require.NoError(t, err)
	require.Equal(t, "testLicenseKey", credsValue)
}

func TestGetActiveProfileValue_InvalidKey(t *testing.T) {
	envVarResolver = &mockEnvResolver{}
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
	m := &mockEnvResolver{}
	envVarResolver = m
	mockConfigFiles := createMockConfigFiles(t)
	defer mockConfigFiles.teardown()

	m.GetenvVal = "newAPIKey"

	credsValue, err := GetActiveProfileValue(APIKey)
	require.NoError(t, err)
	require.Equal(t, "newAPIKey", credsValue)
}

func TestSetProfileValue_Basic(t *testing.T) {
	envVarResolver = &mockEnvResolver{}
	mockConfigFiles := createMockConfigFiles(t)
	defer mockConfigFiles.teardown()

	err := SetAPIKey("default", "NRAK-abc123")
	require.NoError(t, err)

	credsValue, err := GetActiveProfileValue(APIKey)
	require.NoError(t, err)
	require.Equal(t, "NRAK-abc123", credsValue)

	err = SetRegion("default", "US")
	require.NoError(t, err)

	credsValue, err = GetActiveProfileValue(Region)
	require.NoError(t, err)
	require.Equal(t, "US", credsValue)

	err = SetAccountID("default", 123456789)
	require.NoError(t, err)

	credsValue, err = GetActiveProfileValue(AccountID)
	require.NoError(t, err)
	require.Equal(t, float64(123456789), credsValue)

	err = SetLicenseKey("default", "license")
	require.NoError(t, err)

	credsValue, err = GetActiveProfileValue(LicenseKey)
	require.NoError(t, err)
	require.Equal(t, "license", credsValue)
}

func TestSetProfileValue_FileNotExists(t *testing.T) {
	envVarResolver = &mockEnvResolver{}
	setupBlankSlateScenario(t)

	credsFilePath := path.Join(configDir, credsFilename)
	defaultProfileFilePath := path.Join(configDir, defaultProfileFilename)

	_, err := os.Stat(credsFilePath)
	require.True(t, os.IsNotExist(err))

	credsValue, err := GetActiveProfileValue(APIKey)
	require.NoError(t, err)
	require.Nil(t, credsValue)

	err = SetAPIKey("default", "NRAK-abc123")
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
	envVarResolver = &mockEnvResolver{}
	setupBlankSlateScenario(t)

	defaultProfileFilePath := path.Join(configDir, defaultProfileFilename)

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
	configDir = filepath.Dir(configFile.Name())
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
	configDir = os.TempDir()
}

const numBytes = "0123456789"

func randNumBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = numBytes[rand.Intn(len(numBytes))]
	}
	return string(b)
}
