package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockRecipeRecommender struct {
	RecommendVal       []types.OpenInstallationRecipe
	RecommendCallCount int
	RecommendErr       error
}

func NewMockRecipeRecommender() *MockRecipeRecommender {
	return &MockRecipeRecommender{}
}

func (r *MockRecipeRecommender) Recommend(ctx context.Context, m *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error) {
	r.RecommendCallCount++

	if r.RecommendErr != nil {
		return nil, r.RecommendErr
	}

	return r.RecommendVal, nil
}
