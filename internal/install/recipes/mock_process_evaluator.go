package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockProcessEvaluator struct {
	processes []types.GenericProcess
}

func NewMockProcessEvaluator() *MockProcessEvaluator {
	return &MockProcessEvaluator{
		processes: []types.GenericProcess{},
	}
}

func (pe *MockProcessEvaluator) WithProcesses(processes []types.GenericProcess) {
	pe.processes = processes
}

func (pe *MockProcessEvaluator) GetOrLoadProcesses(ctx context.Context) []types.GenericProcess {
	return pe.processes
}

func (pe *MockProcessEvaluator) DetectionStatus(ctx context.Context, r *types.OpenInstallationRecipe) execution.RecipeStatusType {
	return execution.RecipeStatusTypes.AVAILABLE
}
