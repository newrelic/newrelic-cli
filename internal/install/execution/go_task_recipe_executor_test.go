// +build integration

package execution

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/go-task/task/v3/taskfile"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestExecute_SystemVariableInterpolation(t *testing.T) {
	p := credentials.Profile{
		LicenseKey:        "testLicenseKey",
		InsightsInsertKey: "testInsightsInsertKey",
		AccountID:         12345,
	}
	credentials.SetDefaultProfile(p)

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
		\"kernelVersion\": \"{{.KERNEL_VERSION}}\",
		\"accountID\": \"{{.NEW_RELIC_ACCOUNT_ID}}\",
		\"licenseKey\": \"{{.NEW_RELIC_LICENSE_KEY}}\",
		\"apiKey\": \"{{.NEW_RELIC_API_KEY}}\",
		\"region\": \"{{.NEW_RELIC_REGION}}\",
		\"insightsInsertKey\": \"{{.NEW_RELIC_INSIGHTS_INSERT_KEY}}\",
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

	var actualM types.DiscoveryManifest
	if err := json.Unmarshal(dat, &actualM); err != nil {
		t.Fatalf("error unmarshaling temp file contents: %s", err)
	}

	var actualP profile
	if err := json.Unmarshal(dat, &actualP); err != nil {
		t.Fatalf("error unmarshaling temp file contents: %s", err)
	}

	require.NotEmpty(t, string(dat))
	require.Equal(t, m.OS, actualM.OS)
	require.Equal(t, m.Platform, actualM.Platform)
	require.Equal(t, m.PlatformVersion, actualM.PlatformVersion)
	require.Equal(t, m.PlatformFamily, actualM.PlatformFamily)
	require.Equal(t, m.KernelArch, actualM.KernelArch)
	require.Equal(t, m.KernelVersion, actualM.KernelVersion)

	require.Equal(t, p.APIKey, actualP.APIKey)
	require.Equal(t, strconv.Itoa(p.AccountID), actualP.AccountID)
	require.Equal(t, p.LicenseKey, actualP.LicenseKey)
	require.Equal(t, p.Region, actualP.Region)
	require.Equal(t, p.InsightsInsertKey, actualP.InsightsInsertKey)
}

type profile struct {
	APIKey            string `json:"apiKey"`
	InsightsInsertKey string `json:"insightsInsertKey"`
	Region            string `json:"region"`
	AccountID         string `json:"accountID"`
	LicenseKey        string `json:"licenseKey"`
}
