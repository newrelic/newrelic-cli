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

func Test_shouldExpandPreInstall(t *testing.T) {
	m := make(map[string]interface{})
	mm := make(map[interface{}]interface{})
	m["preInstall"] = mm
	mm["info"] = "A"
	mm["prompt"] = "B"
	mm["requireAtDiscovery"] = "C"

	c := expandPreInstall(m)
	require.Equal(t, "A", c.Info, "Preinstall info should equal")
	require.Equal(t, "B", c.Prompt, "Preinstall prompt should equal")
	require.Equal(t, "C", c.RequireAtDiscovery, "Preinstall requiredAtDiscovery should equal")
}

func Test_shouldExpandDiscoveryMode(t *testing.T) {
	m := make(map[string]interface{})
	dm := expandDiscoveryMode(m)

	require.Equal(t, 2, len(dm), "Omit discovery mode should return both guided and targeted")
	require.Equal(t, OpenInstallationDiscoveryModeTypes.GUIDED, dm[0], "Omit discovery mode should return both guided and targeted")
	require.Equal(t, OpenInstallationDiscoveryModeTypes.TARGETED, dm[1], "Omit discovery mode should return both guided and targeted")

	m["discoveryMode"] = []interface{}{"guided"}
	dm = expandDiscoveryMode(m)
	require.Equal(t, 1, len(dm), "Omit discovery mode should return both guided and targeted")
	require.Equal(t, OpenInstallationDiscoveryModeTypes.GUIDED, dm[0], "Only guided mode")

	m["discoveryMode"] = []interface{}{"targeted"}
	dm = expandDiscoveryMode(m)
	require.Equal(t, 1, len(dm), "Omit discovery mode should return both guided and targeted")
	require.Equal(t, OpenInstallationDiscoveryModeTypes.TARGETED, dm[0], "Only target mode")

	m["discoveryMode"] = []interface{}{"badMode"}
	dm = expandDiscoveryMode(m)
	require.Equal(t, 0, len(dm), "Bad value should return nothing")

	m["discoveryMode"] = []interface{}{"badMode", "guided"}
	dm = expandDiscoveryMode(m)
	require.Equal(t, 1, len(dm), "One good value should be parsed")
}

func TestOhiKeywordInRecipeWorks(t *testing.T) {
	r := OpenInstallationRecipe{}
	r.Keywords = []string{"ohi"}

	require.Equal(t, r.IsOhi(), true)
}
