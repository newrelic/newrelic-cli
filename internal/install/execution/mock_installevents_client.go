package execution

import "github.com/newrelic/newrelic-client-go/pkg/installevents"

type MockInstalleventsClient struct {
	CreateInstallEventVal          *installevents.InstallEvent
	CreateInstallMetadataVal       *installevents.InstallMetadata
	CreateInstallEventErr          error
	CreateInstallMetadataErr       error
	CreateInstallEventCallCount    int
	CreateInstallMetadataCallCount int
}

func NewMockInstalleventsClient() *MockInstalleventsClient {
	return &MockInstalleventsClient{}
}

func (c *MockInstalleventsClient) CreateInstallEvent(installevents.InstallStatus) (*installevents.InstallEvent, error) {
	c.CreateInstallEventCallCount++
	return c.CreateInstallEventVal, c.CreateInstallEventErr
}

func (c *MockInstalleventsClient) CreateInstallMetadata(installevents.InputInstallMetadata) (*installevents.InstallMetadata, error) {
	c.CreateInstallEventCallCount++
	return c.CreateInstallMetadataVal, c.CreateInstallMetadataErr
}
