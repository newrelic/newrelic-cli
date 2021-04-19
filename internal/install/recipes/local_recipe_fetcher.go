package recipes

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

//go:embed recipes/*
var localRecipes embed.FS

// ErrRecipeNotFound is used when a recipe is requested by name, but does not exist for the given constraint.
var ErrRecipeNotFound = errors.New("recipe not found")

// ErrInvalidRecipeFile is used when a recipe file fails to Unmarshal into a Recipe.
var ErrInvalidRecipeFile = errors.New("invalid recipe file")

type LocalRecipeFetcher struct{}

// FetchRecipe fetches a recommended recipe by name.
func (f *LocalRecipeFetcher) FetchRecipe(ctx context.Context, manifest *types.DiscoveryManifest, friendlyName string) (*types.Recipe, error) {

	recipes, err := f.FetchRecommendations(ctx, manifest)
	if err != nil {
		return nil, err
	}

	for _, recipe := range recipes {
		if recipe.Name == friendlyName {
			return &recipe, nil
		}
	}

	return nil, fmt.Errorf("%s: %w", friendlyName, ErrRecipeNotFound)
}

// FetchRecommendations fetches the recipes based on the manifest constraints.
func (f *LocalRecipeFetcher) FetchRecommendations(ctx context.Context, manifest *types.DiscoveryManifest) ([]types.Recipe, error) {

	recipes, err := f.FetchRecipes(ctx, manifest)
	if err != nil {
		return nil, err
	}

	return manifest.ConstrainRecipes(recipes), nil
}

// FetchRecipes fetches all recipes.
func (f *LocalRecipeFetcher) FetchRecipes(ctx context.Context, manifest *types.DiscoveryManifest) ([]types.Recipe, error) {
	var recipes []types.Recipe

	files, err := recurseDirectory("recipes/newrelic/infrastructure")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		var r types.Recipe

		content, err := localRecipes.ReadFile(file)
		if err != nil {
			return nil, errors.Wrap(err, "error reading recipe file")
		}

		err = yaml.Unmarshal(content, &r)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", file, ErrInvalidRecipeFile)
		}

		recipes = append(recipes, r)
	}

	return recipes, nil
}

func recurseDirectory(startDir string) ([]string, error) {
	log.Debugf("recursing %s", startDir)
	var fileNames []string
	results, err := localRecipes.ReadDir(startDir)
	if err != nil {
		return nil, err
	}

	for _, r := range results {
		if r.Type().IsDir() {
			files, err := recurseDirectory(filepath.Join(startDir, r.Name()))
			if err != nil {
				return nil, err
			}

			fileNames = append(fileNames, files...)
		}

		if r.Type().IsRegular() {
			fileNames = append(fileNames, filepath.Join(startDir, r.Name()))
		}
	}

	return fileNames, nil
}
