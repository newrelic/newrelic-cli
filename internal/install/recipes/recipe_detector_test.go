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

func TestDetectBundleRecipe_NoDependency_Available(t *testing.T) {
	b := NewRecipeDetectorTestBuilder()
	br := NewRecipeBuilder().Name("r1").BuildBundleRecipe()
	b.WithProcessEvaluatorRecipeStatus(br.Recipe, execution.RecipeStatusTypes.AVAILABLE)
	detector := b.Build()

	detector.detectBundleRecipe(context.Background(), br)
	require.True(t, br.HasStatus(execution.RecipeStatusTypes.AVAILABLE))
}

func TestDetectBundleRecipe_SingleDependencyNotAvailable(t *testing.T) {
	b := NewRecipeDetectorTestBuilder()
	br := NewRecipeBuilder().Name("r1").DependencyBuilder(NewRecipeBuilder().Name("dep1")).BuildBundleRecipe()
	b.WithProcessEvaluatorRecipeStatus(br.Recipe, execution.RecipeStatusTypes.AVAILABLE)
	b.WithProcessEvaluatorRecipeStatus(br.Dependencies[0].Recipe, execution.RecipeStatusTypes.UNSUPPORTED)
	detector := b.Build()

	detector.detectBundleRecipe(context.Background(), br)
	require.False(t, br.HasStatus(execution.RecipeStatusTypes.AVAILABLE))
	require.True(t, br.Dependencies[0].HasStatus(execution.RecipeStatusTypes.UNSUPPORTED))
}

func TestDetectBundleRecipe_TwoDependencyAndDepUnsupported(t *testing.T) {
	b := NewRecipeDetectorTestBuilder()
	infra := NewRecipeBuilder().Name("infra")
	log := NewRecipeBuilder().Name("log").DependencyBuilder(infra)
	infraBundleRecipe := infra.BuildBundleRecipe()
	logBundleRecipe := log.BuildBundleRecipe()
	b.WithProcessEvaluatorRecipeStatus(logBundleRecipe.Recipe, execution.RecipeStatusTypes.AVAILABLE)
	b.WithProcessEvaluatorRecipeStatus(infraBundleRecipe.Recipe, execution.RecipeStatusTypes.UNSUPPORTED)
	detector := b.Build()

	detector.detectBundleRecipe(context.Background(), infraBundleRecipe)
	detector.detectBundleRecipe(context.Background(), logBundleRecipe)

	require.False(t, logBundleRecipe.HasStatus(execution.RecipeStatusTypes.AVAILABLE))
	require.True(t, infraBundleRecipe.HasStatus(execution.RecipeStatusTypes.UNSUPPORTED))
}

func TestDetectBundleRecipe_ParentShouldEvaluateDependency_Unsupported(t *testing.T) {
	b := NewRecipeDetectorTestBuilder()
	infra := NewRecipeBuilder().Name("infra")
	log := NewRecipeBuilder().Name("log").DependencyBuilder(infra)
	infraBundleRecipe := infra.BuildBundleRecipe()
	logBundleRecipe := log.BuildBundleRecipe()
	b.WithProcessEvaluatorRecipeStatus(logBundleRecipe.Recipe, execution.RecipeStatusTypes.AVAILABLE)
	b.WithProcessEvaluatorRecipeStatus(infraBundleRecipe.Recipe, execution.RecipeStatusTypes.UNSUPPORTED)
	detector := b.Build()

	detector.detectBundleRecipe(context.Background(), logBundleRecipe)

	require.False(t, logBundleRecipe.HasStatus(execution.RecipeStatusTypes.AVAILABLE))
	require.True(t, infraBundleRecipe.HasStatus(execution.RecipeStatusTypes.UNSUPPORTED))
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
		recipeEvaluated:  make(map[string][]*DetectedStatusType),
	}
}
