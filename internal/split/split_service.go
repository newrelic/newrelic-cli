package split

import (
	"fmt"
	"io/ioutil"
	"os"
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
	prodKey    = "8me2vu6v8lhssdkrpenp1uunl9s3bdc8njqp"
	stagingKey = "mcf9oimts3laqli01e2ktrjdudkdbh8dg42a"
	accountID  = configAPI.GetActiveProfileAccountID()
)

type Service struct {
	client *client.SplitClient
}

// Creates a new instance of a Split Factory
// Using "localhost" as the apiKey allows us to use Split.io
// in localhost mode as defined in their documentation
func NewService(region string) (*Service, error) {
	var apiKey = getAPIKeyByRegion(region)
	cfg := conf.Default()
	if region == "localhost" {
		apiKey = "localhost"
		cfg.SplitFile = createMockSplits()
	}

	factory, err := client.NewSplitFactory(apiKey, cfg)
	if err != nil {
		log.Errorf("Split SDK init error: %s\n", err)
		return nil, fmt.Errorf("split SDK init error: %s", err)
	}

	client := factory.Client()
	err = client.BlockUntilReady(10)
	if err != nil {
		return nil, fmt.Errorf("split SDK timeout: %s", err)
	}

	return &Service{client: client}, nil
}

// Get a treatment from Split.io given the name of the split
func (s *Service) Get(split string) string {
	return s.client.Treatment(accountID, split, nil)
}

// Get all treatments from Split.io given a list of splits
func (s *Service) GetAll(splits []string) map[string]string {
	log.Debugf("Retrieving treatments for splits: %s", splits)
	return s.client.Treatments(accountID, splits, nil)
}

// Creates a temporary file with splits used for unit-testing
func createMockSplits() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("could not get user home directory: %s", err)
	}
	// Create a temporary file that holds test splits for testing purposes
	blob := []byte(MockSplits)
	filename := dir + "/mock.split"
	err = ioutil.WriteFile(filename, blob, 0777)
	if err != nil {
		log.Errorf("could not create temp file: %s", err)
	}

	return filename
}

func getAPIKeyByRegion(region string) string {
	if strings.EqualFold(region, "staging") {
		return stagingKey
	}
	return prodKey
}
