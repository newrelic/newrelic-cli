package install

import (
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockRecipeDetector struct {
	availableRecipes   recipes.RecipeDetectionResults
	unavailableRecipes recipes.RecipeDetectionResults
	Err                error
}

type RecipeDetectionResult struct {
	Recipe     *types.OpenInstallationRecipe
	Status     execution.RecipeStatusType
	DurationMs int64
}

func (mrd *MockRecipeDetector) AddRecipeDetectionResult(detectionResult *recipes.RecipeDetectionResult) {
	if detectionResult.Status == execution.RecipeStatusTypes.AVAILABLE {
		mrd.availableRecipes = append(mrd.availableRecipes, detectionResult)
	} else {
		mrd.unavailableRecipes = append(mrd.unavailableRecipes, detectionResult)
	}
}

func (mrd *MockRecipeDetector) GetDetectedRecipes() (recipes.RecipeDetectionResults, recipes.RecipeDetectionResults, error) {
	if mrd.Err != nil {
		return nil, nil, mrd.Err
	}
	return mrd.availableRecipes, mrd.unavailableRecipes, nil
}
