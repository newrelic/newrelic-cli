// +build unit

package install

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
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
	if err != nil {
		t.Fatal("error creating temp file")
	}

	defer os.Remove(tmpFile.Name())

	ff := recipes.NewRecipeFileFetcher()

	f, err := ff.LoadRecipeFile(tmpFile.Name())
	require.NoError(t, err)
	require.NotNil(t, f)
}

func TestFetchRecipeFile_FailedStatusCode(t *testing.T) {
	ff := recipes.RecipeFileFetcherImpl{}

	makeHTTPGetFunc := func(statusCode int) func(string) (*http.Response, error) {
		return func(recipeURL string) (*http.Response, error) {
			return &http.Response{
				StatusCode: statusCode,
				Body:       ioutil.NopCloser(os.Stdin),
			}, nil
		}
	}

	u, err := url.Parse("https://localhost/404")
	require.NoError(t, err)

	ff.HTTPGetFunc = makeHTTPGetFunc(404)
	f, err := ff.FetchRecipeFile(u)
	require.Error(t, err)
	require.Nil(t, f)

	ff.HTTPGetFunc = makeHTTPGetFunc(199)
	f, err = ff.FetchRecipeFile(u)
	require.Error(t, err)
	require.Nil(t, f)

	ff.HTTPGetFunc = makeHTTPGetFunc(200)
	f, err = ff.FetchRecipeFile(u)
	require.NoError(t, err)
	require.NotNil(t, f)

	ff.HTTPGetFunc = makeHTTPGetFunc(299)
	f, err = ff.FetchRecipeFile(u)
	require.NoError(t, err)
	require.NotNil(t, f)
}

func TestNewRecipeFile(t *testing.T) {
	var expected recipes.RecipeFile
	err := yaml.Unmarshal([]byte(testRecipeFileString), &expected)
	require.NoError(t, err)

	actual, err := recipes.NewRecipeFile(testRecipeFileString)
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(&expected, actual))
}

func TestString(t *testing.T) {
	var f recipes.RecipeFile
	err := yaml.Unmarshal([]byte(testRecipeFileString), &f)
	require.NoError(t, err)

	s, err := f.String()
	require.NoError(t, err)
	require.NotEmpty(t, s)
}

func TestToRecipe(t *testing.T) {
	var f recipes.RecipeFile
	err := yaml.Unmarshal([]byte(testRecipeFileString), &f)
	require.NoError(t, err)

	r, err := f.ToRecipe()
	require.NoError(t, err)
	require.NotEmpty(t, r)
	require.NotEmpty(t, r.File)
	require.Equal(t, f.Name, r.Name)
	require.Equal(t, f.Description, r.Description)
	require.Equal(t, f.Repository, r.Repository)
	require.Equal(t, f.ValidationNRQL, r.ValidationNRQL)

	require.NotEmpty(t, f.Keywords, r.Keywords)
	require.NotEmpty(t, f.ProcessMatch, r.ProcessMatch)
}
