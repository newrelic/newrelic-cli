package execution

import "github.com/newrelic/newrelic-client-go/pkg/installevents"

type MockInstalleventsClient struct {
	CreateInstallEventVal        *installevents.InstallationRecipeEvent
	CreateInstallStatusVal       *installevents.InstallationInstallStatus
	CreateInstallEventErr        error
	CreateInstallEventCallCount  int
	CreateInstallStatusErr       error
	CreateInstallStatusCallCount int

	// Mock recipe event calls
	CreateRecipeEventCallCount int
}

func NewMockInstallEventsClient() *MockInstalleventsClient {
	return &MockInstalleventsClient{}
}

func (c *MockInstalleventsClient) InstallationCreateRecipeEvent(int, installevents.InstallationRecipeStatus) (*installevents.InstallationRecipeEvent, error) {
	c.CreateRecipeEventCallCount++
	return c.CreateInstallEventVal, c.CreateInstallEventErr
}

func (c *MockInstalleventsClient) InstallationCreateInstallStatus(int, installevents.InstallationInstallStatusInput) (*installevents.InstallationInstallStatus, error) {
	c.CreateInstallStatusCallCount++
	return c.CreateInstallStatusVal, c.CreateInstallStatusErr
}
