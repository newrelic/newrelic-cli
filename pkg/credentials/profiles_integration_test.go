package credentials

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDefaultProfile(t *testing.T) {

	c := Credentials{
		Profiles: make(map[string]Profile),
	}

	c.Profiles["testCase"] = Profile{}

	err := c.SetDefaultProfile("notTestCase")
	assert.Error(t, err, "no profile found")

	err = c.SetDefaultProfile("testCase")
	assert.Error(t, err, "credential ConfigDirectory is empty")

	c.ConfigDirectory = "/tmp/newrelic"

	err = c.SetDefaultProfile("testCase")
	assert.Error(t, err)

	os.Mkdir(c.ConfigDirectory, 0700)
	defer os.RemoveAll(c.ConfigDirectory)

	err = c.SetDefaultProfile("testCase")
	assert.NoError(t, err)
}
