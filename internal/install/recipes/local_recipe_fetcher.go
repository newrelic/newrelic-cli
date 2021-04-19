package recipes

import (
	"context"
	"embed"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

//go:embed recipes/*
var localRecipes embed.FS

type LocalRecipeFetcher struct{}

func (f *LocalRecipeFetcher) FetchRecipe(ctx context.Context, manifest *types.DiscoveryManifest, friendlyName string) (*types.Recipe, error) {

	return nil, nil
}

func (f *LocalRecipeFetcher) FetchRecommendations(ctx context.Context, manifest *types.DiscoveryManifest) ([]types.Recipe, error) {

	return nil, nil
}

func (f *LocalRecipeFetcher) FetchRecipes(ctx context.Context, manifest *types.DiscoveryManifest) ([]types.Recipe, error) {
	var recipes []types.Recipe

	files := recurseDirectory("recipes/newrelic/infrastructure")

	for _, file := range files {
		var r types.Recipe

		content, err := localRecipes.ReadFile(file)
		if err != nil {
			return nil, errors.Wrap(err, "error reading recipe file")
		}

		err = yaml.Unmarshal(content, &r)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshaling recipe file into Recipe")
		}

		recipes = append(recipes, r)
	}

	log.Debugf("%d embedded recipes found", len(recipes))

	return recipes, nil
}

func recurseDirectory(startDir string) []string {
	log.Debugf("recursing %s", startDir)
	var fileNames []string
	results, err := localRecipes.ReadDir(startDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range results {
		if r.Type().IsDir() {
			files := recurseDirectory(filepath.Join(startDir, r.Name()))
			fileNames = append(fileNames, files...)
		}

		if r.Type().IsRegular() {
			fileNames = append(fileNames, filepath.Join(startDir, r.Name()))
		}
	}

	return fileNames
}
