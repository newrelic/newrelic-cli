package install

import "context"

type mockRecipeExecutor struct {
	result bool
}

func newMockRecipeExecutor() *mockRecipeExecutor {
	return &mockRecipeExecutor{
		result: false,
	}
}

func (m *mockRecipeExecutor) execute(ctx context.Context, dm discoveryManifest, r recipe) error {
	return nil
}
