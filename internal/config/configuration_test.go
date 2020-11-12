// +build integration

package config

import (
	"fmt"
	"os"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func getTestConfigDir() string {
	home, err := homedir.Dir()
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s/.newrelic-test", home)
}

func TestConfigure(t *testing.T) {
	t.Parallel()

	testConfigDir := getTestConfigDir()

	config, err := Configure(testConfigDir)

	require.NoError(t, err)
	require.IsType(t, new(viper.Viper), config.ViperInstance)

	// Clean up test config file
	defer func() {
		// e := os.RemoveAll(testConfigDir)

		e := os.Remove(config.GetConfigFilePath())
		if e != nil {
			log.Warnf("error deleting temporary configuration directory '%s': %s", testConfigDir, err)
		}
	}()
}
