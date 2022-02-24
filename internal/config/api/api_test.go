//go:build integration
// +build integration

package api

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/config"
)

var (
	testConfig = `{
		"*":{
			"loglevel":"debug",
			"plugindir": ".newrelic/plugins",
			"prereleasefeatures": "NOT_ASKED",
			"sendusagedata": "NOT_ASKED",
		}
	}`

	testCredentials = `{
		"default": {
			"apiKey": "testApiKey",
			"region": "testRegion",
			"accountID": 12345,
			"licenseKey": "testLicenseKey"
		},
		"another": {
			"apiKey": "anotherTestApiKey",
			"region": "anotherTestRegion",
			"accountID": 67890,
			"licenseKey": "anotherTestLicenseKey"
		},
	}`
)

func TestGetActiveProfileValues(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	err = ioutil.WriteFile(filepath.Join(dir, config.CredentialsFileName), []byte(testCredentials), 0644)
	require.NoError(t, err)

	os.Unsetenv("NEW_RELIC_API_KEY")
	os.Unsetenv("NEW_RELIC_LICENSE_KEY")
	os.Unsetenv("NEW_RELIC_REGION")
	os.Unsetenv("NEW_RELIC_ACCOUNT_ID")

	require.Equal(t, "testApiKey", GetActiveProfileString("apiKey"))
	require.Equal(t, "testRegion", GetActiveProfileString("region"))
	require.Equal(t, "testLicenseKey", GetActiveProfileString("licenseKey"))
	require.Equal(t, 12345, GetActiveProfileAccountID())
}

func TestGetActiveProfileValues_EnvVarOverride(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	err = ioutil.WriteFile(filepath.Join(dir, config.CredentialsFileName), []byte(testCredentials), 0644)
	require.NoError(t, err)

	os.Setenv("NEW_RELIC_API_KEY", "apiKeyOverride")
	os.Setenv("NEW_RELIC_LICENSE_KEY", "licenseKeyOverride")
	os.Setenv("NEW_RELIC_REGION", "regionOverride")
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "67890")

	require.Equal(t, "apiKeyOverride", GetActiveProfileString("apiKey"))
	require.Equal(t, "regionOverride", GetActiveProfileString("region"))
	require.Equal(t, "licenseKeyOverride", GetActiveProfileString("licenseKey"))
	require.Equal(t, 67890, GetActiveProfileAccountID())
}

func TestGetConfigValues(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	err = ioutil.WriteFile(filepath.Join(dir, config.ConfigFileName), []byte(testConfig), 0644)
	require.NoError(t, err)

	os.Unsetenv("NEW_RELIC_CLI_LOG_LEVEL")
	os.Unsetenv("NEW_RELIC_CLI_PLUGIN_DIR")
	os.Unsetenv("NEW_RELIC_CLI_PRERELEASEFEATURES")
	os.Unsetenv("NEW_RELIC_CLI_SENDUSAGEDATA")

	require.Equal(t, "debug", GetConfigString("loglevel"))
	require.Equal(t, ".newrelic/plugins", GetConfigString("plugindir"))
	require.Equal(t, config.TernaryValues.Unknown, GetConfigTernary("prereleasefeatures"))
	require.Equal(t, config.TernaryValues.Unknown, GetConfigTernary("sendusagedata"))
}

func TestGetConfigValues_EnvVarOverride(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	err = ioutil.WriteFile(filepath.Join(dir, config.ConfigFileName), []byte(testConfig), 0644)
	require.NoError(t, err)

	os.Setenv("NEW_RELIC_CLI_LOG_LEVEL", "trace")
	os.Setenv("NEW_RELIC_CLI_PLUGIN_DIR", "/tmp")
	os.Setenv("NEW_RELIC_CLI_PRERELEASEFEATURES", "ALLOW")
	os.Setenv("NEW_RELIC_CLI_SENDUSAGEDATA", "ALLOW")

	require.Equal(t, "trace", GetConfigString("loglevel"))
	require.Equal(t, "/tmp", GetConfigString("plugindir"))
	require.Equal(t, config.TernaryValues.Allow, GetConfigTernary("prereleasefeatures"))
	require.Equal(t, config.TernaryValues.Allow, GetConfigTernary("sendusagedata"))
}

func TestGetConfigValues_DefaultValues(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	os.Unsetenv("NEW_RELIC_CLI_LOG_LEVEL")
	os.Unsetenv("NEW_RELIC_CLI_PLUGIN_DIR")
	os.Unsetenv("NEW_RELIC_CLI_PRERELEASEFEATURES")
	os.Unsetenv("NEW_RELIC_CLI_SENDUSAGEDATA")

	require.Equal(t, "info", GetConfigString("loglevel"))
	require.Equal(t, filepath.Join(dir, config.DefaultPluginDir), GetConfigString("plugindir"))
	require.Equal(t, config.TernaryValues.Unknown, GetConfigTernary("prereleasefeatures"))
	require.Equal(t, config.TernaryValues.Unknown, GetConfigTernary("sendusagedata"))
}

func TestGetProfileInt(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	os.Unsetenv("NEW_RELIC_ACCOUNT_ID")
	err = ioutil.WriteFile(filepath.Join(dir, config.CredentialsFileName), []byte(testCredentials), 0644)
	require.NoError(t, err)

	config.Init(dir)

	a := GetProfileInt("another", config.AccountID)
	require.Equal(t, 67890, a)
}

func TestGetProfileInt_NotFound(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	os.Unsetenv("NEW_RELIC_ACCOUNT_ID")
	err = ioutil.WriteFile(filepath.Join(dir, config.CredentialsFileName), []byte(testCredentials), 0644)
	require.NoError(t, err)

	config.Init(dir)

	a := GetProfileInt("dne", config.AccountID)
	require.Equal(t, 0, a)
}

func TestGetProfileString(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	os.Unsetenv("NEW_RELIC_API_KEY")
	err = ioutil.WriteFile(filepath.Join(dir, config.CredentialsFileName), []byte(testCredentials), 0644)
	require.NoError(t, err)

	config.Init(dir)

	a := GetProfileString("another", config.APIKey)
	require.Equal(t, "anotherTestApiKey", a)
}

func TestGetProfileString_NotFound(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	os.Unsetenv("NEW_RELIC_API_KEY")
	err = ioutil.WriteFile(filepath.Join(dir, config.CredentialsFileName), []byte(testCredentials), 0644)
	require.NoError(t, err)

	config.Init(dir)

	a := GetProfileString("dne", config.APIKey)
	require.Equal(t, "", a)
}

func TestGetConfigString_NotFound(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	a := GetConfigString(config.FieldKey("dne"))
	require.Equal(t, "", a)
}

func TestGetConfigTernary_NotFound(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	a := GetConfigTernary(config.FieldKey("dne"))
	require.Equal(t, config.Ternary(""), a)
}

func TestGetLogLevel(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	err = ioutil.WriteFile(filepath.Join(dir, config.ConfigFileName), []byte(testConfig), 0644)
	require.NoError(t, err)

	os.Unsetenv("NEW_RELIC_CLI_LOG_LEVEL")

	require.Equal(t, "debug", GetLogLevel())
}

func TestGetLogLevel_FlagOverride(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	err = ioutil.WriteFile(filepath.Join(dir, config.ConfigFileName), []byte(testConfig), 0644)
	require.NoError(t, err)

	os.Unsetenv("NEW_RELIC_CLI_LOG_LEVEL")

	err = SetConfigValue(config.LogLevel, "info")
	require.NoError(t, err)

	config.FlagDebug = true
	require.Equal(t, "debug", GetLogLevel())

	config.FlagTrace = true
	require.Equal(t, "trace", GetLogLevel())

	// clean up
	config.FlagDebug = false
	config.FlagTrace = false
}

func TestGetLogLevel_EnvVarOveriide(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	err = ioutil.WriteFile(filepath.Join(dir, config.ConfigFileName), []byte(testConfig), 0644)
	require.NoError(t, err)

	os.Setenv("NEW_RELIC_CLI_LOG_LEVEL", "trace")

	config.FlagDebug = true

	require.Equal(t, "trace", GetLogLevel())

	// clean up
	config.FlagDebug = false
}

func TestSetProfileValue(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	os.Unsetenv("NEW_RELIC_API_KEY")
	err = ioutil.WriteFile(filepath.Join(dir, config.CredentialsFileName), []byte(testCredentials), 0644)
	require.NoError(t, err)

	config.Init(dir)

	err = SetProfileValue("another", config.APIKey, "override")
	require.NoError(t, err)

	a := GetProfileString("another", config.APIKey)
	require.Equal(t, "override", a)
}
func TestRemoveProfile(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	filename := filepath.Join(dir, config.CredentialsFileName)
	err = ioutil.WriteFile(filename, []byte(testCredentials), 0644)
	require.NoError(t, err)

	err = RemoveProfile("default")
	require.NoError(t, err)

	profiles := GetProfileNames()
	require.Equal(t, 1, len(profiles))
	require.Equal(t, "another", profiles[0])
}

func TestRemoveProfile_SetRemainingAsDefault(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	filename := filepath.Join(dir, config.CredentialsFileName)
	err = ioutil.WriteFile(filename, []byte(testCredentials), 0644)
	require.NoError(t, err)

	config.Init(dir)

	err = RemoveProfile("default")
	require.NoError(t, err)

	p := GetActiveProfileName()
	require.Equal(t, "another", p)
}

func TestRemoveProfile_RemoveDefaultProfileFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	filename := filepath.Join(dir, config.CredentialsFileName)
	err = ioutil.WriteFile(filename, []byte(testCredentials), 0644)
	require.NoError(t, err)

	config.Init(dir)

	err = RemoveProfile("default")
	require.NoError(t, err)

	err = RemoveProfile("another")
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(dir, config.DefaultProfileFileName))
	require.True(t, os.IsNotExist(err))
}

func TestGetActiveProfileName_FlagOverride(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	filename := filepath.Join(dir, config.CredentialsFileName)
	err = ioutil.WriteFile(filename, []byte(testCredentials), 0644)
	require.NoError(t, err)

	config.FlagProfileName = "override"

	p := GetActiveProfileName()
	require.Equal(t, "override", p)

	// clean up
	config.FlagProfileName = ""
}

func TestGetActiveProfileName_DefaultProfile(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	filename := filepath.Join(dir, config.CredentialsFileName)
	err = ioutil.WriteFile(filename, []byte(testCredentials), 0644)
	require.NoError(t, err)

	defaultProfileFilePath := filepath.Join(dir, config.DefaultProfileFileName)
	err = ioutil.WriteFile(defaultProfileFilePath, []byte("\"another\""), 0644)
	require.NoError(t, err)

	p := GetActiveProfileName()
	require.Equal(t, "another", p)
}

func TestSetDefaultProfile(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	filename := filepath.Join(dir, config.CredentialsFileName)
	err = ioutil.WriteFile(filename, []byte(testCredentials), 0644)
	require.NoError(t, err)

	err = SetDefaultProfile("another")
	require.NoError(t, err)

	p := GetActiveProfileName()
	require.Equal(t, "another", p)
}

func TestSetDefaultProfile_ProfileDoesNotExist(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	err = SetDefaultProfile("dne")
	require.Error(t, err)
}

func TestDeleteConfigValue(t *testing.T) {
	dir, err := ioutil.TempDir("", "newrelic-cli.config_test.*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	config.Init(dir)

	err = ioutil.WriteFile(filepath.Join(dir, config.ConfigFileName), []byte(testConfig), 0644)
	require.NoError(t, err)

	os.Unsetenv("NEW_RELIC_CLI_LOG_LEVEL")

	err = DeleteConfigValue(config.LogLevel)
	require.NoError(t, err)
	require.Equal(t, config.DefaultLogLevel, GetLogLevel())
}

func TestGetConfigFieldDefinition(t *testing.T) {
	fd := GetConfigFieldDefinition(config.LogLevel)
	require.NotNil(t, fd)
	require.Equal(t, config.LogLevel, fd.Key)
	require.Equal(t, config.DefaultLogLevel, fd.Default)
	require.Equal(t, "NEW_RELIC_CLI_LOG_LEVEL", fd.EnvVar)
	require.Contains(t, getFunctionName(fd.SetValidationFunc), "StringInStrings")

	fd = GetConfigFieldDefinition(config.PluginDir)
	require.NotNil(t, fd)
	require.Equal(t, config.PluginDir, fd.Key)
	require.Equal(t, filepath.Join(config.BasePath, config.DefaultPluginDir), fd.Default)
	require.Equal(t, "NEW_RELIC_CLI_PLUGIN_DIR", fd.EnvVar)

	fd = GetConfigFieldDefinition(config.PreReleaseFeatures)
	require.NotNil(t, fd)
	require.Equal(t, config.PreReleaseFeatures, fd.Key)
	require.Equal(t, config.TernaryValues.Unknown, fd.Default)
	require.Equal(t, "NEW_RELIC_CLI_PRERELEASEFEATURES", fd.EnvVar)
	require.Contains(t, getFunctionName(fd.SetValidationFunc), "IsTernary")

	fd = GetConfigFieldDefinition(config.SendUsageData)
	require.NotNil(t, fd)
	require.Equal(t, config.SendUsageData, fd.Key)
	require.Equal(t, config.TernaryValues.Unknown, fd.Default)
	require.Equal(t, "NEW_RELIC_CLI_SENDUSAGEDATA", fd.EnvVar)
	require.Contains(t, getFunctionName(fd.SetValidationFunc), "IsTernary")
}

func TestForEachConfigFieldDefinition(t *testing.T) {
	count := 0
	fn := func(fd config.FieldDefinition) {
		count++
	}

	ForEachConfigFieldDefinition(fn)
	require.Equal(t, 4, count)
}

func TestForEachProfileFieldDefinition(t *testing.T) {
	count := 0
	fn := func(fd config.FieldDefinition) {
		count++
	}

	ForEachProfileFieldDefinition("default", fn)
	require.Equal(t, 4, count)
}

func TestGetValidConfigFieldKeys(t *testing.T) {
	k := GetValidConfigFieldKeys()
	require.Equal(t, 4, len(k))
}

func getFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
