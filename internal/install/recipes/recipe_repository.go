package recipes

import (
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeRepository struct {
	RecipeLoaderFunc func() []types.OpenInstallationRecipe
}

// NewRecipeRepository returns a new instance of types.RecipeRepository.
func NewRecipeRepository(loaderFunc func() []types.OpenInstallationRecipe) *RecipeRepository {
	rr := RecipeRepository{
		RecipeLoaderFunc: loaderFunc,
	}

	return &rr
}

func (rf *RecipeRepository) FindAll(m types.DiscoveryManifest) []types.OpenInstallationRecipe {
	log.Debugf("Find all recipes available for host")
	return []types.OpenInstallationRecipe{}
}
