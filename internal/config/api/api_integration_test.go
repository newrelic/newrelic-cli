// +build integration

package api

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/config"
)

func TestConfigSetLogLevel(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	config.Init(f)

	// Set the valid log levels
	for _, l := range []string{
		"ERROR",
		"WARN",
		"INFO",
		"DEBUG",
		"TRACE",
	} {
		err = SetConfigString("logLevel", l)
		assert.NoError(t, err)

		actual := GetConfigString(config.LogLevel)
		assert.Equal(t, l, actual)
	}

	err = SetConfigString("logLevel", "INVALID_VALUE")
	assert.Error(t, err)
}

func TestConfigSetSendUsageData(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	config.Init(f)

	// Set the valid sendUsageData values
	for _, l := range []config.Ternary{
		config.TernaryValues.Allow,
		config.TernaryValues.Disallow,
		config.TernaryValues.Unknown,
	} {
		err = SetConfigValue("sendUsageData", l)
		assert.NoError(t, err)
		assert.Equal(t, GetConfigTernary(config.SendUsageData), l)
	}

	err = SetConfigValue("sendUsageData", "INVALID_VALUE")
	assert.Error(t, err)
}

func TestConfigSetPreReleaseFeatures(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	config.Init(f)

	// Set the valid pre-release feature values
	for _, l := range []config.Ternary{
		config.TernaryValues.Allow,
		config.TernaryValues.Disallow,
		config.TernaryValues.Unknown,
	} {
		err = SetConfigValue("preReleaseFeatures", l)
		assert.NoError(t, err)
		assert.Equal(t, GetConfigTernary(config.PreReleaseFeatures), l)
	}

	err = SetConfigValue("preReleaseFeatures", "INVALID_VALUE")
	assert.Error(t, err)
}

func TestConfigSetPluginDir(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	config.Init(f)

	err = SetConfigString(config.PluginDir, "test")
	assert.NoError(t, err)
	assert.Equal(t, "test", GetConfigString(config.PluginDir))
}
