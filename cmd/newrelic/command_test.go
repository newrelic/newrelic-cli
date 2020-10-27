// +build unit

package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	assert.NotEmptyf(t, Command.Use, "Need to set Command.%s on Command %s", "Use", Command.CalledAs())
	assert.NotEmptyf(t, Command.Short, "Need to set Command.%s on Command %s", "Short", Command.CalledAs())
}

func TestInitialization(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	// // Initialize the new configuration directory
	// c, err := LoadCredentials(f)
	// assert.NoError(t, err)
	// assert.Equal(t, len(c.Profiles), 0)
	// assert.Equal(t, c.DefaultProfile, "")
	// assert.Equal(t, c.ConfigDirectory, f)

	// // Adding a profile with the same name should result in an error
	// err = c.AddProfile("testCase1", "us", "foot", "")
	// assert.Error(t, err)
	// assert.Equal(t, len(c.Profiles), 1)
	// t st -sassert.True(t, c.profileExists("testCase1"))
}
