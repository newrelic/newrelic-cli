package uninstall

import (
	"context"
	"strings"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// RepositoryRecipeFinder resolves recipe names using the same host-filtered
// recipes.RecipeRepository the install command uses, so uninstall stays in
// sync with whatever recipes/InstallTargets install would use.
type RepositoryRecipeFinder struct {
	loaderFunc func() ([]*types.OpenInstallationRecipe, error)
	repo       *recipes.RecipeRepository
}

// NewRepositoryRecipeFinder returns a RecipeFinder backed by the given recipe loader and host manifest.
func NewRepositoryRecipeFinder(loaderFunc func() ([]*types.OpenInstallationRecipe, error), manifest *types.DiscoveryManifest) *RepositoryRecipeFinder {
	return &RepositoryRecipeFinder{
		loaderFunc: loaderFunc,
		repo:       recipes.NewRecipeRepository(loaderFunc, manifest),
	}
}

// FindRecipe returns the recipe if it's known and compatible with this host.
// It returns ErrRecipeUnsupportedOnHost if the recipe exists but targets a
// different host, or ErrRecipeNotFound if the name is unknown entirely.
func (f *RepositoryRecipeFinder) FindRecipe(ctx context.Context, name string) (*types.OpenInstallationRecipe, error) {
	if r := f.repo.FindRecipeByName(name); r != nil {
		return r, nil
	}

	all, err := f.loaderFunc()
	if err != nil {
		return nil, err
	}

	for _, r := range all {
		if strings.EqualFold(r.Name, name) {
			return nil, ErrRecipeUnsupportedOnHost
		}
	}

	return nil, ErrRecipeNotFound
}
