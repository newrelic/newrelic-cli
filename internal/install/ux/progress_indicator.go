package ux

import "github.com/newrelic/newrelic-cli/internal/install/types"

type ProgressIndicator interface {
	Fail(types.Recipe)
	Success(types.Recipe)
	Start(types.Recipe)
	Stop()
}
