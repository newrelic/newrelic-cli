// +build unit

package configuration

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testTernaryInfo = []struct {
	Name  Ternary
	Value string
	Bool  bool
	Err   error
}{
	{
		Name:  TernaryValues.Allow,
		Value: "ALLOW",
		Bool:  true,
	},
	{
		Name:  TernaryValues.Disallow,
		Value: "DISALLOW",
		Bool:  false,
	},
	{
		Name:  TernaryValues.Unknown,
		Value: "NOT_ASKED",
		Bool:  false,
	},
	{
		Name:  Ternary("invalid"),
		Value: "invalid",
		Bool:  false,
		Err:   errors.New("\"invalid\" is not a valid value; Please use one of: {ALLOW DISALLOW NOT_ASKED}"),
	},
}

func TestTernaryString(t *testing.T) {
	t.Parallel()

	// Set the valid pre-release feature values
	for _, info := range testTernaryInfo {
		assert.Equal(t, info.Value, info.Name.String())
	}
}

func TestTernaryValid(t *testing.T) {
	t.Parallel()

	// Set the valid pre-release feature values
	for _, info := range testTernaryInfo {
		assert.Equal(t, info.Err, info.Name.Valid())
	}

}

func TestTernaryBool(t *testing.T) {
	t.Parallel()

	for _, info := range testTernaryInfo {
		assert.Equal(t, info.Bool, info.Name.Bool())
	}

	// Invalid data
	assert.Equal(t, false, Ternary("asdf").Bool())
}
