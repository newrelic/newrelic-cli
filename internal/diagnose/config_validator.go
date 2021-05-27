package diagnose

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/shirou/gopsutil/host"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-cli/internal/utils/validation"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/apiaccess"
)

const (
	validationEventType          = "NrIntegrationError"
	DefaultPostRetryDelaySec     = 5
	DefaultPostMaxRetries        = 20
	DefaultMaxValidationAttempts = 20
)

type ConfigValidator struct {
	client *newrelic.NewRelic
	*validation.PollingNRQLValidator
	profile           *credentials.Profile
	PostRetryDelaySec int
	PostMaxRetries    int
}

type ValidationTracerEvent struct {
	EventType string `json:"eventType"`
	Hostname  string `json:"hostname"`
	Purpose   string `json:"purpose"`
	GUID      string `json:"guid"`
}

func NewConcreteConfigValidator(client *newrelic.NewRelic) *ConfigValidator {
	v := validation.NewPollingNRQLValidator(&client.Nrdb)
	v.MaxAttempts = DefaultMaxValidationAttempts

	return &ConfigValidator{
		client:               client,
		PollingNRQLValidator: v,
		profile:              credentials.DefaultProfile(),
		PostRetryDelaySec:    DefaultPostRetryDelaySec,
		PostMaxRetries:       DefaultPostMaxRetries,
	}
}

func (c *ConfigValidator) ValidateConfig(ctx context.Context) error {
	if err := c.validateKeys(c.profile); err != nil {
		return err
	}

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

	postEvent := func() error {
		if err = c.client.Events.CreateEvent(c.profile.AccountID, evt); err != nil {
			log.Error(err)
			return ErrPostEvent
		}

		return nil
	}

	r := utils.NewRetry(c.PostMaxRetries, c.PostRetryDelaySec, postEvent)
	if err = r.ExecWithRetries(); err != nil {
		return err
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

func (c *ConfigValidator) validateKeys(profile *credentials.Profile) error {
	validateKeyFunc := func() error {
		if err := c.validateLicenseKey(profile); err != nil {
			return err
		}

		if err := c.validateInsightsInsertKey(profile); err != nil {
			return err
		}
		return nil
	}

	r := utils.NewRetry(c.PostMaxRetries, c.PostRetryDelaySec, validateKeyFunc)
	if err := r.ExecWithRetries(); err != nil {
		return err
	}
	return nil
}

func (c *ConfigValidator) validateInsightsInsertKey(profile *credentials.Profile) error {
	insightsInsertKeys, err := c.client.APIAccess.ListInsightsInsertKeys(profile.AccountID)
	if err != nil {
		return fmt.Errorf(ErrConnectionStringFormat, err)
	}

	for _, k := range insightsInsertKeys {
		if k.Key == profile.InsightsInsertKey {
			return nil
		}
	}

	return ErrInsightsInsertKey
}

func (c *ConfigValidator) validateLicenseKey(profile *credentials.Profile) error {
	params := apiaccess.APIAccessKeySearchQuery{
		Scope: apiaccess.APIAccessKeySearchScope{
			AccountIDs: []int{profile.AccountID},
		},
		Types: []apiaccess.APIAccessKeyType{
			apiaccess.APIAccessKeyTypeTypes.INGEST,
		},
	}

	licenseKeys, err := c.client.APIAccess.SearchAPIAccessKeys(params)
	if err != nil {
		return fmt.Errorf(ErrConnectionStringFormat, err)
	}

	for _, k := range licenseKeys {
		if k.Key == profile.LicenseKey {
			return nil
		}
	}

	return ErrLicenseKey
}
