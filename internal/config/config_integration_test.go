// +build integration

package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigSetLogLevel(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	// Initialize the new configuration directory
	c, err := LoadConfig(f)
	assert.NoError(t, err)
	assert.Equal(t, c.configDir, f)

	// Set the valid log levels
	for _, l := range []string{
		"ERROR",
		"WARN",
		"INFO",
		"DEBUG",
		"TRACE",
	} {
		err = c.Set("logLevel", l)
		assert.NoError(t, err)
		assert.Equal(t, l, c.LogLevel)

		// Double check that the config is written to disk
		c2, err := LoadConfig(f)
		assert.NoError(t, err)
		assert.Equal(t, l, c2.LogLevel)
	}

	err = c.Set("logLevel", "INVALID_VALUE")
	assert.Error(t, err)

	err = c.Set("loglevel", "Info")
	assert.Error(t, err)

	err = c.Set("Loglevel", "Debug")
	assert.Error(t, err)

}

func TestConfigSetSendUsageData(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	// Initialize the new configuration directory
	c, err := LoadConfig(f)
	assert.NoError(t, err)
	assert.Equal(t, c.configDir, f)

	// Set the valid sendUsageData values
	for _, l := range []Ternary{
		TernaryValues.Allow,
		TernaryValues.Disallow,
		TernaryValues.Unknown,
	} {
		err = c.Set("sendUsageData", l)
		assert.NoError(t, err)
		assert.Equal(t, l, c.SendUsageData)
	}

	err = c.Set("sendUsageData", "INVALID_VALUE")
	assert.Error(t, err)
}

func TestConfigSetPreReleaseFeatures(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	// Initialize the new configuration directory
	c, err := LoadConfig(f)
	assert.NoError(t, err)
	assert.Equal(t, c.configDir, f)

	// Set the valid pre-release feature values
	for _, l := range []Ternary{
		TernaryValues.Allow,
		TernaryValues.Disallow,
		TernaryValues.Unknown,
	} {
		err = c.Set("preReleaseFeatures", l)
		assert.NoError(t, err)
		assert.Equal(t, l, c.PreReleaseFeatures)
	}

	err = c.Set("preReleaseFeatures", "INVALID_VALUE")
	assert.Error(t, err)
}

func TestConfigSetPluginDir(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	// Initialize the new configuration directory
	c, err := LoadConfig(f)
	assert.NoError(t, err)
	assert.Equal(t, c.configDir, f)

	err = c.Set("pluginDir", "test")
	assert.NoError(t, err)
	assert.Equal(t, "test", c.PluginDir)
}
