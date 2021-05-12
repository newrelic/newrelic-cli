package execution

import "strings"

type MockSuccessLinkGenerator struct {
	GenerateExplorerLinkCallCount int
	GenerateEntityLinkCallCount   int
	GenerateExplorerLinkVal       string
	GenerateEntityLinkVal         string
}

func NewMockSuccessLinkGenerator() *MockSuccessLinkGenerator {
	return &MockSuccessLinkGenerator{}
}

func (g *MockSuccessLinkGenerator) GenerateExplorerLink(filter string) string {
	g.GenerateExplorerLinkCallCount++
	return g.GenerateExplorerLinkVal
}

func (g *MockSuccessLinkGenerator) GenerateEntityLink(entityGUID string) string {
	g.GenerateEntityLinkCallCount++
	return g.GenerateEntityLinkVal
}

func (g *MockSuccessLinkGenerator) GenerateRedirectURL(status InstallStatus) string {
	if status.hasAnyRecipeStatus(RecipeStatusTypes.INSTALLED) {
		switch t := status.successLinkConfig.Type; {
		case strings.EqualFold(string(t), "explorer"):
			return g.GenerateExplorerLink(status.successLinkConfig.Filter)
		default:
			return g.GenerateEntityLink(status.HostEntityGUID())
		}
	}

	return ""
}
