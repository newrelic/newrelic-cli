// +build unit

package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToStringByFieldName(t *testing.T) {
	data := map[string]interface{}{
		"intField":    9600, // e.g. the port used for a db connection URL
		"stringField": "stringValue",
		"boolField":   false,
	}

	intAsString := toStringByFieldName("intField", data)
	require.Equal(t, "9600", intAsString)

	stringAsString := toStringByFieldName("stringField", data)
	require.Equal(t, "stringValue", stringAsString)

	boolAsString := toStringByFieldName("boolField", data)
	require.Equal(t, "false", boolAsString)
}

func TestExpandValidation(t *testing.T) {
	data := map[string]interface{}{
		"validation": map[interface{}]interface{}{
			"agentUrl": "http://localhost:18003/v1/status/entity",
			"nrql":     "SELECT count(*) from SystemSample",
		},
	}

	output := expandValidation(data)
	require.Equal(t, "http://localhost:18003/v1/status/entity", output.AgentURL)
	require.Equal(t, "SELECT count(*) from SystemSample", output.NRQL)
}

func TestExpandValidationEmpty(t *testing.T) {
	expected := expandValidation(map[string]interface{}{})

	require.Equal(t, OpenInstallationRecipeValidationConfig{}, expected)
}

func Test_shouldDisplayShortString(t *testing.T) {

	recipe := OpenInstallationRecipe{
		Name:           "test-recipe",
		DisplayName:    "my verbose looking recipe name",
		ValidationNRQL: "testNrql",
		InstallTargets: []OpenInstallationRecipeInstallTarget{
			{
				Type:            OpenInstallationTargetTypeTypes.HOST,
				Os:              OpenInstallationOperatingSystemTypes.DARWIN,
				Platform:        OpenInstallationPlatformTypes.AMAZON,
				PlatformVersion: "2",
				KernelArch:      "x86",
			},
		},
	}

	output := recipe.ToShortDisplayString()
	require.Equal(t, "test-recipe (darwin/amazon/2/x86)", output)
}

func Test_shouldDisplayShortStringMultipleTargets(t *testing.T) {

	recipe := OpenInstallationRecipe{
		Name:           "test-recipe",
		DisplayName:    "my verbose looking recipe name",
		ValidationNRQL: "testNrql",
		InstallTargets: []OpenInstallationRecipeInstallTarget{
			{
				Type:            OpenInstallationTargetTypeTypes.HOST,
				Os:              OpenInstallationOperatingSystemTypes.DARWIN,
				Platform:        OpenInstallationPlatformTypes.AMAZON,
				PlatformVersion: "2",
				KernelArch:      "x86",
			},
			{
				Type:            OpenInstallationTargetTypeTypes.HOST,
				Os:              OpenInstallationOperatingSystemTypes.LINUX,
				Platform:        OpenInstallationPlatformTypes.REDHAT,
				PlatformVersion: "8",
				KernelArch:      "arm",
			},
		},
	}

	output := recipe.ToShortDisplayString()
	require.True(t, strings.Contains(output, "test-recipe"))
	require.True(t, strings.Contains(output, "darwin/amazon/2/x86"))
	require.True(t, strings.Contains(output, "linux/redhat/8/arm"))
}
