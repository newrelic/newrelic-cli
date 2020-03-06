package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
