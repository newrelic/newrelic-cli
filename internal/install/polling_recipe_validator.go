package install

import (
	"context"
	"strings"
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

func (m *pollingRecipeValidator) validate(ctx context.Context, dm discoveryManifest, r recipe) (bool, error) {
	count := 0
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		if count == m.maxAttempts {
			return false, nil
		}

		log.Debugf("Validation attempt #%d...", count+1)
		ok, err := m.tryValidate(ctx, dm, r)
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

func (m *pollingRecipeValidator) tryValidate(ctx context.Context, dm discoveryManifest, r recipe) (bool, error) {
	query := substituteHostname(dm, r)
	results, err := m.executeQuery(ctx, query)
	if err != nil {
		return false, err
	}

	if len(results) == 0 {
		return false, nil
	}

	// The query is assumed to use a count aggregate function
	count := results[0]["count"].(float64)

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// TODO: replace with go templates
func substituteHostname(dm discoveryManifest, r recipe) string {
	return strings.Replace(r.ValidationNRQL, "HOSTNAME", dm.Hostname, -1)
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
