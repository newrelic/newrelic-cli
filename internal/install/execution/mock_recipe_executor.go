package execution

import (
	"context"
	"encoding/json"
	"log"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockRecipeExecutor struct {
	ExecuteErr   error
	OutputParser *OutputParser
}

func NewMockRecipeExecutor() *MockRecipeExecutor {
	return &MockRecipeExecutor{
		OutputParser: NewOutputParser(map[string]interface{}{}),
	}
}

func (m *MockRecipeExecutor) Execute(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return m.ExecuteErr
}

func (m *MockRecipeExecutor) ExecutePreInstall(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return m.ExecuteErr
}

func (m *MockRecipeExecutor) GetOutput() *OutputParser {
	return m.OutputParser
}

func (m *MockRecipeExecutor) SetOutput(value string) {
	if value != "" {
		var values map[string]interface{}
		if err := json.Unmarshal([]byte(value), &values); err == nil {
			m.OutputParser = NewOutputParser(values)
			return
		}
		log.Fatalf("couldn't unmarshal json for mock with %s", value)
	}
}
