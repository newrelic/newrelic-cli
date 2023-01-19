package mocks

import (
	"context"
	"errors"

	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"
)

type MockNRDBClient struct {
	results  func() []nrdb.NRDBResult
	attempts int
	error    string
}

func NewMockNRDBClient() *MockNRDBClient {
	return &MockNRDBClient{
		results: func() []nrdb.NRDBResult {
			return []nrdb.NRDBResult{}
		},
	}
}

func (c *MockNRDBClient) QueryWithContext(ctx context.Context, accountID int, nrql nrdb.NRQL) (*nrdb.NRDBResultContainer, error) {
	c.attempts++

	if c.error != "" {
		return nil, errors.New(c.error)
	}

	return &nrdb.NRDBResultContainer{
		Results: c.results(),
	}, nil
}

func (c *MockNRDBClient) ThrowError(message string) {
	c.error = message
}

func (c *MockNRDBClient) ReturnResultsAfterNAttempts(before []nrdb.NRDBResult, after []nrdb.NRDBResult, attempts int) {
	c.results = func() []nrdb.NRDBResult {
		if c.attempts < attempts {
			return before
		}

		return after
	}
}

func (c *MockNRDBClient) Attempts() int {
	return c.attempts
}
