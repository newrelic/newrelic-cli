package validation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

const (
	defaultMaxAttempts = 60
	defaultInterval    = 5 * time.Second
)

// PollingNRQLValidator polls NRDB to assert data is being reported for the given query.
type PollingNRQLValidator struct {
	MaxAttempts       int
	Interval          time.Duration
	ProgressIndicator ux.ProgressIndicator
	client            utils.NRDBClient
}

// NewPollingNRQLValidator returns a new instance of PollingNRQLValidator.
func NewPollingNRQLValidator(c utils.NRDBClient) *PollingNRQLValidator {
	v := PollingNRQLValidator{
		client:            c,
		MaxAttempts:       defaultMaxAttempts,
		Interval:          defaultInterval,
		ProgressIndicator: ux.NewSpinner(),
	}

	return &v
}

// Validate polls NRDB to assert data is being reported for the given query.
func (m *PollingNRQLValidator) Validate(ctx context.Context, query string) (string, error) {
	return m.waitForData(ctx, query)
}

func (m *PollingNRQLValidator) waitForData(ctx context.Context, query string) (string, error) {
	count := 0
	ticker := time.NewTicker(m.Interval)
	defer ticker.Stop()

	progressMsg := "Checking for data in New Relic (this may take a few minutes)..."
	m.ProgressIndicator.Start(progressMsg)
	defer m.ProgressIndicator.Stop()

	for {
		if count == m.MaxAttempts {
			m.ProgressIndicator.Fail("")
			return "", fmt.Errorf("reached max validation attempts")
		}

		ok, entityGUID, err := m.tryValidate(ctx, query)
		if err != nil {
			m.ProgressIndicator.Fail("")
			return "", err
		}

		count++

		if ok {
			m.ProgressIndicator.Success("")
			return entityGUID, nil
		}

		select {
		case <-ticker.C:
			continue

		case <-ctx.Done():
			m.ProgressIndicator.Fail("")
			return "", fmt.Errorf("validation cancelled")
		}
	}
}

func (m *PollingNRQLValidator) tryValidate(ctx context.Context, query string) (bool, string, error) {
	results, err := m.executeQuery(ctx, query)
	if err != nil {
		return false, "", err
	}

	if len(results) == 0 {
		return false, "", nil
	}

	// The query is assumed to use a count aggregate function
	count := results[0]["count"].(float64)

	if count > 0 {
		// Try and parse an entity GUID from the results.  The query is assumed to
		// optionally use a facet over entityGuid.  The standard case seems to be
		// that all entities contain a facet of "entityGuid", and so if we find it
		// here, we return it.
		if entityGUID, ok := results[0]["entityGuid"]; ok {
			return true, entityGUID.(string), nil
		}

		// In the logs integration, the facet doesn't contain "entityGuid", but
		// does contain, "entity.guid", so here we check for that also.
		if entityGUID, ok := results[0]["entity.guids"]; ok {
			return true, entityGUID.(string), nil
		}

		return true, "", nil
	}

	return false, "", nil
}

func (m *PollingNRQLValidator) executeQuery(ctx context.Context, query string) ([]nrdb.NRDBResult, error) {
	profile := credentials.DefaultProfile()
	if profile == nil || profile.AccountID == 0 {
		return nil, errors.New("no account ID found in default profile")
	}

	nrql := nrdb.NRQL(query)

	result, err := m.client.QueryWithContext(ctx, profile.AccountID, nrql)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}
