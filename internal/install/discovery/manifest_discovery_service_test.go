package discovery

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/newrelic/newrelic-cli/internal/install/discovery/mocks"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestDiscoverManifestReturnsValidManifest(t *testing.T) {
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.On("Discover", mock.Anything).Return(&types.DiscoveryManifest{
		OS:              "darwin",
		PlatformVersion: "10.13",
	}, nil)

	manifest, err := DiscoverManifest(context.Background(), mockDiscoverer)

	assert.NoError(t, err)
	assert.NotNil(t, manifest)
	assert.Equal(t, "darwin", manifest.OS)
	assert.Equal(t, "10.13", manifest.PlatformVersion)
}

func TestDiscoverManifestErrorsOnInvalidManifest(t *testing.T) {
	expectedError := errors.New("discovery errored")
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.On("Discover", mock.Anything).Return(nil, expectedError)

	manifest, err := DiscoverManifest(context.Background(), mockDiscoverer)

	assert.Nil(t, manifest)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "discovery errored"))
}

func TestValidateManifestReturnsNoErrorOnValidManifest(t *testing.T) {
	mockValidator := mocks.NewValidator(t)
	mockValidator.On("Validate", mock.Anything).Return(nil)

	err := ValidateManifest(&types.DiscoveryManifest{OS: "valid"}, mockValidator)

	assert.NoError(t, err)
}

func TestValidateManifestReturnsErrorOnInvalidManifest(t *testing.T) {
	expectedError := errors.New("validation error")
	mockValidator := mocks.NewValidator(t)
	mockValidator.On("Validate", mock.Anything).Return(expectedError)

	err := ValidateManifest(&types.DiscoveryManifest{OS: "invalid"}, mockValidator)

	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "validation error"))
}
