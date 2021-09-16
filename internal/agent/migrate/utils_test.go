// +build unit

package migrate

import (
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPopulateV4ConfigRelativePath(t *testing.T) {
	command := []string{
		"./relative/path/testBinary",
		"--metrics",
		"-events",
		"--not_standard_flag",
		"not_standard_value",
	}
	labelMap := map[string]string{
		"oneKeyLabel": "oneValueLabel",
		"twoKeyLabel": "twoValueLabel",
	}
	argMap := map[string]string{
		"oneKeyArg": "oneValueArg",
		"twoKeyArg": "twoValueArg",
	}
	plugin := getPluginExample(command)
	pluginWrapper := getPluginWrapperExample(argMap, labelMap)

	v4Config, err := populateV4Config(plugin, pluginWrapper)
	require.NoError(t, err, "no error is expected")
	require.Len(t, v4Config.Integrations, 1)

	instance := v4Config.Integrations[0]
	assert.Equal(t, "testUser", instance.User)
	assert.Equal(t, "inventoryPrefix", instance.InventorySource)
	assert.Equal(t, "12345s", instance.Interval)
	assert.Equal(t, "true", instance.Env["METRICS"])
	assert.Equal(t, "true", instance.Env["EVENTS"])
	assert.Equal(t, instance.Labels, labelMap)
	for k, v := range argMap {
		assert.Equal(t, instance.Env[strings.ToUpper(k)], v)
	}
	assert.Equal(t, ShlexOpt(nil), instance.Exec)
	assert.Equal(t, []string{"--not_standard_flag", "not_standard_value"}, instance.CLIArgs)
	assert.Equal(t, "testBinary", instance.InstanceName)
}

func TestPopulateV4ConfigAbsolutePath(t *testing.T) {
	command := []string{
		"/absolute/path/testBinary",
		"--metrics",
		"-events",
		"--not_standard_flag",
		"not_standard_value",
	}

	plugin := getPluginExample(command)
	pluginWrapper := getPluginWrapperExample(nil, nil)
	v4Config, err := populateV4Config(plugin, pluginWrapper)

	require.NoError(t, err, "no error is expected")
	require.Len(t, v4Config.Integrations, 1)

	instance := v4Config.Integrations[0]
	assert.Equal(t, "testUser", instance.User)
	assert.Equal(t, "inventoryPrefix", instance.InventorySource)
	assert.Equal(t, "12345s", instance.Interval)
	assert.NotEqual(t, "true", instance.Env["METRICS"])
	assert.NotEqual(t, "true", instance.Env["EVENTS"])
	assert.Equal(t, []string(nil), instance.CLIArgs)
	assert.Equal(t, ShlexOpt{
		"/absolute/path/testBinary",
		"--metrics",
		"-events",
		"--not_standard_flag",
		"not_standard_value",
	}, instance.Exec)
	assert.Equal(t, "testBinary", instance.InstanceName)

}

func TestYAMLValidity(t *testing.T) {
	tmpPath := t.TempDir()
	tmpFile := path.Join(tmpPath, "testingFile")
	tmpCommentFile := path.Join(tmpPath, "commentFile")
	err := ioutil.WriteFile(tmpCommentFile, []byte("testComment"), 0666)
	require.NoError(t, err, "no error is expected")

	command := []string{
		"./relative/path/testBinary",
		"--metrics",
		"-events",
		"--not_standard_flag",
		"not_standard_value",
	}
	argMap := map[string]string{
		"oneKeyArg": "oneValueArg",
		"twoKeyArg": "twoValueArg",
	}

	plugin := getPluginExample(command)
	pluginWrapper := getPluginWrapperExample(argMap, nil)
	v4Config, err := populateV4Config(plugin, pluginWrapper)
	require.NoError(t, err, "no error is expected")

	err = writeOutput(v4Config, tmpCommentFile, tmpCommentFile, tmpFile)
	require.NoError(t, err, "no error is expected")

	testV4 := v4{}
	err = readAndUnmarshallConfig(tmpFile, &testV4)
	require.NoError(t, err, "no error is expected")
}

func TestPopulateV4ConfigWrongName(t *testing.T) {
	v4Config, err := populateV4Config(
		Plugin{
			Name: "testIntegration",
		},
		PluginInstanceWrapper{
			IntegrationName: "different",
		},
	)

	require.Error(t, err, "error is expected")
	require.Nil(t, v4Config, 0)
}

func getPluginExample(command []string) Plugin {
	plugin := Plugin{
		Name: "testIntegration",
		Commands: map[string]*PluginV1Command{
			"testCommand": {
				Command:  command,
				Prefix:   "inventoryPrefix",
				Interval: 12345,
			},
		},
	}
	return plugin
}

func getPluginWrapperExample(argMap map[string]string, labelMap map[string]string) PluginInstanceWrapper {
	pluginWrapper := PluginInstanceWrapper{
		IntegrationName: "testIntegration",
		Instances: []*PluginV1Instance{
			{
				Name:            "testCommandName",
				Command:         "testCommand",
				Arguments:       argMap,
				Labels:          labelMap,
				IntegrationUser: "testUser",
			},
			{
				Name:            "testCommandNameTwo",
				Command:         "testCommandTwo",
				Arguments:       map[string]string{},
				Labels:          map[string]string{},
				IntegrationUser: "testUserTwo",
			},
		},
	}
	return pluginWrapper
}
