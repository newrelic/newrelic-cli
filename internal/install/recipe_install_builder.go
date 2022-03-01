package install

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/diagnose"
	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeInstallBuilder struct {
	configValidator   *diagnose.MockConfigValidator
	recipeFetcher     *recipes.MockRecipeFetcher
	discoverer        *discovery.MockDiscoverer
	status            *execution.InstallStatus
	mockOsValidator   *discovery.MockOsValidator
	manifestValidator *discovery.ManifestValidator
	licenseKeyFetcher *MockLicenseKeyFetcher
	shouldInstallCore func() bool
	bundler           RecipeBundler
	bundleInstaller   RecipeBundleInstaller
	installerContext  types.InstallerContext
	recipeVarProvider *execution.MockRecipeVarProvider
	recipeExecutor    *execution.MockRecipeExecutor
}

func NewRecipeInstallBuilder() *RecipeInstallBuilder {

	rib := &RecipeInstallBuilder{
		configValidator: diagnose.NewMockConfigValidator(),
		recipeFetcher:   recipes.NewMockRecipeFetcher(),
		discoverer:      discovery.NewMockDiscoverer(),
	}

	statusReporter := execution.NewMockStatusReporter()
	statusReporters := []execution.StatusSubscriber{statusReporter}
	status := execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rib.status = status

	rib.mockOsValidator = discovery.NewMockOsValidator()
	rib.manifestValidator = discovery.NewMockManifestValidator(rib.mockOsValidator)
	// Default to not skip core
	rib.shouldInstallCore = func() bool { return true }
	rib.installerContext = types.InstallerContext{}
	rib.licenseKeyFetcher = NewMockLicenseKeyFetcher()
	rib.recipeVarProvider = execution.NewMockRecipeVarProvider()
	rib.recipeVarProvider.Vars = map[string]string{}
	rib.recipeExecutor = execution.NewMockRecipeExecutor()

	return rib
}

func (rib *RecipeInstallBuilder) WithLibraryVersion(libraryVersion string) *RecipeInstallBuilder {
	rib.recipeFetcher.LibraryVersion = libraryVersion
	return rib
}

func (rib *RecipeInstallBuilder) WithLicenseKeyFetchResult(result error) *RecipeInstallBuilder {
	rib.licenseKeyFetcher.FetchLicenseKeyFunc = func(ctx context.Context) (string, error) {
		return "", result
	}
	return rib
}

func (rib *RecipeInstallBuilder) WithConfigValidatorError(err error) *RecipeInstallBuilder {
	rib.configValidator.Error = err
	return rib
}

func (rib *RecipeInstallBuilder) WithDiscovererError(err error) *RecipeInstallBuilder {
	rib.discoverer.Error = err
	return rib
}

func (rib *RecipeInstallBuilder) WithStatusReporter(statusReporter *execution.MockStatusSubscriber) *RecipeInstallBuilder {

	statusReporters := []execution.StatusSubscriber{statusReporter}
	status := execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rib.status = status
	return rib
}

func (rib *RecipeInstallBuilder) WithDiscovererValidatorError(err error) *RecipeInstallBuilder {
	rib.mockOsValidator.Error = err
	return rib
}

func (rib *RecipeInstallBuilder) withShouldInstallCore(shouldSkipCore func() bool) *RecipeInstallBuilder {
	rib.shouldInstallCore = shouldSkipCore
	return rib
}

func (rib *RecipeInstallBuilder) WithBundler(bundler RecipeBundler) *RecipeInstallBuilder {
	rib.bundler = bundler
	return rib
}

func (rib *RecipeInstallBuilder) WithBundleInstaller(bundlerInstaller RecipeBundleInstaller) *RecipeInstallBuilder {
	rib.bundleInstaller = bundlerInstaller
	return rib
}

func (rib *RecipeInstallBuilder) WithTargetRecipeName(name string) *RecipeInstallBuilder {
	rib.installerContext.RecipeNames = append(rib.installerContext.RecipeNames, name)
	return rib
}

func (rib *RecipeInstallBuilder) WithRecipeExecutionResult(err error) *RecipeInstallBuilder {
	rib.recipeExecutor.ExecuteErr = err
	return rib
}

func (rib *RecipeInstallBuilder) WithRecipeVarValues(vars map[string]string, err error) *RecipeInstallBuilder {
	rib.recipeVarProvider.Vars = vars
	rib.recipeVarProvider.Error = err
	return rib
}

func (rib *RecipeInstallBuilder) Build() *RecipeInstall {
	recipeInstall := &RecipeInstall{}
	recipeInstall.discoverer = rib.discoverer
	recipeInstall.configValidator = rib.configValidator
	recipeInstall.recipeFetcher = rib.recipeFetcher
	recipeInstall.status = rib.status
	recipeInstall.manifestValidator = rib.manifestValidator
	recipeInstall.bundlerFactory = func(ctx context.Context, repo *recipes.RecipeRepository) RecipeBundler {
		return rib.bundler
	}
	recipeInstall.bundleInstallerFactory = func(ctx context.Context, manifest *types.DiscoveryManifest, recipeInstallerInterface RecipeInstaller, statusReporter StatusReporter) RecipeBundleInstaller {
		return rib.bundleInstaller
	}
	recipeInstall.shouldInstallCore = rib.shouldInstallCore
	recipeInstall.InstallerContext = rib.installerContext
	recipeInstall.licenseKeyFetcher = rib.licenseKeyFetcher
	recipeInstall.recipeVarPreparer = rib.recipeVarProvider
	recipeInstall.recipeExecutor = rib.recipeExecutor

	return recipeInstall
}
