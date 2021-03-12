package execution

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
