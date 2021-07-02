package execution

import "github.com/newrelic/newrelic-client-go/pkg/installevents"

type MockInstalleventsClient struct {
	CreateInstallEventVal       *installevents.InstallEvent
	CreateInstallEventErr       error
	CreateInstallEventCallCount int
}

func NewMockInstalleventsClient() *MockInstalleventsClient {
	return &MockInstalleventsClient{}
}

func (c *MockInstalleventsClient) CreateInstallEvent(installevents.InstallStatus) (*installevents.InstallEvent, error) {
	c.CreateInstallEventCallCount++
	return c.CreateInstallEventVal, c.CreateInstallEventErr
}
