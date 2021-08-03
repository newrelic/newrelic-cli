package validation

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

const (
	DefaultMaxAttempts      = 60
	DefaultIntervalSeconds  = 5
	ValidationInProgressMsg = "Checking for data in New Relic (this may take a few minutes)..."
	ReachedMaxValidationMsg = "reached max validation attempts"
)

// PollingNRQLValidator polls NRDB to assert data is being reported for the given query.
type PollingNRQLValidator struct {
	MaxAttempts          int
	IntervalMilliSeconds int
	ProgressIndicator    ux.ProgressIndicator
	client               utils.NRDBClient
}

// NewPollingNRQLValidator returns a new instance of PollingNRQLValidator.
func NewPollingNRQLValidator(c utils.NRDBClient) *PollingNRQLValidator {
	v := PollingNRQLValidator{
		client:               c,
		MaxAttempts:          DefaultMaxAttempts,
		IntervalMilliSeconds: DefaultIntervalSeconds * 1000,
		ProgressIndicator:    ux.NewSpinner(),
	}

	return &v
}

// Validate polls NRDB to assert data is being reported for the given query.
func (m *PollingNRQLValidator) Validate(ctx context.Context, query string) (string, error) {
	ticker := time.NewTicker(time.Duration(m.IntervalMilliSeconds) * time.Millisecond)
	defer ticker.Stop()

	m.ProgressIndicator.Start(ValidationInProgressMsg)
	defer m.ProgressIndicator.Stop()

	entityGUID, err := m.tryValidate(ctx, query)
	if err != nil {
		m.ProgressIndicator.Fail("")
		if strings.Contains(err.Error(), "context canceled") {
			return "", err
		}
		return "", fmt.Errorf("%s: %s", ReachedMaxValidationMsg, err)
	}
	m.ProgressIndicator.Success("")
	return entityGUID, nil
}

func (m *PollingNRQLValidator) tryValidate(ctx context.Context, query string) (string, error) {
	guid := ""
	validatorFunc := func() error {
		results, err := m.executeQuery(ctx, query)
		if err != nil {
			return err
		}

		if len(results) == 0 {
			return errors.New("no results returned")
		}

		// The query is assumed to use a count aggregate function
		count := results[0]["count"].(float64)

		if count > 0 {
			// Try and parse an entity GUID from the results.  The query is assumed to
			// optionally use a facet over entityGuid.  The standard case seems to be
			// that all entities contain a facet of "entityGuid", and so if we find it
			// here, we return it.
			if entityGUID, ok := results[0]["entityGuid"]; ok {
				guid = entityGUID.(string)
				return nil
			}

			// In the logs integration, the facet doesn't contain "entityGuid", but
			// does contain, "entity.guid", so here we check for that also.
			if entityGUID, ok := results[0]["entity.guids"]; ok {
				guid = entityGUID.(string)
				return nil
			}

			// entity guid is optional, no error returned
			return nil
		}

		return errors.New("no count found in results")
	}

	retry := utils.NewRetry(m.MaxAttempts, m.IntervalMilliSeconds, validatorFunc)
	if err := retry.ExecWithRetries(ctx); err != nil {
		return "", err
	}

	return guid, nil
}

func (m *PollingNRQLValidator) executeQuery(ctx context.Context, query string) ([]nrdb.NRDBResult, error) {
	accountID := configAPI.RequireActiveProfileAccountID()

	nrql := nrdb.NRQL(query)

	result, err := m.client.QueryWithContext(ctx, accountID, nrql)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}
