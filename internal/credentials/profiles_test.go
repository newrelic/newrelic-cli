// +build unit

package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfileMarshal(t *testing.T) {
	t.Parallel()

	p := Profile{
		APIKey: "testAPIKey",
		Region: "TEST",
	}

	// Ensure that the region name is Lowercase
	m, err := p.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `{"apiKey":"testAPIKey","region":"test"}`, string(m))
}
