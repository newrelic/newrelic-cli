package diagnose

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/shirou/gopsutil/host"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/utils/validation"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

const (
	validationEventType = "NrIntegrationError"
)

type ConfigValidator struct {
	client *newrelic.NewRelic
	*validation.PollingNRQLValidator
}

type ValidationTracerEvent struct {
	EventType string `json:"eventType"`
	Hostname  string `json:"hostname"`
	Purpose   string `json:"purpose"`
	GUID      string `json:"guid"`
}

func NewConfigValidator(client *newrelic.NewRelic) *ConfigValidator {
	v := validation.NewPollingNRQLValidator(&client.Nrdb)
	v.MaxAttempts = 20

	return &ConfigValidator{
		client:               client,
		PollingNRQLValidator: v,
	}
}

func (c *ConfigValidator) ValidateConfig(ctx context.Context) error {
	defaultProfile := credentials.DefaultProfile()

	i, err := host.InfoWithContext(ctx)
	if err != nil {
		log.Error(err)
		return ErrDiscovery
	}

	evt := ValidationTracerEvent{
		EventType: validationEventType,
		Hostname:  i.Hostname,
		Purpose:   "New Relic CLI configuration validation",
		GUID:      uuid.NewString(),
	}

	log.Printf("Sending tracer event to New Relic.")

	if err = c.client.Events.CreateEvent(defaultProfile.AccountID, evt); err != nil {
		log.Error(reflect.TypeOf(err))
		log.Error(err)
		return ErrPostEvent
	}

	query := fmt.Sprintf(`
	FROM %s
	SELECT count(*)
	WHERE hostname LIKE '%s%%'
	AND guid = '%s'
	SINCE 10 MINUTES AGO
	`, evt.EventType, evt.Hostname, evt.GUID)

	if _, err = c.Validate(ctx, query); err != nil {
		log.Error(err)
		err = ErrValidation
	}

	return err
}

var ErrDiscovery = errors.New("discovery failed")
var ErrPostEvent = errors.New("posting an event failed")
var ErrValidation = errors.New("validation failed")
