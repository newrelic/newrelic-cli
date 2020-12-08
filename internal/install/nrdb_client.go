package install

import (
	"context"

	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

type nrdbClient interface {
	QueryWithContext(context.Context, int, nrdb.NRQL) (*nrdb.NRDBResultContainer, error)
}
