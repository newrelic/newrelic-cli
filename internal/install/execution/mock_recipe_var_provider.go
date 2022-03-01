package execution

import (
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockRecipeVarProvider struct {
	Vars  map[string]string
	Error error
}

func NewMockRecipeVarProvider() *MockRecipeVarProvider {
	return &MockRecipeVarProvider{}
}

func (rvp *MockRecipeVarProvider) Prepare(m types.DiscoveryManifest, r types.OpenInstallationRecipe, assumeYes bool, licenseKey string) (types.RecipeVars, error) {
	return rvp.Vars, rvp.Error
}
