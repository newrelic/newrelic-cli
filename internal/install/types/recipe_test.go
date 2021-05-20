// +build unit

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToStringByFieldName(t *testing.T) {
	data := map[string]interface{}{
		"intField":    9600, // e.g. the port used for a db connection URL
		"stringField": "stringValue",
		"boolField":   false,
	}

	intAsString := toStringByFieldName("intField", data)
	require.Equal(t, "9600", intAsString)

	stringAsString := toStringByFieldName("stringField", data)
	require.Equal(t, "stringValue", stringAsString)

	boolAsString := toStringByFieldName("boolField", data)
	require.Equal(t, "false", boolAsString)
}
