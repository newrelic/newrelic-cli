// build +unit

package obfuscate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObfuscateStringWithKey(t *testing.T) {
	t.Parallel()

	var result string

	//empty text to obfuscate
	result = ObfuscateStringWithKey("", "RandomKey")
	assert.Equal(t, "", result)

	//empty key
	result = ObfuscateStringWithKey("RandomString", "")
	assert.Equal(t, "", result)

	//empty text to obfuscate and empty key
	result = ObfuscateStringWithKey("", "")
	assert.Equal(t, "", result)

	//text to obfuscate with length longer than key
	result = ObfuscateStringWithKey("ThisIs18Characters", "ThisIs13Chars")
	assert.NotEmpty(t, result)

	//text to obfuscate with length shorter than key
	result = ObfuscateStringWithKey("ThisIs13Chars", "ThisIs18Characters")
	assert.NotEmpty(t, result)

	//expected results
	result = ObfuscateStringWithKey("XYZ", "123")
	assert.Equal(t, "aWtp", result)

	result = ObfuscateStringWithKey("XYZ", "456")
	assert.Equal(t, "bGxs", result)

	result = ObfuscateStringWithKey("国字 kokuji", "123")
	assert.Equal(t, "1KmO1J+kEVlcWkdZWA==", result)

}
