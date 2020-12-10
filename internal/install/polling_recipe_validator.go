package install

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"time"

	"github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

type contextKey int

const (
	defaultMaxAttempts            = 20
	defaultInterval               = 5 * time.Second
	TestIdentifierKey  contextKey = iota
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

	if isTest := ctx.Value(TestIdentifierKey); isTest == nil {
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Start()
		defer s.Stop()
	}

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
	query, err := substituteHostname(dm, r)
	if err != nil {
		return false, err
	}

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

func substituteHostname(dm discoveryManifest, r recipe) (string, error) {
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

func (m *pollingRecipeValidator) executeQuery(ctx context.Context, query string) ([]nrdb.NRDBResult, error) {
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
