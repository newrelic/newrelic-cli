package split

import (
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
)

var Service *Srvc

func Init() {
	region := configAPI.GetConfigString(config.Region)
	service, err := NewService(region)
	if err != nil {
		log.Fatalf("could not initialize SplitService: %s\n", err)
	}
	Service = service
}

func init() {
	Init()
}
