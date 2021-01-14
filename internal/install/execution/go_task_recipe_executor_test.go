// +build integration

package execution

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/go-task/task/v3/taskfile"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-cli/internal/configuration"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestExecute_SystemVariableInterpolation(t *testing.T) {
	configuration.EnvVarResolver = configuration.NewMockEnvResolver()
	err := configuration.SetActiveProfileValue(configuration.LicenseKey, "testLicenseKey")
	require.NoError(t, err)

	e := NewGoTaskRecipeExecutor()

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

	f := recipes.RecipeFile{
		Install: map[string]interface{}{
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
		},
	}

	fs, err := yaml.Marshal(f)
	if err != nil {
		t.Fatal("could not marshal recipe file")
	}

	r := types.Recipe{
		File: string(fs),
	}

	v, err := e.Prepare(context.Background(), m, r, false)
	require.NoError(t, err)

	err = e.Execute(context.Background(), m, r, v)
	require.NoError(t, err)

	dat, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("error reading temp file %s", tmpFile.Name())
	}

	var actual types.DiscoveryManifest
	if err := json.Unmarshal(dat, &actual); err != nil {
		t.Fatalf("error unmarshaling temp file contents")
	}

	require.NotEmpty(t, string(dat))
	require.Equal(t, m.OS, actual.OS)
	require.Equal(t, m.Platform, actual.Platform)
	require.Equal(t, m.PlatformVersion, actual.PlatformVersion)
	require.Equal(t, m.PlatformFamily, actual.PlatformFamily)
	require.Equal(t, m.KernelArch, actual.KernelArch)
	require.Equal(t, m.KernelVersion, actual.KernelVersion)
}
