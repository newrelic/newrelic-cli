package diagnose

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/utils/validation"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

const (
	validationEventType = "NrIntegrationError"
)

type ConfigValidator struct {
	client *newrelic.NewRelic
	*validation.PollingNRQLValidator
	discovery.Discoverer
}

type ValidationTracerEvent struct {
	EventType string `json:"eventType"`
	Hostname  string `json:"hostname"`
}

func NewConfigValidator(client *newrelic.NewRelic) *ConfigValidator {
	pf := discovery.NewNoOpProcessFilterer()

	return &ConfigValidator{
		client:               client,
		PollingNRQLValidator: validation.NewPollingNRQLValidator(&client.Nrdb),
		Discoverer:           discovery.NewPSUtilDiscoverer(pf),
	}
}

func (c *ConfigValidator) ValidateConfig(ctx context.Context) error {
	defaultProfile := credentials.DefaultProfile()
	manifest, err := c.Discover(ctx)
	if err != nil {
		return err
	}

	evt := ValidationTracerEvent{
		EventType: validationEventType,
		Hostname:  manifest.Hostname,
	}

	log.Printf("Sending tracer event to New Relic.")

	err = c.client.Events.CreateEvent(defaultProfile.AccountID, evt)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
	FROM %s
	SELECT count(*)
	WHERE hostname LIKE '%s%%'
	SINCE 10 MINUTES AGO
	`, validationEventType, manifest.Hostname)

	_, err = c.Validate(ctx, query)

	return err
}
