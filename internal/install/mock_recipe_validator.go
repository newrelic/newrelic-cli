package install

import "context"

type mockRecipeValidator struct {
	result func(discoveryManifest, recipe) (bool, string, error)
}

func newMockRecipeValidator() *mockRecipeValidator {
	return &mockRecipeValidator{
		result: func(discoveryManifest, recipe) (bool, string, error) { return false, "", nil },
	}
}

func (m *mockRecipeValidator) validate(ctx context.Context, dm discoveryManifest, r recipe) (bool, string, error) {
	return m.result(dm, r)
}
