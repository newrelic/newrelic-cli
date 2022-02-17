package install

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/newrelic/newrelic-cli/internal/diagnose"
	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/install/validation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type bundleInstallerTest struct {
	installedRecipes map[string]bool
	ctx              context.Context
	manifest         *types.DiscoveryManifest
	recipeInstaller  *RecipeInstaller
	statusReporter   *mockStatusReporter
	bundleInstaller  *BundleInstaller
}

var (
	bundleInstallerTestImpl *bundleInstallerTest
)

type mockInstallBundleRecipe struct {
	mock.Mock
}

func setup(err error) {
	ic := types.InstallerContext{
		RecipePaths: []string{"testRecipePath"},
		RecipeNames: []string{"testRecipeName"},
	}

	d := discovery.NewMockDiscoverer()
	mv := discovery.NewEmptyManifestValidator()
	f := recipes.NewMockRecipeFetcher()
	e := execution.NewMockRecipeExecutor()

	if err != nil {
		e.ExecuteErr = err
	}

	v := validation.NewMockRecipeValidator()
	ff := recipes.NewMockRecipeFileFetcher()
	statusReporters := []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status := execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	p := ux.NewMockPrompter()
	pi := ux.NewMockProgressIndicator()
	sp := ux.NewMockProgressIndicator()
	lkf := NewMockLicenseKeyFetcher()
	cv := diagnose.NewMockConfigValidator()
	rvp := execution.NewRecipeVarProvider()
	av := validation.NewAgentValidator()
	rf := recipes.NewRecipeFilterRunner(ic, status)
	i := RecipeInstaller{ic, d, mv, f, e, v, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}

	manifest := types.DiscoveryManifest{
		DiscoveredProcesses: []types.GenericProcess{},
	}

	bundleInstallerTestImpl = &bundleInstallerTest{
		statusReporter:  &mockStatusReporter{},
		recipeInstaller: &i,
		manifest:        &manifest,
		ctx:             context.Background(),
	}

	bundleInstallerTestImpl.bundleInstaller = NewBundleInstaller(
		bundleInstallerTestImpl.ctx,
		bundleInstallerTestImpl.manifest,
		bundleInstallerTestImpl.recipeInstaller,
		bundleInstallerTestImpl.statusReporter)
}

// public functions
func TestBundleInstallerStopsOnError(t *testing.T) {
	expectedError := "I am an error"
	errorPrefix := "execution failed for : "
	setup(errors.New(expectedError))

	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe: recipes.NewRecipeBuilder().ID("0").Name("recipe1").Build(),
				RecipeStatuses: []recipes.RecipeStatus{
					{
						Status:     execution.RecipeStatusTypes.AVAILABLE,
						StatusTime: time.Now(),
					},
				},
			},
		},
	}

	actualError := bundleInstallerTestImpl.bundleInstaller.InstallStopOnError(&bundle, true)

	require.Equal(t, errorPrefix+expectedError, actualError.Error())
}

func TestBundleInstallerContinueOnError(t *testing.T) {
	expectedError := "I am an error"
	setup(errors.New(expectedError))

	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe: recipes.NewRecipeBuilder().ID("0").Name("recipe1").Build(),
				RecipeStatuses: []recipes.RecipeStatus{
					{
						Status:     execution.RecipeStatusTypes.AVAILABLE,
						StatusTime: time.Now(),
					},
				},
			},
		},
	}

	//TODO: Need to find out how to verify error was thrown
	bundleInstallerTestImpl.bundleInstaller.InstallContinueOnError(&bundle, true)
}

func TestBundleInstallerReportsStatus(t *testing.T) {
	setup(nil)
	bundle := givenBundle(types.InfraAgentRecipeName)
	bundle.BundleRecipes[0].AddStatus(execution.RecipeStatusTypes.AVAILABLE, time.Now())

	bundleInstallerTestImpl.bundleInstaller.reportStatus(bundle)

	actual := bundleInstallerTestImpl.statusReporter.counter
	expected := len(bundle.BundleRecipes[0].RecipeStatuses)
	require.Equal(t, expected, actual)

}

func givenBundle(recipeName string) *recipes.Bundle {
	bundle := &recipes.Bundle{}
	r := &types.OpenInstallationRecipe{
		Name: recipeName,
	}
	br := &recipes.BundleRecipe{
		Recipe: r,
	}
	bundle.AddRecipe(br)
	return bundle
}

// func TestBundleInstallerInstallsBundleRecipes(t *testing.T) {
// 	require.Fail(t, "Implement me")
// }

// func TestBundleInstallerInstallsBundleRecipesWithDependencies(t *testing.T) {
// 	require.Fail(t, "Implement me")
// }

type mockStatusReporter struct {
	counter int
}

func (sr *mockStatusReporter) ReportStatus(status execution.RecipeStatusType, recipe types.OpenInstallationRecipe) {
	sr.counter++
}
