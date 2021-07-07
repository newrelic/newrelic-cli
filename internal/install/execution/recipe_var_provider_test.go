// +build integration

package execution

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/go-task/task/v3/taskfile"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestRecipeVarProvider_Basic(t *testing.T) {
	e := NewRecipeVarProvider()

	licenseKey := "testLicenseKey"

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

	v, err := e.Prepare(m, r, false, licenseKey)
	require.NoError(t, err)
	require.Contains(t, m.OS, v["OS"])
	require.Contains(t, m.Platform, v["Platform"])
	require.Contains(t, m.PlatformVersion, v["PlatformVersion"])
	require.Contains(t, m.PlatformFamily, v["PlatformFamily"])
	require.Contains(t, m.KernelArch, v["KernelArch"])
	require.Contains(t, m.KernelVersion, v["KernelVersion"])
	require.Contains(t, "https://download.newrelic.com/", v["NEW_RELIC_DOWNLOAD_URL"])
}

func TestRecipeVarProvider_OverrideDownloadURL(t *testing.T) {
	e := NewRecipeVarProvider()

	licenseKey := "testLicenseKey"
	os.Setenv("NEW_RELIC_LICENSE_KEY", licenseKey)

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

	v, err := e.Prepare(m, r, false, licenseKey)
	require.NoError(t, err)
	require.Contains(t, "https://another.download.newrelic.com/", v["NEW_RELIC_DOWNLOAD_URL"])
}

func TestRecipeVarProvider_AllowInfraStaging(t *testing.T) {
	e := NewRecipeVarProvider()

	licenseKey := "testLicenseKey"
	p := credentials.Profile{
		LicenseKey: "",
	}
	credentials.SetDefaultProfile(p)

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
	os.Setenv("NEW_RELIC_DOWNLOAD_URL", "https://nr-downloads-ohai-staging.s3-website-us-east-1.amazonaws.com/")

	v, err := e.Prepare(m, r, false, licenseKey)
	require.NoError(t, err)
	require.Contains(t, "https://nr-downloads-ohai-staging.s3-website-us-east-1.amazonaws.com/", v["NEW_RELIC_DOWNLOAD_URL"])
}

func TestRecipeVarProvider_AllowInfraTesting(t *testing.T) {
	e := NewRecipeVarProvider()

	licenseKey := "testLicenseKey"
	p := credentials.Profile{
		LicenseKey: "",
	}
	credentials.SetDefaultProfile(p)

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
	os.Setenv("NEW_RELIC_DOWNLOAD_URL", "https://nr-downloads-ohai-testing.s3-website-us-east-1.amazonaws.com/")

	v, err := e.Prepare(m, r, false, licenseKey)
	require.NoError(t, err)
	require.Contains(t, "https://nr-downloads-ohai-testing.s3-website-us-east-1.amazonaws.com/", v["NEW_RELIC_DOWNLOAD_URL"])
}

func TestRecipeVarProvider_DisallowUnknownInfraTesting(t *testing.T) {
	e := NewRecipeVarProvider()

	licenseKey := "testLicenseKey"
	p := credentials.Profile{
		LicenseKey: "",
	}
	credentials.SetDefaultProfile(p)

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
	os.Setenv("NEW_RELIC_DOWNLOAD_URL", "https://nr-downloads-ohai-unknown.s3-website-us-east-1.amazonaws.com/")

	v, err := e.Prepare(m, r, false, licenseKey)
	require.NoError(t, err)
	require.Contains(t, "https://download.newrelic.com/", v["NEW_RELIC_DOWNLOAD_URL"])
}

func TestRecipeVarProvider_OverrideDownloadURL_RefusedNotHttps(t *testing.T) {
	e := NewRecipeVarProvider()

	licenseKey := "testLicenseKey"
	os.Setenv("NEW_RELIC_LICENSE_KEY", licenseKey)

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

	v, err := e.Prepare(m, r, false, licenseKey)
	require.NoError(t, err)
	require.Contains(t, "https://download.newrelic.com/", v["NEW_RELIC_DOWNLOAD_URL"])
}

func TestRecipeVarProvider_OverrideDownloadURL_RefusedNotNewRelic(t *testing.T) {
	e := NewRecipeVarProvider()
	licenseKey := "testLicenseKey"
	os.Setenv("NEW_RELIC_LICENSE_KEY", licenseKey)

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

	v, err := e.Prepare(m, r, false, licenseKey)
	require.NoError(t, err)
	require.Contains(t, "https://download.newrelic.com/", v["NEW_RELIC_DOWNLOAD_URL"])
}
