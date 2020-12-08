package install

import (
	"context"
	"errors"

	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

type mockNrdbClient struct {
	results  func() []nrdb.NRDBResult
	attempts int
	error    string
}

func newMockNrdbClient() *mockNrdbClient {
	return &mockNrdbClient{
		results: func() []nrdb.NRDBResult {
			return []nrdb.NRDBResult{}
		},
	}
}

func (c *mockNrdbClient) QueryWithContext(ctx context.Context, accountID int, nrql nrdb.NRQL) (*nrdb.NRDBResultContainer, error) {
	c.attempts++

	if c.error != "" {
		return nil, errors.New(c.error)
	}

	return &nrdb.NRDBResultContainer{
		Results: c.results(),
	}, nil
}

func (c *mockNrdbClient) ThrowError(message string) {
	c.error = message
}

func (c *mockNrdbClient) ReturnResultsAfterNAttempts(before []nrdb.NRDBResult, after []nrdb.NRDBResult, attempts int) {
	c.results = func() []nrdb.NRDBResult {
		if c.attempts < attempts {
			return before
		}

		return after
	}
}

func (c *mockNrdbClient) Attempts() int {
	return c.attempts
}
