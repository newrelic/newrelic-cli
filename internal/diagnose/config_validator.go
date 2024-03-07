package diagnose

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/shirou/gopsutil/v3/host"

	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/segment"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-cli/internal/utils/validation"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/apiaccess"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
)

const (
	validationEventType          = "NrIntegrationError"
	DefaultMaxValidationAttempts = 20
)

var sg = segment.Init()

type ConfigValidator struct {
	client *newrelic.NewRelic
	*validation.PollingNRQLValidator
	PostRetryDelaySec int
	PostMaxRetries    int
}

type ValidationTracerEvent struct {
	EventType string `json:"eventType"`
	Hostname  string `json:"hostname"`
	Purpose   string `json:"purpose"`
	GUID      string `json:"guid"`
}

func NewConfigValidator(client *newrelic.NewRelic) *ConfigValidator {
	v := validation.NewPollingNRQLValidator(&client.Nrdb)
	v.MaxAttempts = DefaultMaxValidationAttempts

	return &ConfigValidator{
		client:               client,
		PollingNRQLValidator: v,
		PostRetryDelaySec:    config.DefaultPostRetryDelaySec,
		PostMaxRetries:       config.DefaultMaxTimeoutSeconds / config.DefaultPostRetryDelaySec,
	}
}

func (c *ConfigValidator) Validate(ctx context.Context) error {
	accountID := configAPI.GetActiveProfileAccountID()

	if err := c.validateKeys(ctx); err != nil {
		return err
	}

	i, err := host.InfoWithContext(ctx)
	if err != nil {
		log.Debug(err)
		return types.NewDetailError(types.EventTypes.UnableToDiscover, types.ErrDiscovery.Error())
	}

	evt := ValidationTracerEvent{
		EventType: validationEventType,
		Hostname:  i.Hostname,
		Purpose:   "New Relic CLI configuration validation",
		GUID:      uuid.NewString(),
	}

	postEvent := func() error {
		if err = c.client.Events.CreateEventWithContext(ctx, accountID, evt); err != nil {
			log.Debug(err)

			if e, ok := err.(*nrErrors.PaymentRequiredError); ok {
				return e
			}

			return types.NewDetailError(types.EventTypes.UnableToPostData, err.Error())
		}

		return nil
	}

	r := utils.NewRetry(c.PostMaxRetries, c.PostRetryDelaySec*1000, postEvent)
	retryCtx := r.ExecWithRetries(ctx)

	if !retryCtx.Success {
		return types.NewDetailError(types.EventTypes.InvalidIngestKey, retryCtx.MostRecentError().Error())
	}

	query := fmt.Sprintf(`
	FROM %s
	SELECT count(*)
	WHERE hostname LIKE '%s%%'
	AND guid = '%s'
	SINCE 10 MINUTES AGO
	`, evt.EventType, evt.Hostname, evt.GUID)

	sg.Track("ValidateNrqlStart")

	start := time.Now()
	_, err = c.PollingNRQLValidator.Validate(ctx, query)
	durationMs := time.Since(start).Milliseconds()

	ei := segment.NewEventInfo("ValidateNrqlEnd", "")
	ei.WithAdditionalInfo("durationMs", durationMs)
	ei.WithAdditionalInfo("hasError", err != nil)
	sg.TrackInfo(ei)

	if err != nil {
		log.Debug(err)
		return types.NewDetailError(types.EventTypes.UnableToLocatePostedData, types.ErrValidation.Error()+" "+err.Error())
	}

	return err
}

func (c *ConfigValidator) validateKeys(ctx context.Context) error {
	validateKeyFunc := func() error {
		return c.validateLicenseKey(ctx)
	}

	r := utils.NewRetry(c.PostMaxRetries, c.PostRetryDelaySec*1000, validateKeyFunc)
	retryCtx := r.ExecWithRetries(ctx)

	if !retryCtx.Success {
		return retryCtx.MostRecentError()
	}

	return nil
}

func (c *ConfigValidator) validateLicenseKey(ctx context.Context) error {
	accountID := configAPI.GetActiveProfileAccountID()
	licenseKey := configAPI.GetActiveProfileString(config.LicenseKey)
	params := apiaccess.APIAccessKeySearchQuery{
		Scope: apiaccess.APIAccessKeySearchScope{
			AccountIDs: []int{accountID},
		},
		Types: []apiaccess.APIAccessKeyType{
			apiaccess.APIAccessKeyTypeTypes.INGEST,
		},
	}

	licenseKeys, err := c.client.APIAccess.SearchAPIAccessKeysWithContext(ctx, params)
	if err != nil {
		return types.NewDetailError(types.EventTypes.UnableToConnect, fmt.Sprintf(types.ErrConnectionStringFormat, err))
	}

	for _, k := range licenseKeys {
		if k.Key == licenseKey {
			return nil
		}
	}

	return types.NewDetailError(types.EventTypes.InvalidIngestKey, types.ErrLicenseKey.Error())
}
