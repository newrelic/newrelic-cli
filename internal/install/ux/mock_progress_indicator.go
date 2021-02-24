package ux

type MockProgressIndicator struct{}

func NewMockProgressIndicator() *MockProgressIndicator {
	return &MockProgressIndicator{}
}

func (s *MockProgressIndicator) Fail(string) {
}

func (s MockProgressIndicator) Success(string) {
}

func (s *MockProgressIndicator) Start(string) {
}

func (s MockProgressIndicator) Stop() {
}
