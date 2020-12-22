package ux

type MockProgressIndicator struct{}

func NewMockProgressIndicator() *MockProgressIndicator {
	return &MockProgressIndicator{}
}

func (s *MockProgressIndicator) Fail() {
}

func (s MockProgressIndicator) Success() {
}

func (s *MockProgressIndicator) Start(string) {
}

func (s MockProgressIndicator) Stop() {
}
