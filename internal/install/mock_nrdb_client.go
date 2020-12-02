package install

import (
	"context"
	"errors"
	"fmt"

	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

type mockNrdbClient struct {
	results  func() []nrdb.NrdbResult
	attempts int
	error    string
}

func newMockNrdbClient() *mockNrdbClient {
	return &mockNrdbClient{
		results: func() []nrdb.NrdbResult {
			return []nrdb.NrdbResult{}
		},
	}
}

func (c *mockNrdbClient) QueryWithContext(ctx context.Context, accountID int, nrql nrdb.Nrql) (*nrdb.NrdbResultContainer, error) {
	c.attempts++

	if c.error != "" {
		return nil, errors.New(c.error)
	}

	return &nrdb.NrdbResultContainer{
		Results: c.results(),
	}, nil
}

func (c *mockNrdbClient) ThrowError(message string) {
	c.error = message
}

func (c *mockNrdbClient) ReturnResultsAfterNAttempts(results []nrdb.NrdbResult, attempts int) {
	c.results = func() []nrdb.NrdbResult {
		fmt.Printf("c.attempts: %d\n", c.attempts)
		fmt.Printf("attempts: %d\n", attempts)
		if c.attempts < attempts {
			return []nrdb.NrdbResult{}
		}

		fmt.Printf("returning desired results\n")
		return results
	}
}

func (c *mockNrdbClient) Attempts() int {
	return c.attempts
}
