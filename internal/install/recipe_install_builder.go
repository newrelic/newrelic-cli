package install

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/validation/mocks"

	"github.com/newrelic/newrelic-cli/internal/diagnose"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
)

type RecipeInstallBuilder struct {
	configValidator    *diagnose.MockConfigValidator
	recipeFetcher      *recipes.MockRecipeFetcher
	status             *execution.InstallStatus
	licenseKeyFetcher  *MockLicenseKeyFetcher
	shouldInstallCore  func() bool
	installerContext   types.InstallerContext
	recipeLogForwarder *execution.MockRecipeLogForwarder
	recipeVarProvider  *execution.MockRecipeVarProvider
	recipeExecutor     *execution.MockRecipeExecutor
	progressIndicator  *ux.SpinnerProgressIndicator
	agentValidator     *mocks.MockAgentValidator
	recipeValidator    *mocks.MockRecipeValidator
	recipeDetector     *MockRecipeDetector
	processes          []types.GenericProcess
}

func NewRecipeInstallBuilder() *RecipeInstallBuilder {
	rib := &RecipeInstallBuilder{
		configValidator: diagnose.NewMockConfigValidator(),
		recipeFetcher:   recipes.NewMockRecipeFetcher(),
		processes:       []types.GenericProcess{},
	}

	statusReporter := execution.NewMockStatusReporter()
	statusReporters := []execution.StatusSubscriber{statusReporter}
	status := execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rib.status = status

	// Default to not skip core
	rib.shouldInstallCore = func() bool { return true }
	rib.installerContext = types.InstallerContext{}
	rib.licenseKeyFetcher = NewMockLicenseKeyFetcher()
	rib.recipeLogForwarder = execution.NewMockRecipeLogForwarder()
	rib.recipeVarProvider = execution.NewMockRecipeVarProvider()
	rib.recipeVarProvider.Vars = map[string]string{}
	rib.recipeExecutor = execution.NewMockRecipeExecutor()
	rib.progressIndicator = ux.NewSpinnerProgressIndicator()
	rib.agentValidator = &mocks.MockAgentValidator{}
	rib.recipeValidator = &mocks.MockRecipeValidator{}
	rib.recipeDetector = &MockRecipeDetector{}

	return rib
}

func (rib *RecipeInstallBuilder) WithLibraryVersion(libraryVersion string) *RecipeInstallBuilder {
	rib.recipeFetcher.LibraryVersion = libraryVersion
	return rib
}

func (rib *RecipeInstallBuilder) WithFetchRecipesVal(fetchRecipesVal []*types.OpenInstallationRecipe) *RecipeInstallBuilder {
	rib.recipeFetcher.FetchRecipesVal = fetchRecipesVal
	return rib
}

func (rib *RecipeInstallBuilder) WithRecipeDetectionResult(detectionResult *recipes.RecipeDetectionResult) *RecipeInstallBuilder {
	rib.recipeDetector.AddRecipeDetectionResult(detectionResult)
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

func (rib *RecipeInstallBuilder) WithStatusReporter(statusReporter *execution.MockStatusSubscriber) *RecipeInstallBuilder {
	statusReporters := []execution.StatusSubscriber{statusReporter}
	status := execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rib.status = status
	return rib
}

func (rib *RecipeInstallBuilder) withShouldInstallCore(shouldSkipCore func() bool) *RecipeInstallBuilder {
	rib.shouldInstallCore = shouldSkipCore
	return rib
}

func (rib *RecipeInstallBuilder) WithTargetRecipeName(name string) *RecipeInstallBuilder {
	rib.installerContext.RecipeNames = append(rib.installerContext.RecipeNames, name)
	return rib
}

func (rib *RecipeInstallBuilder) WithRecipeExecutionError(err error) *RecipeInstallBuilder {
	rib.recipeExecutor.ExecuteErr = err
	return rib
}

func (rib *RecipeInstallBuilder) WithOutput(value string) *RecipeInstallBuilder {
	rib.recipeExecutor.SetOutput(value)
	return rib
}

func (rib *RecipeInstallBuilder) WithRecipeOutput(value []string) *RecipeInstallBuilder {
	return rib
}

func (rib *RecipeInstallBuilder) WithRecipeVarValues(vars map[string]string, err error) *RecipeInstallBuilder {
	rib.recipeVarProvider.Vars = vars
	rib.recipeVarProvider.Error = err
	return rib
}

func (rib *RecipeInstallBuilder) WithRecipeLogForwarder(optIn bool) *RecipeInstallBuilder {
	rib.recipeLogForwarder.SetUserOptedIn(optIn)
	return rib
}

func (rib *RecipeInstallBuilder) WithProgressIndicator(i *ux.SpinnerProgressIndicator) *RecipeInstallBuilder {
	rib.progressIndicator = i
	return rib
}

func (rib *RecipeInstallBuilder) WithAgentValidationError(e error) *RecipeInstallBuilder {
	rib.agentValidator.Error = e
	return rib
}

func (rib *RecipeInstallBuilder) WithRecipeValidationError(e error) *RecipeInstallBuilder {
	rib.recipeValidator.Error = e
	return rib
}

func (rib *RecipeInstallBuilder) WithRunningProcess(cmd string, name string) *RecipeInstallBuilder {
	p := recipes.NewMockProcess(cmd, name, 0)
	rib.processes = append(rib.processes, p)
	return rib
}

func (rib *RecipeInstallBuilder) Build() *RecipeInstall {
	recipeInstall := &RecipeInstall{}
	recipeInstall.configValidator = rib.configValidator
	recipeInstall.recipeFetcher = rib.recipeFetcher
	recipeInstall.status = rib.status
	recipeInstall.bundlerFactory = func(ctx context.Context, detections recipes.RecipeDetectionResults) RecipeBundler {
		return recipes.NewBundler(ctx, detections)
	}
	recipeInstall.bundleInstallerFactory = func(ctx context.Context, manifest *types.DiscoveryManifest, recipeInstallerInterface RecipeInstaller, statusReporter StatusReporter) RecipeBundleInstaller {
		return NewBundleInstaller(context.Background(), &types.DiscoveryManifest{}, recipeInstall, rib.status)
	}
	recipeInstall.shouldInstallCore = rib.shouldInstallCore
	recipeInstall.InstallerContext = rib.installerContext
	recipeInstall.licenseKeyFetcher = rib.licenseKeyFetcher
	recipeInstall.recipeLogForwarder = rib.recipeLogForwarder
	recipeInstall.recipeVarPreparer = rib.recipeVarProvider
	recipeInstall.recipeExecutor = rib.recipeExecutor
	recipeInstall.progressIndicator = rib.progressIndicator
	recipeInstall.agentValidator = rib.agentValidator
	recipeInstall.recipeValidator = rib.recipeValidator
	recipeInstall.recipeDetectorFactory = func(ctx context.Context, repo *recipes.RecipeRepository) RecipeStatusDetector {
		return rib.recipeDetector
	}
	mockProcessEvaluator := recipes.NewMockProcessEvaluator()
	mockProcessEvaluator.WithProcesses(rib.processes)
	recipeInstall.processEvaluator = mockProcessEvaluator

	return recipeInstall
}
