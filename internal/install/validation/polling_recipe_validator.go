package validation

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

type contextKey int

const (
	defaultMaxAttempts            = 20
	defaultInterval               = 5 * time.Second
	TestIdentifierKey  contextKey = iota
)

// PollingRecipeValidator is an implementation of the RecipeValidator interface
// that polls NRDB to assert data is being reported for the given recipe.
type PollingRecipeValidator struct {
	maxAttempts int
	interval    time.Duration
	client      nrdbClient
}

// NewPollingRecipeValidator returns a new instance of PollingRecipeValidator.
func NewPollingRecipeValidator(c nrdbClient) *PollingRecipeValidator {
	v := PollingRecipeValidator{
		maxAttempts: defaultMaxAttempts,
		interval:    defaultInterval,
		client:      c,
	}

	return &v
}

// Validate polls NRDB to assert data is being reported for the given recipe.
func (m *PollingRecipeValidator) Validate(ctx context.Context, dm types.DiscoveryManifest, r types.Recipe) (string, error) {
	return m.waitForData(ctx, dm, r)
}

func (m *PollingRecipeValidator) waitForData(ctx context.Context, dm types.DiscoveryManifest, r types.Recipe) (string, error) {
	count := 0
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	if isTest := ctx.Value(TestIdentifierKey); isTest == nil {
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Start()
		defer s.Stop()
	}

	for {
		if count == m.maxAttempts {
			return "", fmt.Errorf("reached max validation attempts")
		}

		log.Debugf("Validation attempt #%d...", count+1)
		ok, entityGUID, err := m.tryValidate(ctx, dm, r)
		if err != nil {
			return "", err
		}

		count++

		if ok {
			return entityGUID, nil
		}

		select {
		case <-ticker.C:
			continue

		case <-ctx.Done():
			return "", fmt.Errorf("validation cancelled")
		}
	}
}

func (m *PollingRecipeValidator) tryValidate(ctx context.Context, dm types.DiscoveryManifest, r types.Recipe) (bool, string, error) {
	query, err := substituteHostname(dm, r)
	if err != nil {
		return false, "", err
	}

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
		// Try and parse an entity GUID from the results
		// The query is assumed to optionally use a facet over entityGuid
		if entityGUID, ok := results[0]["entityGuid"]; ok {
			return true, entityGUID.(string), nil
		}

		return true, "", nil
	}

	return false, "", nil
}

func substituteHostname(dm types.DiscoveryManifest, r types.Recipe) (string, error) {
	tmpl, err := template.New("validationNRQL").Parse(r.ValidationNRQL)
	if err != nil {
		panic(err)
	}

	v := struct {
		HOSTNAME string
	}{
		HOSTNAME: dm.Hostname,
	}

	var tpl bytes.Buffer
	if err = tmpl.Execute(&tpl, v); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func (m *PollingRecipeValidator) executeQuery(ctx context.Context, query string) ([]nrdb.NRDBResult, error) {
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
