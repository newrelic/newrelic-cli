package configuration

import "os"

type envResolver interface {
	Getenv(key string) string
}

type osEnvResolver struct{}

func (r *osEnvResolver) Getenv(key string) string {
	return os.Getenv(key)
}

type MockEnvResolver struct {
	GetenvVal  string
	GetenvVals map[string]string
}

func NewMockEnvResolver() *MockEnvResolver {
	return &MockEnvResolver{
		GetenvVals: map[string]string{},
	}
}

func (r *MockEnvResolver) Getenv(key string) string {
	if r.GetenvVals != nil {
		return r.GetenvVals[key]
	}
	return r.GetenvVal
}
