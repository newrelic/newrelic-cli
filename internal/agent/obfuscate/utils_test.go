// build +unit

package obfuscate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringWithKey(t *testing.T) {
	t.Parallel()

	var result string

	//empty text to obfuscate
	result = StringWithKey("", "RandomKey")
	assert.Equal(t, "", result)

	//empty key
	result = StringWithKey("RandomString", "")
	assert.Equal(t, "", result)

	//empty text to obfuscate and empty key
	result = StringWithKey("", "")
	assert.Equal(t, "", result)

	//text to obfuscate with length longer than key
	result = StringWithKey("ThisIs18Characters", "ThisIs13Chars")
	assert.NotEmpty(t, result)

	//text to obfuscate with length shorter than key
	result = StringWithKey("ThisIs13Chars", "ThisIs18Characters")
	assert.NotEmpty(t, result)

	//expected results
	result = StringWithKey("XYZ", "123")
	assert.Equal(t, "aWtp", result)

	result = StringWithKey("XYZ", "456")
	assert.Equal(t, "bGxs", result)

	result = StringWithKey("国字 kokuji", "123")
	assert.Equal(t, "1KmO1J+kEVlcWkdZWA==", result)

}
