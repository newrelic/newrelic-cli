//go:build unit

package utils

import (
	"errors"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructToMap(t *testing.T) {
	t.Parallel()

	// Must have json tags.
	type testStruct struct {
		Name    string  `json:"name,omitempty"`
		GUID    string  `json:"guid,omitempty"`
		Exclude bool    `json:"exclude,omitempty"`
		ID      int     `json:"id,omitempty"`
		Float   float64 `json:"float,omitempty"`
	}

	item := testStruct{
		Name:    "example",
		GUID:    "abc123",
		Exclude: true,
		ID:      1234,
		Float:   1.123,
	}

	fieldsFilter := []string{"name", "guid", "id", "float"}

	result := StructToMap(item, fieldsFilter)

	expected := map[string]interface{}{
		"name":  "example",
		"guid":  "abc123",
		"id":    1234,
		"float": 1.123,
	}

	assert.Equal(t, expected, result)
}

func TestIsAbsoluteURL(t *testing.T) {
	urls := []struct {
		URL        string
		IsAbsolute bool
	}{
		{
			URL:        "https://one.newrelic.com",
			IsAbsolute: true,
		},
		{
			URL:        "one.newrelic.com",
			IsAbsolute: false,
		},
		{
			URL:        "http://localhost:18003/v1/status/entity",
			IsAbsolute: true,
		},
		{
			URL:        "/v1/status/entity",
			IsAbsolute: false,
		},
	}

	for _, u := range urls {
		result := IsAbsoluteURL(u.URL)
		require.Equal(t, u.IsAbsolute, result)
	}
}

func TestLogIfFatal(t *testing.T) {
	cases := []struct {
		param       error
		expectFatal bool
	}{
		{
			param:       nil,
			expectFatal: false,
		},
		{
			param:       errors.New("invalid"),
			expectFatal: true,
		},
		{
			param:       errors.New("403"),
			expectFatal: true,
		},
	}

	defer func() { log.StandardLogger().ExitFunc = nil }()
	var fatal bool
	log.StandardLogger().ExitFunc = func(int) { fatal = true }

	for _, c := range cases {
		fatal = false
		LogIfFatal(c.param)
		require.Equal(t, c.expectFatal, fatal)
	}
}

func TestIsValidUserAPIKeyFormat_Valid(t *testing.T) {
	result := IsValidUserAPIKeyFormat("NRAK-ABCDEBFGJIJKLMNOPQRSTUVWXYZ")
	assert.True(t, result)
}

func TestIsValidUserAPIKeyFormat_Invalid(t *testing.T) {
	// Invalid prefix
	result := IsValidUserAPIKeyFormat("NRBR-ABCDEBFGJIJKLMNOPQRSTUVWXYZ")
	assert.False(t, result)

	// A license key is not a valid User API key (note, this is not a real license key)
	result = IsValidUserAPIKeyFormat("4321abcsksd344ndlsdm20231mwd21230md12cbsdhk2")
	assert.False(t, result)

	// Special characters are invaid after the prefix hyphen
	result = IsValidUserAPIKeyFormat("NRAK-@$%^!")
	assert.False(t, result)
}

func TestIsValidLicenseKeyFormat_Valid(t *testing.T) {
	result := IsValidLicenseKeyFormat("0123456789abcdefABCDEF0123456789abcdNRAL")
	assert.True(t, result)
}

func TestIsValidLicenseKeyFormat_Valid_EU(t *testing.T) {
	result := IsValidLicenseKeyFormat("eu01xx6789abcdefABCDEF0123456789abcdNRAL")
	assert.True(t, result)
}

func TestIsValidLicenseKeyFormat_Invalid(t *testing.T) {
	// Invalid length
	result := IsValidLicenseKeyFormat("0123456789abcefNRAL")
	assert.False(t, result)

	// Invalid suffix
	result = IsValidLicenseKeyFormat("0123456789abcdefABCDEF0123456789abcdNRAK")
	assert.False(t, result)

	// Invalid characters
	result = IsValidLicenseKeyFormat("0123456789ghijklGHIJKL0123456789ghijNRAL")
	assert.False(t, result)

	// Invalid characters EU
	result = IsValidLicenseKeyFormat("eu01xx6789ghijklGHIJKL0123456789ghijNRAL")
	assert.False(t, result)
}
