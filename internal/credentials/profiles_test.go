// +build unit

package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/pkg/region"
)

func TestProfileMarshal(t *testing.T) {
	t.Parallel()

	p := Profile{
		APIKey: "testAPIKey",
		Region: region.Name("TEST"),
	}

	// Ensure that the region name is Lowercase
	m, err := p.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `{"apiKey":"testAPIKey","region":"test"}`, string(m))
}
