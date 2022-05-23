package install

import "github.com/newrelic/newrelic-cli/internal/install/recipes"

type MockBundleInstaller struct {
	installedRecipes map[string]bool
	errors           map[string]error
}

func NewMockBundleInstaller() *MockBundleInstaller {
	return &MockBundleInstaller{
		installedRecipes: make(map[string]bool),
		errors:           make(map[string]error),
	}
}

func (mbi *MockBundleInstaller) InstallStopOnError(bundle *recipes.Bundle, assumeYes bool) error {
	for _, recipe := range bundle.BundleRecipes {
		err, foundError := mbi.errors[recipe.Recipe.Name]
		if !foundError {
			mbi.installedRecipes[recipe.Recipe.Name] = true
		} else {
			return err
		}
	}

	return nil
}

func (mbi *MockBundleInstaller) InstallContinueOnError(bundle *recipes.Bundle, assumeYes bool) {
	for _, recipe := range bundle.BundleRecipes {
		_, foundError := mbi.errors[recipe.Recipe.Name]
		if !foundError {
			mbi.installedRecipes[recipe.Recipe.Name] = true
		}
	}
}

func (mbi *MockBundleInstaller) WithInstallError(recipeName string, err error) {
	mbi.errors[recipeName] = err
}

func (mbi *MockBundleInstaller) InstalledRecipesCount() int {
	return len(mbi.installedRecipes)
}
