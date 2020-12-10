package install

import "context"

type mockFileFilterer struct{}

func newMockFileFilterer() *mockFileFilterer {
	return &mockFileFilterer{}
}

func (m *mockFileFilterer) filter(ctx context.Context, recipes []recipe) ([]logMatch, error) {
	return nil, nil
}
