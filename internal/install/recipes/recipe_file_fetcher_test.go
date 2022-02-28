package recipes

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var (
	testRecipeFileString = `
---
description: testDescription
keywords:
  - testKeyword
name: testName
processMatch:
  - testProcessMatch
repository: testRepository
validationNrql: testValidationNrql
inputVars:
  - name: testName
    prompt: testPrompt
    secret: true
    default: testDefault
installTargets:
  - type: testType
    os: testOS
    platform: testPlatform
    platformFamily: testPlatformFamily
    platformVersion: testPlatformVersion
    kernelVersion: testKerrnelVersion
    kernelArch: testKernelArch
logMatch:
  - name: testName
    file: testFile
    attributes:
      logtype: testlogtype
    pattern: testPattern
    systemd: testSystemd
`
)

func TestLoadRecipeFile(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), t.Name())
	_, _ = io.WriteString(tmpFile, testRecipeFileString)

	if err != nil {
		t.Fatal("error creating temp file")
	}

	defer os.Remove(tmpFile.Name())

	ff := NewRecipeFileFetcher()

	f, err := ff.LoadRecipeFile(tmpFile.Name())

	require.NoError(t, err)
	require.NotNil(t, f)
	require.Equal(t, "testName", f.Name)
}

func TestFetchRecipeFile_FailedStatusCode(t *testing.T) {
	type testCase struct {
		statusCode  int
		shouldError bool
	}

	stubbedHTTPGetFunction := func(statusCode int) func(string) (*http.Response, error) {
		return func(recipeURL string) (*http.Response, error) {
			return &http.Response{
				StatusCode: statusCode,
				Body:       ioutil.NopCloser(os.Stdin),
			}, nil
		}
	}

	ff := RecipeFileFetcher{}
	u, err := url.Parse("https://localhost/valid-url")
	assert.NoError(t, err)

	tests := []testCase{
		{statusCode: 404, shouldError: true},
		{statusCode: 199, shouldError: true},
		{statusCode: 200, shouldError: false},
		{statusCode: 299, shouldError: false},
	}

	for _, testCondition := range tests {
		ff.HTTPGetFunc = stubbedHTTPGetFunction(testCondition.statusCode)
		f, err := ff.FetchRecipeFile(u)

		switch testCondition.shouldError {
		case true:
			assert.Error(t, err)
			assert.Nil(t, f)

		case false:
			assert.NoError(t, err)
			assert.NotNil(t, f)
		}
	}
}

func TestNewRecipeFile(t *testing.T) {
	var expected types.OpenInstallationRecipe
	err := yaml.Unmarshal([]byte(testRecipeFileString), &expected)
	require.NoError(t, err)

	actual, err := NewRecipeFile(testRecipeFileString)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(&expected, actual))
}
