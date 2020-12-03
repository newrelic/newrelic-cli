package install

import "context"

// nolint:unused
type mockRecipeValidator struct {
	result func(recipeFile) (bool, error)
}

func (m *mockRecipeValidator) validate(ctx context.Context, r recipeFile) (bool, error) {
	return m.result(r)
}
