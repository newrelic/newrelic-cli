package install

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

const (
	defaultMaxAttempts = 20
	defaultInterval    = 5 * time.Second
)

type pollingRecipeValidator struct {
	maxAttempts int
	interval    time.Duration
	client      nrdbClient
}

func newPollingRecipeValidator(c nrdbClient) *pollingRecipeValidator {
	v := pollingRecipeValidator{
		maxAttempts: defaultMaxAttempts,
		interval:    defaultInterval,
		client:      c,
	}

	return &v
}

func (m *pollingRecipeValidator) validate(ctx context.Context, r recipe) (bool, error) {
	count := 0
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		if count == m.maxAttempts {
			return false, nil
		}

		log.Debugf("Validation attempt #%d...", count+1)
		ok, err := m.tryValidate(ctx, r)
		if err != nil {
			return false, err
		}

		count++

		if ok {
			return true, nil
		}

		select {
		case <-ticker.C:
			continue

		case <-ctx.Done():
			return false, nil
		}
	}
}

func (m *pollingRecipeValidator) tryValidate(ctx context.Context, r recipe) (bool, error) {
	results, err := m.executeQuery(ctx, r.ValidationNRQL)
	if err != nil {
		return false, err
	}

	if len(results) > 0 {
		return true, nil
	}

	return false, nil
}

func (m *pollingRecipeValidator) executeQuery(ctx context.Context, query string) ([]nrdb.NrdbResult, error) {
	accountID := credentials.DefaultProfile().AccountID
	nrql := nrdb.Nrql(query)

	result, err := m.client.QueryWithContext(ctx, accountID, nrql)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}
