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
	bundler           RecipeBundler
	bundleInstaller   RecipeBundleInstaller
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

	getEnvVariable = func(name string) string {
		return ""
	}

	return rib
}

func (rib *RecipeInstallBuilder) WithLibraryVersion(libraryVersion string) *RecipeInstallBuilder {
	rib.recipeFetcher.LibraryVersion = libraryVersion
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

func (rib *RecipeInstallBuilder) WithBundler(bundler RecipeBundler) *RecipeInstallBuilder {
	rib.bundler = bundler
	return rib
}

func (rib *RecipeInstallBuilder) WithBundleInstaller(bundlerInstaller RecipeBundleInstaller) *RecipeInstallBuilder {
	rib.bundleInstaller = bundlerInstaller
	return rib
}

func (rib *RecipeInstallBuilder) Build() *RecipeInstall {
	recipeInstall := &RecipeInstall{}
	recipeInstall.discoverer = rib.discoverer
	recipeInstall.configValidator = rib.configValidator
	recipeInstall.recipeFetcher = rib.recipeFetcher
	recipeInstall.status = rib.status
	recipeInstall.manifestValidator = rib.manifestValidator
	getBundler = func(ctx context.Context, repo *recipes.RecipeRepository) RecipeBundler {
		return rib.bundler
	}
	getBundleInstaller = func(ctx context.Context, manifest *types.DiscoveryManifest, recipeInstallerInterface RecipeInstaller, statusReporter StatusReporter) RecipeBundleInstaller {
		return rib.bundleInstaller
	}

	return recipeInstall
}
