package execution

import "github.com/newrelic/newrelic-client-go/pkg/installevents"

type MockInstalleventsClient struct {
	CreateInstallEventVal       *installevents.InstallationRecipeEvent
	CreateInstallEventErr       error
	CreateInstallEventCallCount int
}

func NewMockInstallEventsClient() *MockInstalleventsClient {
	return &MockInstalleventsClient{}
}

func (c *MockInstalleventsClient) InstallationCreateRecipeEvent(int, installevents.InstallationRecipeStatus) (*installevents.InstallationRecipeEvent, error) {
	c.CreateInstallEventCallCount++
	return c.CreateInstallEventVal, c.CreateInstallEventErr
}
