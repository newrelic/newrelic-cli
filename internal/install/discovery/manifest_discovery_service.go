package discovery

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func DiscoverManifest(ctx context.Context, discoverer Discoverer) (*types.DiscoveryManifest, error) {
	log.Debug("discovering system information")
	m, err := discoverer.Discover(ctx)
	if err != nil {
		return nil, fmt.Errorf("there was an error discovering system info: %s", err)
	}

	return m, nil
}

func ValidateManifest(manifest *types.DiscoveryManifest, validator Validator) error {
	err := validator.Validate(manifest)
	if err != nil {
		return err
	}
	log.Debugf("Done asserting valid operating system for OS:%s and PlatformVersion:%s", manifest.OS, manifest.PlatformVersion)
	return nil
}
