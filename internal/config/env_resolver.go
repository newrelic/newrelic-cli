package config

import "os"

type envResolver interface {
	Getenv(key string) string
}

type OSEnvResolver struct{}

func (r *OSEnvResolver) Getenv(key string) string {
	return os.Getenv(key)
}

func NewOSEnvResolver() *OSEnvResolver {
	return &OSEnvResolver{}
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

func (r *MockEnvResolver) Setenv(key string, val string) {
	r.GetenvVals[key] = val
}

func (r *MockEnvResolver) Getenv(key string) string {
	if r.GetenvVals != nil {
		return r.GetenvVals[key]
	}
	return r.GetenvVal
}
