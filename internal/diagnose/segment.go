package diagnose

import (
	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/segment"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http/httpproxy"
)

func initSegment() *segment.Segment {
	accountID := configAPI.GetActiveProfileAccountID()
	region := configAPI.GetActiveProfileString(config.Region)
	isProxyConfigured := isProxyConfigured()
	writeKey, err := recipes.NewEmbeddedRecipeFetcher().GetSegmentWriteKey()
	if err != nil {
		log.Debug("segment: error reading write key, cannot write to segment", err)
		return segment.NewNoOp()
	}

	return segment.New(writeKey, accountID, region, isProxyConfigured)
}
func isProxyConfigured() bool {
	proxyConfig := httpproxy.FromEnvironment()
	return proxyConfig.HTTPProxy != "" || proxyConfig.HTTPSProxy != "" || proxyConfig.NoProxy != ""
}
