package discovery

import (
	"context"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

//go:generate mockgen -source=psutil_discoverer.go -destination=../discovery/mocks/mock_discoverer.go -package=mocks

type Discoverer interface {
	Discover(ctx context.Context) (*types.DiscoveryManifest, error)
}

//go:generate mockgen -source=validator.go -destination=../discovery/mocks/mock_validator.go -package=mocks

type Validator interface {
	Validate(manifest *types.DiscoveryManifest) error
}
