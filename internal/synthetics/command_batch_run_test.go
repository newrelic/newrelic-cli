//go:build unit

package synthetics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

func TestCreateConfigurationUsingGUIDs_ReturnsIsBlockingSetToTrueForAllMonitors(t *testing.T) {
	guids := []string{
		"abc",
		"xyz",
	}

	config := createConfigurationUsingGUIDs(guids)

	for i, guid := range guids {
		assert.Equal(t, config.Tests[i].MonitorGUID, synthetics.EntityGUID(guid))
		assert.Equal(t, config.Tests[i].Config.IsBlocking, true)
	}
}

func TestGetTestsMissingIsBlockingInConfig_ShouldReturnFirstTestOnly(t *testing.T) {
	var data = `
config:
  batchName: test
  branch: dev
  commit: abc123
  deepLink: https://example.com
  platform: CLI
  repository: repo
tests:
  - monitorGuid: NDE3N
    config:
  - monitorGuid: ZDFOQ
    config:
      isBlocking: false
  - monitorGuid: USHxN
    config:
      isBlocking: true
      overrides:
        location: east-1
        domain:
          - domain: http://example.org
            override: https://example.com
`
	indexes, err := getTestsMissingIsBlockingInConfig([]byte(data))

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(indexes))
	assert.Equal(t, 0, indexes[0])
}

func TestGetTestsMissingIsBlockingInConfig_ShouldReturnLastTestOnly(t *testing.T) {
	var data = `
config:
  batchName: test
  branch: dev
  commit: abc123
  deepLink: https://example.com
  platform: CLI
  repository: repo
tests:
  - monitorGuid: NDE3N
    config:
      isBlocking: true
  - monitorGuid: ZDFOQ
    config:
      isBlocking: false
  - monitorGuid: USHxN
    config:
      overrides:
        location: east-1
        domain:
          - domain: http://example.org
            override: https://example.com
`
	indexes, err := getTestsMissingIsBlockingInConfig([]byte(data))

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(indexes))
	assert.Equal(t, 2, indexes[0])
}

func TestGetTestsMissingIsBlockingInConfig_ShouldReturnThreeTests(t *testing.T) {
	var data = `
config:
  batchName: test
  branch: dev
  commit: abc123
tests:
  - monitorGuid: NDE3N
    config:
  - monitorGuid: ZDFOQ
    config:
  - monitorGuid: USHxN
    config:
      overrides:
        location: east-1
        domain:
          - domain: http://example.com
            override: https://example.com
`
	indexes, err := getTestsMissingIsBlockingInConfig([]byte(data))

	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(indexes))
}

func TestGetTestsMissingIsBlockingInConfig_ShouldReturnNone(t *testing.T) {
	var data = `
config:
  batchName: test
  branch: dev
  commit: abc123
tests:
  - monitorGuid: NDE3N
    config:
      isBlocking: false
      overrides:
        location: east-1
  - monitorGuid: ZDFOQ
    config:
      isBlocking: false
      overrides:
        location: east-1
        domain:
          - domain: http://example.com
            override: https://example.com
  - monitorGuid: USHxN
    config:
      isBlocking: true
`
	indexes, err := getTestsMissingIsBlockingInConfig([]byte(data))

	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(indexes))
}
