package recipes

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestRecipeDetectorShouldFailBecauseOfProcessEvaluation(t *testing.T) {
	recipe := NewRecipeBuilder().Build()
	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorRecipeStatus(recipe, execution.RecipeStatusTypes.NULL)
	b.WithScriptEvaluatorRecipeStatus(recipe, execution.RecipeStatusTypes.AVAILABLE)
	detector := b.Build()

	_, ua, _ := detector.GetDetectedRecipes()
	actual, ok := ua.GetRecipeDetection(recipe.Name)

	require.True(t, ok)
	require.Equal(t, execution.RecipeStatusTypes.NULL, actual.Status)
}

func TestRecipeDetectorShouldBeAvailableWhenRecipeScriptDetectionIsMissingScript(t *testing.T) {
	recipe := NewRecipeBuilder().Build()
	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorRecipeStatus(recipe, execution.RecipeStatusTypes.AVAILABLE)
	b.WithScriptEvaluatorRecipeStatus(recipe, execution.RecipeStatusTypes.DETECTED)
	detector := b.Build()

	a, _, _ := detector.GetDetectedRecipes()
	actual, ok := a.GetRecipeDetection(recipe.Name)
	require.True(t, ok)
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, actual.Status)
}

func TestRecipeDetectorShouldFailWhenScriptFails(t *testing.T) {
	recipe := NewRecipeBuilder().WithPreInstallScript("pre-install script mock").Build()

	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorRecipeStatus(recipe, execution.RecipeStatusTypes.AVAILABLE)
	b.WithScriptEvaluatorRecipeStatus(recipe, execution.RecipeStatusTypes.NULL)
	detector := b.Build()

	_, ua, _ := detector.GetDetectedRecipes()
	actual, ok := ua.GetRecipeDetection(recipe.Name)

	require.True(t, ok)
	require.Equal(t, execution.RecipeStatusTypes.NULL, actual.Status)
}

func TestRecipeDetectorShouldDetectBecauseOfScriptEvaluation(t *testing.T) {
	recipe := NewRecipeBuilder().WithPreInstallScript("pre-install script mock").Build()

	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorRecipeStatus(recipe, execution.RecipeStatusTypes.AVAILABLE)
	b.WithScriptEvaluatorRecipeStatus(recipe, execution.RecipeStatusTypes.DETECTED)
	detector := b.Build()

	_, ua, _ := detector.GetDetectedRecipes()
	actual, ok := ua.GetRecipeDetection(recipe.Name)

	require.True(t, ok)
	require.Equal(t, execution.RecipeStatusTypes.DETECTED, actual.Status)
}

func TestRecipeDetectorShouldBeAvailableBecauseOfScriptEvaluation(t *testing.T) {
	recipe := NewRecipeBuilder().WithPreInstallScript("pre-install script mock").Build()

	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorRecipeStatus(recipe, execution.RecipeStatusTypes.AVAILABLE)
	b.WithScriptEvaluatorRecipeStatus(recipe, execution.RecipeStatusTypes.AVAILABLE)
	detector := b.Build()

	a, _, _ := detector.GetDetectedRecipes()
	actual, ok := a.GetRecipeDetection(recipe.Name)

	require.True(t, ok)
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, actual.Status)
}

func TestDetectionResultsShouldSortByRecipeName(t *testing.T) {
	detections := []*RecipeDetectionResult{}
	r1 := &types.OpenInstallationRecipe{
		Name: "b",
	}
	r2 := &types.OpenInstallationRecipe{
		Name: "a",
	}
	detections = append(detections, &RecipeDetectionResult{Recipe: r1})
	detections = append(detections, &RecipeDetectionResult{Recipe: r2})
	sort.Sort(RecipeDetectionResults(detections))

	require.Equal(t, detections[0].Recipe.Name, "a")
	require.Equal(t, detections[1].Recipe.Name, "b")
}

type MockRecipesFinder struct {
	recipes []*types.OpenInstallationRecipe
	err     error
}

func (mrf *MockRecipesFinder) FindAll() ([]*types.OpenInstallationRecipe, error) {
	if mrf.err != nil {
		return nil, mrf.err
	}
	return mrf.recipes, nil
}

type RecipeDetectorTestBuilder struct {
	processEvaluator *MockRecipeEvaluator
	scriptEvaluator  *MockRecipeEvaluator
	recipesFinder    *MockRecipesFinder
}

func NewRecipeDetectorTestBuilder() *RecipeDetectorTestBuilder {
	return &RecipeDetectorTestBuilder{
		processEvaluator: NewMockRecipeEvaluator(),
		scriptEvaluator:  NewMockRecipeEvaluator(),
		recipesFinder:    &MockRecipesFinder{},
	}
}

func (b *RecipeDetectorTestBuilder) WithRecipesFinderError(err error) *RecipeDetectorTestBuilder {
	b.recipesFinder.err = err
	return b
}

func (b *RecipeDetectorTestBuilder) WithProcessEvaluatorRecipeStatus(recipe *types.OpenInstallationRecipe, status execution.RecipeStatusType) *RecipeDetectorTestBuilder {
	b.recipesFinder.recipes = append(b.recipesFinder.recipes, recipe)
	b.processEvaluator.WithRecipeStatus(recipe, status)
	return b
}

func (b *RecipeDetectorTestBuilder) WithScriptEvaluatorRecipeStatus(recipe *types.OpenInstallationRecipe, status execution.RecipeStatusType) *RecipeDetectorTestBuilder {
	b.recipesFinder.recipes = append(b.recipesFinder.recipes, recipe)
	b.scriptEvaluator.WithRecipeStatus(recipe, status)
	return b
}

func (b *RecipeDetectorTestBuilder) Build() *RecipeDetector {
	return &RecipeDetector{
		context:          context.Background(),
		repo:             b.recipesFinder,
		processEvaluator: b.processEvaluator,
		scriptEvaluator:  b.scriptEvaluator,
	}
}
