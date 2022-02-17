package install

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func (mri *mockRecipeInstaller) promptIfNotLatestCLIVersion(ctx context.Context) error {
	args := mri.Called(ctx)
	return args.Error(0)
}

func (mri *mockRecipeInstaller) Install() error {
	args := mri.Called()
	return args.Error(0)
}

func (mri *mockRecipeInstaller) install(ctx context.Context) error {
	args := mri.Called(ctx)
	return args.Error(0)
}

func (mri *mockRecipeInstaller) assertDiscoveryValid(ctx context.Context, m *types.DiscoveryManifest) error {
	args := mri.Called(ctx, m)
	return args.Error(0)
}

func (mri *mockRecipeInstaller) discover(ctx context.Context) (*types.DiscoveryManifest, error) {
	args := mri.Called(ctx)
	// FIXME figure out how to return the object
	return nil, args.Error(1)
}

func (mri *mockRecipeInstaller) executeAndValidate(ctx context.Context, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, vars types.RecipeVars) (string, error) {
	args := mri.Called(ctx, m, r, vars)
	return args.String(0), args.Error(1)
}

func (mri *mockRecipeInstaller) validateRecipeViaAllMethods(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest, vars types.RecipeVars) (string, error) {
	args := mri.Called(ctx, r, m, vars)
	return args.String(0), args.Error(1)
}

func (mri *mockRecipeInstaller) executeAndValidateWithProgress(ctx context.Context, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, assumeYes bool) (string, error) {
	//FIXME I never get overridden, even though my interaction is mocked
	args := mri.Called(ctx, m, r, assumeYes)
	return args.String(0), args.Error(1)
}

type mockRecipeInstaller struct {
	mock.Mock
}

func (sr *mockStatusReporter) ReportStatus(status execution.RecipeStatusType, recipe types.OpenInstallationRecipe) {
	sr.counter++
}

type mockStatusReporter struct {
	counter int
}
