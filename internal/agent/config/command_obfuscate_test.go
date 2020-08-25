// +build unit

package apm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestApmProxy(t *testing.T) {
	assert.Equal(t, "config", cmdConfig.Name())

	testcobra.CheckCobraMetadata(t, cmdConfig)
}

func TestApmProxyObfuscate(t *testing.T) {
	assert.Equal(t, "obfuscate", cmdObfuscate.Name())

	testcobra.CheckCobraMetadata(t, cmdObfuscate)
	testcobra.CheckCobraRequiredFlags(t, cmdObfuscate, []string{})
}
