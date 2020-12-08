package install

import "context"

type mockRecipeValidator struct {
	result func(discoveryManifest, recipe) (bool, error)
}

func newMockRecipeValidator() *mockRecipeValidator {
	return &mockRecipeValidator{
		result: func(discoveryManifest, recipe) (bool, error) { return false, nil },
	}
}

func (m *mockRecipeValidator) validate(ctx context.Context, dm discoveryManifest, r recipe) (bool, error) {
	return m.result(dm, r)
}
