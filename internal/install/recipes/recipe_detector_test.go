package recipes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
)

func TestRecipeDetectorShouldFailBecauseOfProcessEvaluation(t *testing.T) {
	recipe := NewRecipeBuilder().Build()
	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorStatus(execution.RecipeStatusTypes.NULL)
	b.WithScriptEvaluatorStatus(execution.RecipeStatusTypes.AVAILABLE)
	detector := b.Build()

	actual, _ := detector.detectRecipe(context.Background(), recipe)
	require.Equal(t, execution.RecipeStatusTypes.NULL, actual)
}

func TestRecipeDetectorShouldBeAvailableWhenRecipeScriptDetectionIsMissingScript(t *testing.T) {
	recipe := NewRecipeBuilder().Build()
	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorStatus(execution.RecipeStatusTypes.AVAILABLE)
	b.WithScriptEvaluatorStatus(execution.RecipeStatusTypes.DETECTED)
	detector := b.Build()

	actual, _ := detector.detectRecipe(context.Background(), recipe)
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, actual)
}

func TestRecipeDetectorShouldFailWhenScriptFails(t *testing.T) {
	recipe := NewRecipeBuilder().WithPreInstallScript("pre-install script mock").Build()

	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorStatus(execution.RecipeStatusTypes.AVAILABLE)
	b.WithScriptEvaluatorStatus(execution.RecipeStatusTypes.NULL)
	detector := b.Build()

	actual, _ := detector.detectRecipe(context.Background(), recipe)
	require.Equal(t, execution.RecipeStatusTypes.NULL, actual)
}

func TestRecipeDetectorShouldDetectBecauseOfScriptEvaluation(t *testing.T) {
	recipe := NewRecipeBuilder().WithPreInstallScript("pre-install script mock").Build()

	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorStatus(execution.RecipeStatusTypes.AVAILABLE)
	b.WithScriptEvaluatorStatus(execution.RecipeStatusTypes.DETECTED)
	detector := b.Build()

	actual, _ := detector.detectRecipe(context.Background(), recipe)
	require.Equal(t, execution.RecipeStatusTypes.DETECTED, actual)
}

func TestRecipeDetectorShouldBeAvailableBecauseOfScriptEvaluation(t *testing.T) {
	recipe := NewRecipeBuilder().WithPreInstallScript("pre-install script mock").Build()

	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorStatus(execution.RecipeStatusTypes.AVAILABLE)
	b.WithScriptEvaluatorStatus(execution.RecipeStatusTypes.AVAILABLE)
	detector := b.Build()

	actual, _ := detector.detectRecipe(context.Background(), recipe)
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, actual)
}

type RecipeDetectorTestBuilder struct {
	processEvaluator *MockRecipeEvaluator
	scriptEvaluator  *MockRecipeEvaluator
}

func NewRecipeDetectorTestBuilder() *RecipeDetectorTestBuilder {
	return &RecipeDetectorTestBuilder{
		processEvaluator: &MockRecipeEvaluator{},
		scriptEvaluator:  &MockRecipeEvaluator{},
	}
}

func (b *RecipeDetectorTestBuilder) WithProcessEvaluatorStatus(status execution.RecipeStatusType) *RecipeDetectorTestBuilder {
	b.processEvaluator.status = status
	return b
}

func (b *RecipeDetectorTestBuilder) WithScriptEvaluatorStatus(status execution.RecipeStatusType) *RecipeDetectorTestBuilder {
	b.scriptEvaluator.status = status
	return b
}

func (b *RecipeDetectorTestBuilder) Build() *RecipeDetector {
	return newRecipeDetector(b.processEvaluator, b.scriptEvaluator)
}

func newRecipeDetector(processEvaluator DetectionStatusProvider, scriptEvaluator DetectionStatusProvider) *RecipeDetector {
	return &RecipeDetector{
		processEvaluator: processEvaluator,
		scriptEvaluator:  scriptEvaluator,
		recipeEvaluated:  make(map[string]bool),
	}
}
