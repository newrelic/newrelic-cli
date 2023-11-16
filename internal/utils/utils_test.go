//go:build unit

package utils

import (
	"errors"
	"reflect"
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

func TestRemoveFromSliceWorks(t *testing.T) {
	type test struct {
		input  []string
		remove string
		want   []string
	}
	tests := []test{
		{input: []string{"a", "b", "c"}, remove: "c", want: []string{"a", "b"}},
		{input: []string{"a", "b", "c"}, remove: "b", want: []string{"a", "c"}},
		{input: []string{"a", "b"}, remove: "a", want: []string{"b"}},
		{input: []string{"a", "b"}, remove: "b", want: []string{"a"}},
		{input: []string{"a"}, remove: "a", want: []string{}},
		{input: []string{}, remove: "a", want: []string{}},
		{input: []string{}, remove: "", want: []string{}},
		{input: nil, remove: "a", want: nil},
	}

	for _, tc := range tests {
		got := RemoveFromSlice(tc.input, tc.remove)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}
