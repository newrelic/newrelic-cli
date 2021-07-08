package execution

import "github.com/newrelic/newrelic-client-go/pkg/installevents"

type MockInstalleventsClient struct {
	CreateInstallEventVal       *installevents.RecipeEvent
	CreateInstallEventErr       error
	CreateInstallEventCallCount int
}

func NewMockInstalleventsClient() *MockInstalleventsClient {
	return &MockInstalleventsClient{}
}

func (c *MockInstalleventsClient) CreateRecipeEvent(int, installevents.RecipeStatus) (*installevents.RecipeEvent, error) {
	c.CreateInstallEventCallCount++
	return c.CreateInstallEventVal, c.CreateInstallEventErr
}
