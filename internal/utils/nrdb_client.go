package utils

import (
	"context"

	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"
)

type NRDBClient interface {
	QueryWithContext(context.Context, int, nrdb.NRQL) (*nrdb.NRDBResultContainer, error)
}
