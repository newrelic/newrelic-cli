package configuration

import "os"

type envResolver interface {
	Getenv(key string) string
}

type osEnvResolver struct{}

func (r *osEnvResolver) Getenv(key string) string {
	return os.Getenv(key)
}

type mockEnvResolver struct {
	GetenvVal string
}

func (r *mockEnvResolver) Getenv(key string) string {
	return r.GetenvVal
}
