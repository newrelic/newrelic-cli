package utils

import (
	"testing"

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

func TestIsJSONString_True(t *testing.T) {
	str := "{\"message\":\"original message\",\"metadata\":{\"someKey\":\"some value\"}}"
	result := IsJSONString(str)
	require.True(t, result)
}

func TestIsJSONString_False(t *testing.T) {
	str := "some string"
	result := IsJSONString(str)
	require.False(t, result)
}
