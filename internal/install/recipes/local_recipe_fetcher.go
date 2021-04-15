package recipes

import (
	"context"
	"embed"

	log "github.com/sirupsen/logrus"

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
	recurseRecipeFiles()

	return nil, nil
}

func recurseRecipeFiles() {

	log.Warnf("content: %+v", localRecipes)

	dirs := recurseDirectories("recipes/newrelic/infrastructure")

	results, err := localRecipes.ReadDir("recipes/newrelic/infrastructure")
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range results {
		log.Warnf("name: %+v\n", r.Name())
	}
}

func recurseDirectories(startDir string) []string {

}
