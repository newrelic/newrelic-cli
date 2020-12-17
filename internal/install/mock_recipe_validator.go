package install

import "context"

type mockRecipeValidator struct {
	validateErr           error
	validateCallCount     int
	validateVal           bool
	validateEntityGUIDVal string
}

func newMockRecipeValidator() *mockRecipeValidator {
	return &mockRecipeValidator{}
}

func (m *mockRecipeValidator) validate(ctx context.Context, dm discoveryManifest, r recipe) (bool, string, error) {
	m.validateCallCount++
	return m.validateVal, m.validateEntityGUIDVal, m.validateErr
}
