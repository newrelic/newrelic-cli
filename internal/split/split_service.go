package split

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/splitio/go-client/v6/splitio/client"
	"github.com/splitio/go-client/v6/splitio/conf"

	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
)

// Split.io client API keys
// These keys are client-facing and are only used to fetch splits.
// There is no security risk in exposing these keys, as the only purpose they
// serve is to retrieve experiments and can not be used to modify anything
// within our internal Split.io system.
var (
	prodKey                          = "8me2vu6v8lhssdkrpenp1uunl9s3bdc8njqp"
	stagingKey                       = "mcf9oimts3laqli01e2ktrjdudkdbh8dg42a"
	accountID                        = configAPI.GetActiveProfileAccountID()
	splitConfig *conf.SplitSdkConfig = conf.Default()
)

type Srvc struct {
	client *client.SplitClient
}

// Creates a new instance of a Split Factory
// Using "localhost" as the apiKey allows us to use Split.io
// in localhost mode as defined in their documentation
func NewService(region string) (*Srvc, error) {
	var apiKey = getAPIKeyByRegion(region)
	if region == "localhost" {
		apiKey = "localhost"
	}

	factory, err := client.NewSplitFactory(apiKey, splitConfig)
	if err != nil {
		log.Errorf("Split SDK init error: %s\n", err)
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
	return s.client.Treatment(accountID, split, nil)
}

// Get all treatments from Split.io given a list of splits
func (s *Srvc) GetAll(splits []string) map[string]string {
	log.Debugf("Retrieving treatments for splits: %s", splits)
	return s.client.Treatments(accountID, splits, nil)
}

func getAPIKeyByRegion(region string) string {
	if strings.EqualFold(region, "staging") {
		return stagingKey
	}
	return prodKey
}
