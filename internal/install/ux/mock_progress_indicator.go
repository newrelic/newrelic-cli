package ux

type MockProgressIndicator struct {
	Msg            string
	spinnerEnabled bool
}

func NewMockProgressIndicator() *MockProgressIndicator {
	return &MockProgressIndicator{}
}

func (s *MockProgressIndicator) Fail(string) {
	s.Msg += "Fail"
}

func (s *MockProgressIndicator) Success(string) {
	s.Msg += "Success"
}

func (s *MockProgressIndicator) Start(string) {
	s.Msg += "Start"
}

func (s *MockProgressIndicator) Stop() {
	s.Msg += "Stop"
}

func (s *MockProgressIndicator) Canceled(m string) {
	s.Msg += "Canceled"
}

func (s *MockProgressIndicator) ShowSpinner(b bool) {
	s.spinnerEnabled = b
}
