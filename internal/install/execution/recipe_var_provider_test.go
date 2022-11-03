package execution

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/yaml.v2"

	"github.com/go-task/task/v3/taskfile"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestRecipeVarProvider_Basic(t *testing.T) {
	e := NewRecipeVarProvider()

	tmpFile, err := ioutil.TempFile(os.TempDir(), t.Name())
	if err != nil {
		t.Fatal("error creating temp file")
	}

	defer os.Remove(tmpFile.Name())

	output := `
  {
    \"hostname\": \"{{.HOSTNAME}}\",
    \"os\": \"{{.OS}}\",
    \"platform\": \"{{.PLATFORM}}\",
    \"platformFamily\": \"{{.PLATFORM_FAMILY}}\",
    \"platformVersion\": \"{{.PLATFORM_VERSION}}\",
    \"kernelArch\": \"{{.KERNEL_ARCH}}\",
    \"kernelVersion\": \"{{.KERNEL_VERSION}}\"
  }`

	// We convert the `install` section of the recipe to a YAML string,
	// which is then used to create a Taskfile for go-task.
	recipeInstallToYaml := map[string]interface{}{
		"version": "3",
		"tasks": taskfile.Tasks{
			"default": &taskfile.Task{
				Cmds: []*taskfile.Cmd{
					{
						Cmd: fmt.Sprintf("echo %s > %s", strings.ReplaceAll(output, "\n", ""), tmpFile.Name()),
					},
				},
				Silent: true,
			},
		},
	}

	installYamlBytes, err := yaml.Marshal(recipeInstallToYaml)
	require.NoError(t, err)

	m := types.DiscoveryManifest{
		Hostname:        "testHostname",
		OS:              "testOS",
		Platform:        "testPlatform",
		PlatformFamily:  "testPlatformFamily",
		PlatformVersion: "testPlatformVersion",
		KernelArch:      "testKernelArch",
		KernelVersion:   "testKernelVersion",
	}

	r := types.OpenInstallationRecipe{
		Install: string(installYamlBytes),
	}

	v, err := e.Prepare(m, r, false, "testLicenseKey")
	require.NoError(t, err)
	require.Contains(t, m.OS, v["OS"])
	require.Contains(t, m.Platform, v["Platform"])
	require.Contains(t, m.PlatformVersion, v["PlatformVersion"])
	require.Contains(t, m.PlatformFamily, v["PlatformFamily"])
	require.Contains(t, m.KernelArch, v["KernelArch"])
	require.Contains(t, m.KernelVersion, v["KernelVersion"])
	require.Equal(t, "", v["NRIA_CUSTOM_ATTRIBUTES"])
	require.Equal(t, "", v["NRIA_PASSTHROUGH_ENVIRONMENT"])
	require.Contains(t, "https://download.newrelic.com/", v["NEW_RELIC_DOWNLOAD_URL"])

	os.Setenv(EnvInstallCustomAttributes, "test:123,bad")
	v, err = e.Prepare(m, r, false, "testLicenseKey")
	require.NoError(t, err)
	require.Equal(t, "custom_attributes:\n  test: \"123\"\n", v["NRIA_CUSTOM_ATTRIBUTES"])

	os.Setenv(EnvNriaCustomAttributes, "{\"owning_team\":\"virtuoso\"}")
	v, err = e.Prepare(m, r, false, "testLicenseKey")
	require.NoError(t, err)
	require.Equal(t, "custom_attributes:\n  owning_team: virtuoso\n  test: \"123\"\n", v["NRIA_CUSTOM_ATTRIBUTES"])
}

func TestRecipeVarProvider_CommandLineEnvarsDirectlyPassedToRecipeContext(t *testing.T) {
	e := NewRecipeVarProvider()

	tmpFile, err := ioutil.TempFile(os.TempDir(), t.Name())
	if err != nil {
		t.Fatal("error creating temp file")
	}

	defer os.Remove(tmpFile.Name())

	output := `
  {
    \"hostname\": \"{{.HOSTNAME}}\",
    \"os\": \"{{.OS}}\",
    \"platform\": \"{{.PLATFORM}}\",
    \"platformFamily\": \"{{.PLATFORM_FAMILY}}\",
    \"platformVersion\": \"{{.PLATFORM_VERSION}}\",
    \"kernelArch\": \"{{.KERNEL_ARCH}}\",
    \"kernelVersion\": \"{{.KERNEL_VERSION}}\"
  }`

	// We convert the `install` section of the recipe to a YAML string,
	// which is then used to create a Taskfile for go-task.
	recipeInstallToYaml := map[string]interface{}{
		"version": "3",
		"tasks": taskfile.Tasks{
			"default": &taskfile.Task{
				Cmds: []*taskfile.Cmd{
					{
						Cmd: fmt.Sprintf("echo %s > %s", strings.ReplaceAll(output, "\n", ""), tmpFile.Name()),
					},
				},
				Silent: true,
			},
		},
	}

	installYamlBytes, err := yaml.Marshal(recipeInstallToYaml)
	require.NoError(t, err)

	m := types.DiscoveryManifest{
		Hostname:        "testHostname",
		OS:              "testOS",
		Platform:        "testPlatform",
		PlatformFamily:  "testPlatformFamily",
		PlatformVersion: "testPlatformVersion",
		KernelArch:      "testKernelArch",
		KernelVersion:   "testKernelVersion",
	}

	r := types.OpenInstallationRecipe{
		Install: string(installYamlBytes),
	}

	anotherDownloadURL := "https://another.download.newrelic.com/"
	os.Setenv("NEW_RELIC_DOWNLOAD_URL", anotherDownloadURL)

	logFilePath := ".newrelic/newrelic-cli.log"
	os.Setenv("NEW_RELIC_CLI_LOG_FILE_PATH", logFilePath)

	clusterName := "sweet-cluster-name"
	os.Setenv("NR_CLI_CLUSTERNAME", clusterName)

	v, err := e.Prepare(m, r, false, "testLicenseKey")
	require.NoError(t, err)
	require.Contains(t, v["OS"], m.OS)
	assert.Contains(t, m.Platform, v["Platform"])
	assert.Contains(t, m.PlatformVersion, v["PlatformVersion"])
	assert.Contains(t, m.PlatformFamily, v["PlatformFamily"])
	assert.Contains(t, m.KernelArch, v["KernelArch"])
	assert.Contains(t, m.KernelVersion, v["KernelVersion"])
	assert.Equal(t, anotherDownloadURL, v["NEW_RELIC_DOWNLOAD_URL"])
	assert.Contains(t, v["NEW_RELIC_CLI_LOG_FILE_PATH"], logFilePath)
	assert.Equal(t, v["NR_CLI_CLUSTERNAME"], clusterName)
}

func Test_yamlFromJSON_convertsValidJsonToYaml(t *testing.T) {
	json := "{\"customAttribute_1\":\"SOME_ATTRIBUTE\",\"customAttribute_2\": \"SOME_ATTRIBUTE_2\"}"

	yaml := yamlFromJSON("key", json, []string{})

	assert.Contains(t, yaml, "custom_attributes:\n")
	assert.Contains(t, yaml, " customAttribute_1: SOME_ATTRIBUTE\n")
	assert.Contains(t, yaml, " customAttribute_2: SOME_ATTRIBUTE_2\n")
}

func Test_yamlFromJSON_convertsValidJsonToYamlWithInvalidTags(t *testing.T) {
	json := "{\"customAttribute_1\":\"SOME_ATTRIBUTE\",\"customAttribute_2\": \"SOME_ATTRIBUTE_2\"}"

	yaml := yamlFromJSON("key", json, []string{"abc"})

	assert.Contains(t, yaml, "custom_attributes:\n")
	assert.Contains(t, yaml, " customAttribute_1: SOME_ATTRIBUTE\n")
	assert.Contains(t, yaml, " customAttribute_2: SOME_ATTRIBUTE_2\n")
	assert.NotContains(t, yaml, "abc")
}

func Test_yamlFromJSON_convertsValidJsonToYamlWithSomeValidTags(t *testing.T) {
	json := "{\"customAttribute_1\":\"SOME_ATTRIBUTE\",\"customAttribute_2\": \"SOME_ATTRIBUTE_2\"}"

	yaml := yamlFromJSON("key", json, []string{"tag1:abc", "tag2:efg", "nocolontag"})

	assert.Contains(t, yaml, "custom_attributes:\n")
	assert.Contains(t, yaml, " customAttribute_1: SOME_ATTRIBUTE\n")
	assert.Contains(t, yaml, " customAttribute_2: SOME_ATTRIBUTE_2\n")
	assert.Contains(t, yaml, " tag1: abc\n")
	assert.Contains(t, yaml, " tag2: efg\n")
	assert.NotContains(t, yaml, " nocolontag\n")
}

func Test_yamlFromJSON_convertsValidJsonToYamTagWithMultipleValue(t *testing.T) {
	json := "{\"customAttribute_1\":\"SOME_ATTRIBUTE\",\"customAttribute_2\": \"SOME_ATTRIBUTE_2\"}"

	yaml := yamlFromJSON("key", json, []string{"tag1:abc def"})

	assert.Contains(t, yaml, "custom_attributes:\n")
	assert.Contains(t, yaml, " customAttribute_1: SOME_ATTRIBUTE\n")
	assert.Contains(t, yaml, " customAttribute_2: SOME_ATTRIBUTE_2\n")
	assert.Contains(t, yaml, " tag1: abc def\n")
}

func Test_yamlFromJSON_ConvertsInvalidJsonToEmptyValidTags(t *testing.T) {
	json := ""
	yaml := yamlFromJSON("key", json, []string{"tag1:abc"})

	assert.Contains(t, yaml, "custom_attributes:\n")
	assert.Contains(t, yaml, " tag1: abc\n")
}

func Test_yamlFromJSON_ConvertsInvalidJsonToEmpty(t *testing.T) {
	assert.Equal(t, "", yamlFromJSON("key", "totally-not-valid; json", []string{}))
}

func Test_yamlFromJSON_ConvertsEmptyStringToEmptyYaml(t *testing.T) {
	assert.Equal(t, "", yamlFromJSON("key", "", []string{}))
}

func Test_yamlFromCommaDelimitedString_convertsStringToYaml(t *testing.T) {
	yaml := yamlFromCommaDelimitedString("key", "value1,value2, value3")

	assert.Contains(t, yaml, "passthrough_environment:\n  - value1\n  - value2\n  - value3\n")
}

func Test_yamlFromCommaDelimitedString_convertsNonCSVStringToYaml(t *testing.T) {
	yaml := yamlFromCommaDelimitedString("key", "value1")

	assert.Contains(t, yaml, "passthrough_environment:\n")
	assert.Contains(t, yaml, "- value1\n")
}

func Test_yamlFromCommaDelimitedString_ConvertsEmptyStringToEmptyYaml(t *testing.T) {
	assert.Equal(t, "", yamlFromCommaDelimitedString("key", ""))
}

func TestRecipeVarProvider_OverrideDownloadURL(t *testing.T) {
	e := NewRecipeVarProvider()

	m := types.DiscoveryManifest{
		Hostname:        "testHostname",
		OS:              "testOS",
		Platform:        "testPlatform",
		PlatformFamily:  "testPlatformFamily",
		PlatformVersion: "testPlatformVersion",
		KernelArch:      "testKernelArch",
		KernelVersion:   "testKernelVersion",
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), t.Name())
	if err != nil {
		t.Fatal("error creating temp file")
	}

	defer os.Remove(tmpFile.Name())

	output := `
  {
    \"hostname\": \"{{.HOSTNAME}}\",
    \"os\": \"{{.OS}}\",
    \"platform\": \"{{.PLATFORM}}\",
    \"platformFamily\": \"{{.PLATFORM_FAMILY}}\",
    \"platformVersion\": \"{{.PLATFORM_VERSION}}\",
    \"kernelArch\": \"{{.KERNEL_ARCH}}\",
    \"kernelVersion\": \"{{.KERNEL_VERSION}}\"
  }`

	// We convert the `install` section of the recipe to a YAML string,
	// which is then used to create a Taskfile for go-task.
	recipeInstallToYaml := map[string]interface{}{
		"version": "3",
		"tasks": taskfile.Tasks{
			"default": &taskfile.Task{
				Cmds: []*taskfile.Cmd{
					{
						Cmd: fmt.Sprintf("echo %s > %s", strings.ReplaceAll(output, "\n", ""), tmpFile.Name()),
					},
				},
				Silent: true,
			},
		},
	}

	installYamlBytes, err := yaml.Marshal(recipeInstallToYaml)
	require.NoError(t, err)

	installYaml := string(installYamlBytes)

	r := types.OpenInstallationRecipe{
		Install: installYaml,
	}

	// Test for NEW_RELIC_DOWNLOAD_URL
	os.Setenv("NEW_RELIC_DOWNLOAD_URL", "https://another.download.newrelic.com/")

	v, err := e.Prepare(m, r, false, "testLicenseKey")
	require.NoError(t, err)
	require.Contains(t, "https://another.download.newrelic.com/", v["NEW_RELIC_DOWNLOAD_URL"])
}

func TestRecipeVarProvider_AllowInfraStaging(t *testing.T) {
	e := NewRecipeVarProvider()

	m := types.DiscoveryManifest{
		Hostname:        "testHostname",
		OS:              "testOS",
		Platform:        "testPlatform",
		PlatformFamily:  "testPlatformFamily",
		PlatformVersion: "testPlatformVersion",
		KernelArch:      "testKernelArch",
		KernelVersion:   "testKernelVersion",
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), t.Name())
	if err != nil {
		t.Fatal("error creating temp file")
	}

	defer os.Remove(tmpFile.Name())

	output := `
  {
    \"hostname\": \"{{.HOSTNAME}}\",
    \"os\": \"{{.OS}}\",
    \"platform\": \"{{.PLATFORM}}\",
    \"platformFamily\": \"{{.PLATFORM_FAMILY}}\",
    \"platformVersion\": \"{{.PLATFORM_VERSION}}\",
    \"kernelArch\": \"{{.KERNEL_ARCH}}\",
    \"kernelVersion\": \"{{.KERNEL_VERSION}}\"
  }`

	// We convert the `install` section of the recipe to a YAML string,
	// which is then used to create a Taskfile for go-task.
	recipeInstallToYaml := map[string]interface{}{
		"version": "3",
		"tasks": taskfile.Tasks{
			"default": &taskfile.Task{
				Cmds: []*taskfile.Cmd{
					{
						Cmd: fmt.Sprintf("echo %s > %s", strings.ReplaceAll(output, "\n", ""), tmpFile.Name()),
					},
				},
				Silent: true,
			},
		},
	}

	installYamlBytes, err := yaml.Marshal(recipeInstallToYaml)
	require.NoError(t, err)

	installYaml := string(installYamlBytes)

	r := types.OpenInstallationRecipe{
		Install: installYaml,
	}

	// Test for NEW_RELIC_DOWNLOAD_URL
	os.Setenv("NEW_RELIC_DOWNLOAD_URL", "https://nr-downloads-ohai-staging.s3.amazonaws.com/")

	v, err := e.Prepare(m, r, false, "testLicenseKey")
	require.NoError(t, err)
	require.Contains(t, "https://nr-downloads-ohai-staging.s3.amazonaws.com/", v["NEW_RELIC_DOWNLOAD_URL"])
}

func TestRecipeVarProvider_AllowInfraTesting(t *testing.T) {
	e := NewRecipeVarProvider()

	m := types.DiscoveryManifest{
		Hostname:        "testHostname",
		OS:              "testOS",
		Platform:        "testPlatform",
		PlatformFamily:  "testPlatformFamily",
		PlatformVersion: "testPlatformVersion",
		KernelArch:      "testKernelArch",
		KernelVersion:   "testKernelVersion",
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), t.Name())
	if err != nil {
		t.Fatal("error creating temp file")
	}

	defer os.Remove(tmpFile.Name())

	output := `
  {
    \"hostname\": \"{{.HOSTNAME}}\",
    \"os\": \"{{.OS}}\",
    \"platform\": \"{{.PLATFORM}}\",
    \"platformFamily\": \"{{.PLATFORM_FAMILY}}\",
    \"platformVersion\": \"{{.PLATFORM_VERSION}}\",
    \"kernelArch\": \"{{.KERNEL_ARCH}}\",
    \"kernelVersion\": \"{{.KERNEL_VERSION}}\"
  }`

	// We convert the `install` section of the recipe to a YAML string,
	// which is then used to create a Taskfile for go-task.
	recipeInstallToYaml := map[string]interface{}{
		"version": "3",
		"tasks": taskfile.Tasks{
			"default": &taskfile.Task{
				Cmds: []*taskfile.Cmd{
					{
						Cmd: fmt.Sprintf("echo %s > %s", strings.ReplaceAll(output, "\n", ""), tmpFile.Name()),
					},
				},
				Silent: true,
			},
		},
	}

	installYamlBytes, err := yaml.Marshal(recipeInstallToYaml)
	require.NoError(t, err)

	installYaml := string(installYamlBytes)

	r := types.OpenInstallationRecipe{
		Install: installYaml,
	}

	// Test for NEW_RELIC_DOWNLOAD_URL
	os.Setenv("NEW_RELIC_DOWNLOAD_URL", "https://nr-downloads-ohai-testing.s3.amazonaws.com/")

	v, err := e.Prepare(m, r, false, "testLicenseKey")
	require.NoError(t, err)
	require.Contains(t, "https://nr-downloads-ohai-testing.s3.amazonaws.com/", v["NEW_RELIC_DOWNLOAD_URL"])
}

func TestRecipeVarProvider_DisallowUnknownInfraTesting(t *testing.T) {
	e := NewRecipeVarProvider()

	m := types.DiscoveryManifest{
		Hostname:        "testHostname",
		OS:              "testOS",
		Platform:        "testPlatform",
		PlatformFamily:  "testPlatformFamily",
		PlatformVersion: "testPlatformVersion",
		KernelArch:      "testKernelArch",
		KernelVersion:   "testKernelVersion",
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), t.Name())
	if err != nil {
		t.Fatal("error creating temp file")
	}

	defer os.Remove(tmpFile.Name())

	output := `
  {
    \"hostname\": \"{{.HOSTNAME}}\",
    \"os\": \"{{.OS}}\",
    \"platform\": \"{{.PLATFORM}}\",
    \"platformFamily\": \"{{.PLATFORM_FAMILY}}\",
    \"platformVersion\": \"{{.PLATFORM_VERSION}}\",
    \"kernelArch\": \"{{.KERNEL_ARCH}}\",
    \"kernelVersion\": \"{{.KERNEL_VERSION}}\"
  }`

	// We convert the `install` section of the recipe to a YAML string,
	// which is then used to create a Taskfile for go-task.
	recipeInstallToYaml := map[string]interface{}{
		"version": "3",
		"tasks": taskfile.Tasks{
			"default": &taskfile.Task{
				Cmds: []*taskfile.Cmd{
					{
						Cmd: fmt.Sprintf("echo %s > %s", strings.ReplaceAll(output, "\n", ""), tmpFile.Name()),
					},
				},
				Silent: true,
			},
		},
	}

	installYamlBytes, err := yaml.Marshal(recipeInstallToYaml)
	require.NoError(t, err)

	installYaml := string(installYamlBytes)

	r := types.OpenInstallationRecipe{
		Install: installYaml,
	}

	// Test for NEW_RELIC_DOWNLOAD_URL
	os.Setenv("NEW_RELIC_DOWNLOAD_URL", "https://nr-downloads-ohai-unknown.s3.amazonaws.com/")

	v, err := e.Prepare(m, r, false, "testLicenseKey")
	require.NoError(t, err)
	require.Contains(t, "https://download.newrelic.com/", v["NEW_RELIC_DOWNLOAD_URL"])
}

func TestRecipeVarProvider_OverrideDownloadURL_RefusedNotHttps(t *testing.T) {
	e := NewRecipeVarProvider()

	m := types.DiscoveryManifest{
		Hostname:        "testHostname",
		OS:              "testOS",
		Platform:        "testPlatform",
		PlatformFamily:  "testPlatformFamily",
		PlatformVersion: "testPlatformVersion",
		KernelArch:      "testKernelArch",
		KernelVersion:   "testKernelVersion",
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), t.Name())
	if err != nil {
		t.Fatal("error creating temp file")
	}

	defer os.Remove(tmpFile.Name())

	output := `
  {
    \"hostname\": \"{{.HOSTNAME}}\",
    \"os\": \"{{.OS}}\",
    \"platform\": \"{{.PLATFORM}}\",
    \"platformFamily\": \"{{.PLATFORM_FAMILY}}\",
    \"platformVersion\": \"{{.PLATFORM_VERSION}}\",
    \"kernelArch\": \"{{.KERNEL_ARCH}}\",
    \"kernelVersion\": \"{{.KERNEL_VERSION}}\"
  }`

	// We convert the `install` section of the recipe to a YAML string,
	// which is then used to create a Taskfile for go-task.
	recipeInstallToYaml := map[string]interface{}{
		"version": "3",
		"tasks": taskfile.Tasks{
			"default": &taskfile.Task{
				Cmds: []*taskfile.Cmd{
					{
						Cmd: fmt.Sprintf("echo %s > %s", strings.ReplaceAll(output, "\n", ""), tmpFile.Name()),
					},
				},
				Silent: true,
			},
		},
	}

	installYamlBytes, err := yaml.Marshal(recipeInstallToYaml)
	require.NoError(t, err)

	installYaml := string(installYamlBytes)

	r := types.OpenInstallationRecipe{
		Install: installYaml,
	}

	// Test for NEW_RELIC_DOWNLOAD_URL
	os.Setenv("NEW_RELIC_DOWNLOAD_URL", "http://another.download.newrelic.com/")

	v, err := e.Prepare(m, r, false, "testLicenseKey")
	require.NoError(t, err)
	require.Contains(t, "https://download.newrelic.com/", v["NEW_RELIC_DOWNLOAD_URL"])
}

func TestRecipeVarProvider_OverrideDownloadURL_RefusedNotNewRelic(t *testing.T) {
	e := NewRecipeVarProvider()

	m := types.DiscoveryManifest{
		Hostname:        "testHostname",
		OS:              "testOS",
		Platform:        "testPlatform",
		PlatformFamily:  "testPlatformFamily",
		PlatformVersion: "testPlatformVersion",
		KernelArch:      "testKernelArch",
		KernelVersion:   "testKernelVersion",
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), t.Name())
	if err != nil {
		t.Fatal("error creating temp file")
	}

	defer os.Remove(tmpFile.Name())

	output := `
  {
    \"hostname\": \"{{.HOSTNAME}}\",
    \"os\": \"{{.OS}}\",
    \"platform\": \"{{.PLATFORM}}\",
    \"platformFamily\": \"{{.PLATFORM_FAMILY}}\",
    \"platformVersion\": \"{{.PLATFORM_VERSION}}\",
    \"kernelArch\": \"{{.KERNEL_ARCH}}\",
    \"kernelVersion\": \"{{.KERNEL_VERSION}}\"
  }`

	// We convert the `install` section of the recipe to a YAML string,
	// which is then used to create a Taskfile for go-task.
	recipeInstallToYaml := map[string]interface{}{
		"version": "3",
		"tasks": taskfile.Tasks{
			"default": &taskfile.Task{
				Cmds: []*taskfile.Cmd{
					{
						Cmd: fmt.Sprintf("echo %s > %s", strings.ReplaceAll(output, "\n", ""), tmpFile.Name()),
					},
				},
				Silent: true,
			},
		},
	}

	installYamlBytes, err := yaml.Marshal(recipeInstallToYaml)
	require.NoError(t, err)

	installYaml := string(installYamlBytes)

	r := types.OpenInstallationRecipe{
		Install: installYaml,
	}

	// Test for NEW_RELIC_DOWNLOAD_URL
	os.Setenv("NEW_RELIC_DOWNLOAD_URL", "http://github.com/")

	v, err := e.Prepare(m, r, false, "testLicenseKey")
	require.NoError(t, err)
	require.Contains(t, "https://download.newrelic.com/", v["NEW_RELIC_DOWNLOAD_URL"])
}
