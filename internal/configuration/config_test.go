// +build integration

package configuration

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	configFile, err := ioutil.TempFile("", "config.json")
	require.NoError(t, err)
	defer os.Remove(configFile.Name())

	configJson := `
{
	"*": {
		"loglevel": "info",
		"plugindir": "/tmp",
		"prereleasefeatures": "NOT_ASKED",
		"sendusagedata": "NOT_ASKED"
	}
}
`
	_, err = configFile.Write([]byte(configJson))
	require.NoError(t, err)

	credsFile, err := ioutil.TempFile("", "credentials.json")
	require.NoError(t, err)
	defer os.Remove(credsFile.Name())

	credsJson := `
	{
		"default": {
			"apiKey": "testApiKey",
			"region": "US",
			"accountID": 12345,
			"licenseKey": "testLicenseKey"
		}
	}
	`
	_, err = credsFile.Write(([]byte(credsJson)))
	require.NoError(t, err)

	defaultProfileFile, err := ioutil.TempFile("", "defaultProfile.json")
	require.NoError(t, err)
	defer os.Remove(defaultProfileFile.Name())

	defaultProfileJson := `"default"`
	_, err = defaultProfileFile.Write(([]byte(defaultProfileJson)))
	require.NoError(t, err)

	// package-level vars
	configFileName = filepath.Base(configFile.Name())
	credsFileName = filepath.Base(credsFile.Name())
	configDir = filepath.Dir(configFile.Name())

	err = load()
	require.NoError(t, err)

	require.Equal(t, "info", GetConfigValue("logLevel"))
	require.Equal(t, "testApiKey", GetProfileValue("apiKey"))
	//require.Equal(t, "testApiKey", c.Profiles["default"].APIKey)
}

// Create config files if they don't already exist.
func TestCreate(t *testing.T) {
	require.True(t, true)
}
