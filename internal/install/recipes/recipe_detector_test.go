package recipes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestRecipeDetectorShouldFailBecauseOfProcessEvaluation(t *testing.T) {
	recipe := NewRecipeBuilder().Build()
	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorStatus(execution.RecipeStatusTypes.NULL)
	b.WithScriptEvaluatorStatus(execution.RecipeStatusTypes.AVAILABLE)
	detector := b.Build()

	actual := detector.detectRecipe(context.Background(), recipe)
	require.Equal(t, execution.RecipeStatusTypes.NULL, actual.Status)
}

func TestRecipeDetectorShouldBeAvailableWhenRecipeScriptDetectionIsMissingScript(t *testing.T) {
	recipe := NewRecipeBuilder().Build()
	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorStatus(execution.RecipeStatusTypes.AVAILABLE)
	b.WithScriptEvaluatorStatus(execution.RecipeStatusTypes.DETECTED)
	detector := b.Build()

	actual := detector.detectRecipe(context.Background(), recipe)
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, actual.Status)
}

func TestRecipeDetectorShouldFailWhenScriptFails(t *testing.T) {
	recipe := NewRecipeBuilder().WithPreInstallScript("pre-install script mock").Build()

	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorStatus(execution.RecipeStatusTypes.AVAILABLE)
	b.WithScriptEvaluatorStatus(execution.RecipeStatusTypes.NULL)
	detector := b.Build()

	actual := detector.detectRecipe(context.Background(), recipe)
	require.Equal(t, execution.RecipeStatusTypes.NULL, actual.Status)
}

func TestRecipeDetectorShouldDetectBecauseOfScriptEvaluation(t *testing.T) {
	recipe := NewRecipeBuilder().WithPreInstallScript("pre-install script mock").Build()

	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorStatus(execution.RecipeStatusTypes.AVAILABLE)
	b.WithScriptEvaluatorStatus(execution.RecipeStatusTypes.DETECTED)
	detector := b.Build()

	actual := detector.detectRecipe(context.Background(), recipe)
	require.Equal(t, execution.RecipeStatusTypes.DETECTED, actual.Status)
}

func TestRecipeDetectorShouldBeAvailableBecauseOfScriptEvaluation(t *testing.T) {
	recipe := NewRecipeBuilder().WithPreInstallScript("pre-install script mock").Build()

	b := NewRecipeDetectorTestBuilder()
	b.WithProcessEvaluatorStatus(execution.RecipeStatusTypes.AVAILABLE)
	b.WithScriptEvaluatorStatus(execution.RecipeStatusTypes.AVAILABLE)
	detector := b.Build()

	actual := detector.detectRecipe(context.Background(), recipe)
	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, actual.Status)
}

type RecipeDetectorTestBuilder struct {
	processEvaluator *MockRecipeEvaluator
	scriptEvaluator  *MockRecipeEvaluator
}

func NewRecipeDetectorTestBuilder() *RecipeDetectorTestBuilder {
	return &RecipeDetectorTestBuilder{
		processEvaluator: NewMockRecipeEvaluator(),
		scriptEvaluator:  NewMockRecipeEvaluator(),
	}
}

func (b *RecipeDetectorTestBuilder) WithProcessEvaluatorStatus(status execution.RecipeStatusType) *RecipeDetectorTestBuilder {
	b.processEvaluator.status = status
	return b
}

func (b *RecipeDetectorTestBuilder) WithProcessEvaluatorRecipeStatus(recipe *types.OpenInstallationRecipe, status execution.RecipeStatusType) *RecipeDetectorTestBuilder {
	b.processEvaluator.WithRecipeStatus(recipe, status)
	return b
}

func (b *RecipeDetectorTestBuilder) WithScriptEvaluatorStatus(status execution.RecipeStatusType) *RecipeDetectorTestBuilder {
	b.scriptEvaluator.status = status
	return b
}

func (b *RecipeDetectorTestBuilder) WithScriptEvaluatorRecipeStatus(recipe *types.OpenInstallationRecipe, status execution.RecipeStatusType) *RecipeDetectorTestBuilder {
	b.scriptEvaluator.WithRecipeStatus(recipe, status)
	return b
}

func (b *RecipeDetectorTestBuilder) Build() *RecipeDetector {
	return newRecipeDetector(b.processEvaluator, b.scriptEvaluator)
}

func newRecipeDetector(processEvaluator DetectionStatusProvider, scriptEvaluator DetectionStatusProvider) *RecipeDetector {
	return &RecipeDetector{
		processEvaluator: processEvaluator,
		scriptEvaluator:  scriptEvaluator,
	}
}
