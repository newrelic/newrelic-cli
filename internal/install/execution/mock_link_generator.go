package execution

import "strings"

type MockPlatformLinkGenerator struct {
	GenerateExplorerLinkCallCount int
	GenerateEntityLinkCallCount   int
	GenerateLoggingLinkCallCount  int
	GenerateFleetLinkCallCount    int
	GenerateExplorerLinkVal       string
	GenerateEntityLinkVal         string
	GenerateLoggingLinkVal        string
	GenerateFleetLinkVal          string
}

func NewMockPlatformLinkGenerator() *MockPlatformLinkGenerator {
	return &MockPlatformLinkGenerator{}
}

func (g *MockPlatformLinkGenerator) GenerateExplorerLink(status InstallStatus) string {
	g.GenerateExplorerLinkCallCount++
	return g.GenerateExplorerLinkVal
}

func (g *MockPlatformLinkGenerator) GenerateEntityLink(entityGUID string) string {
	g.GenerateEntityLinkCallCount++
	return g.GenerateEntityLinkVal
}

func (g *MockPlatformLinkGenerator) GenerateLoggingLink(entityGUID string) string {
	g.GenerateLoggingLinkCallCount++
	return g.GenerateLoggingLinkVal
}

func (g *MockPlatformLinkGenerator) GenerateFleetLink(entityGUID string) string {
	g.GenerateFleetLinkCallCount++
	return g.GenerateFleetLinkVal
}

func (g *MockPlatformLinkGenerator) GenerateRedirectURL(status InstallStatus) string {
	if status.hasAnyRecipeStatus(RecipeStatusTypes.INSTALLED) {
		switch t := status.successLinkConfig.Type; {
		case strings.EqualFold(string(t), "explorer"):
			return g.GenerateExplorerLink(status)
		default:
			return g.GenerateEntityLink(status.HostEntityGUID())
		}
	}

	return ""
}
