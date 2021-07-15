package split

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
)

var (
	SplitService *Service
)

func Init() {
	region := configAPI.GetConfigString(config.Region)
	apiKey := GetAPIKeyByRegion(strings.ToLower(region))

	service, err := NewSplitService(apiKey, region)
	if err != nil {
		log.Errorf("could not initialize SplitService: %s\n", err)
	}
	SplitService = service
}

func init() {
	Init()
}
