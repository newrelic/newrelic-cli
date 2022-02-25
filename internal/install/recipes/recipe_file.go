package recipes

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeFileFetcher struct {
	HTTPGetFunc  func(string) (*http.Response, error)
	readFileFunc func(string) ([]byte, error)
	Paths        []string
}

func (rff *RecipeFileFetcher) FetchLibraryVersion(ctx context.Context) string {
	return ""
}
func (rff *RecipeFileFetcher) FetchRecipes(ctx context.Context) ([]*types.OpenInstallationRecipe, error) {

	var recipesFromPath []*types.OpenInstallationRecipe

	for _, recipePath := range rff.Paths {
		recipeURL, parseErr := url.Parse(recipePath)
		isURL := parseErr == nil && recipeURL.Scheme != "" && strings.HasPrefix(strings.ToLower(recipeURL.Scheme), "http")
		var recipe *types.OpenInstallationRecipe
		var err error

		if isURL {
			recipe, err = rff.FetchRecipeFile(recipeURL)
			if err != nil {
				return recipesFromPath, fmt.Errorf("could not fetch file %s: %s", recipePath, err)
			}
		} else {
			recipe, err = rff.LoadRecipeFile(recipePath)
			if err != nil {
				return recipesFromPath, fmt.Errorf("could not load file %s: %s", recipePath, err)
			}
		}
		recipesFromPath = append(recipesFromPath, recipe)
	}

	return recipesFromPath, nil
}

func NewRecipeFileFetcher() *RecipeFileFetcher {
	f := RecipeFileFetcher{}
	f.HTTPGetFunc = defaultHTTPGetFunc
	f.readFileFunc = defaultReadFileFunc
	return &f
}

func defaultHTTPGetFunc(recipeURL string) (*http.Response, error) {
	return http.Get(recipeURL)
}

func defaultReadFileFunc(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func NewRecipeFile(recipeFileString string) (*types.OpenInstallationRecipe, error) {
	var f types.OpenInstallationRecipe
	err := yaml.Unmarshal([]byte(recipeFileString), &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (rff *RecipeFileFetcher) FetchRecipeFile(recipeURL *url.URL) (*types.OpenInstallationRecipe, error) {
	response, err := rff.HTTPGetFunc(recipeURL.String())

	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("received non-2xx Status code %d when retrieving recipe", response.StatusCode)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return NewRecipeFile(string(body))
}

func (rff *RecipeFileFetcher) LoadRecipeFile(filename string) (*types.OpenInstallationRecipe, error) {
	out, err := rff.readFileFunc(filename)
	if err != nil {
		return nil, err
	}

	return NewRecipeFile(string(out))
}
