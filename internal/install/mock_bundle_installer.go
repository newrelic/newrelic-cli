package install

import "github.com/newrelic/newrelic-cli/internal/install/recipes"

type MockBundleInstaller struct {
	installedRecipes map[string]bool
	Error            error
}

func NewMockBundleInstaller() *MockBundleInstaller {
	return &MockBundleInstaller{
		installedRecipes: make(map[string]bool),
	}
}

func (mbi *MockBundleInstaller) InstallStopOnError(bundle *recipes.Bundle, assumeYes bool) error {
	for _, recipe := range bundle.BundleRecipes {
		mbi.installedRecipes[recipe.Recipe.Name] = true
	}

	return mbi.Error
}
func (mbi *MockBundleInstaller) InstallContinueOnError(bundle *recipes.Bundle, assumeYes bool) {
	for _, recipe := range bundle.BundleRecipes {
		mbi.installedRecipes[recipe.Recipe.Name] = true
	}
}
func (mbi *MockBundleInstaller) InstalledRecipesCount() int {
	return len(mbi.installedRecipes)
}
