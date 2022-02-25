package install

import (
	"github.com/stretchr/testify/mock"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
)

func (sr *mockStatusReporter) ReportStatus(status execution.RecipeStatusType, event execution.RecipeStatusEvent) {
	sr.Called(status, event)
}

type mockStatusReporter struct {
	mock.Mock
}
