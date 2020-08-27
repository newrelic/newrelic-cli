// build +unit

package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObfuscateStringWithKey(t *testing.T) {
	t.Parallel()

	var result string

	//empty text to obfuscate
	result = obfuscateStringWithKey("", "RandomKey")
	assert.Equal(t, "", result)

	//empty key
	result = obfuscateStringWithKey("RandomString", "")
	assert.Equal(t, "", result)

	//empty text to obfuscate and empty key
	result = obfuscateStringWithKey("", "")
	assert.Equal(t, "", result)

	//text to obfuscate with length longer than key
	result = obfuscateStringWithKey("ThisIs18Characters", "ThisIs13Chars")
	assert.NotEmpty(t, result)

	//text to obfuscate with length shorter than key
	result = obfuscateStringWithKey("ThisIs13Chars", "ThisIs18Characters")
	assert.NotEmpty(t, result)

	//expected results
	result = obfuscateStringWithKey("XYZ", "123")
	assert.Equal(t, "aWtp", result)

	result = obfuscateStringWithKey("XYZ", "456")
	assert.Equal(t, "bGxs", result)

	result = obfuscateStringWithKey("国字 kokuji", "123")
	assert.Equal(t, "1KmO1J+kEVlcWkdZWA==", result)

}
