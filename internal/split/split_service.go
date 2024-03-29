package split

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/splitio/go-client/v6/splitio/client"
	"github.com/splitio/go-client/v6/splitio/conf"

	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
)

// Split.io client API keys
// These keys are client-facing and are compiled in via Github secrets and ldflags.
// There is no security risk in exposing these keys, as the only purpose they
// serve is to retrieve experiments and can not be used to modify anything
// within our internal Split.io system.
var (
	prodKey     = "localhost"
	stagingKey  = "localhost"
	accountID   = configAPI.GetActiveProfileAccountID()
	splitConfig = conf.Default()
)

type Srvc struct {
	client *client.SplitClient
}

// Creates a new instance of a Split Factory
// Using "localhost" as the apiKey allows us to use Split.io
// in localhost mode as defined in their documentation
// If the service is not available we return nil on these methods:
// - Get
// - GetAll
func NewService(region string) (*Srvc, error) {
	var apiKey = getAPIKeyByRegion(region)
	if region == "localhost" {
		apiKey = "localhost"
	}

	splitConfig.LoggerConfig.Prefix = "[splitio/client]"
	splitConfig.LoggerConfig.WarningWriter = config.Logger.WriterLevel(log.DebugLevel)
	splitConfig.LoggerConfig.ErrorWriter = config.Logger.WriterLevel(log.DebugLevel)

	factory, err := client.NewSplitFactory(apiKey, splitConfig)
	if err != nil {
		return nil, fmt.Errorf("split SDK init error: %s", err)
	}

	client := factory.Client()
	err = client.BlockUntilReady(10)
	if err != nil {
		return nil, fmt.Errorf("split SDK timeout: %s", err)
	}

	return &Srvc{client: client}, nil
}

// Get a treatment from Split.io given the name of the split
func (s *Srvc) Get(split string) string {
	if s == nil {
		return ""
	}
	return s.client.Treatment(accountID, split, nil)
}

// Get all treatments from Split.io given a list of splits
func (s *Srvc) GetAll(splits []string) map[string]string {
	if s == nil {
		return nil
	}
	log.Debugf("Retrieving treatments for splits: %s", splits)
	return s.client.Treatments(accountID, splits, nil)
}

func getAPIKeyByRegion(region string) string {
	if strings.EqualFold(region, "staging") {
		return stagingKey
	}
	return prodKey
}
