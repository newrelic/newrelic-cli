package configuration

import (
	"fmt"
	"strings"
)

// Ternary is like a bool, but includes the unknown state
type Ternary string

// TernaryValues provides the set of Ternary values
var TernaryValues = struct {
	// Allow the option
	Allow Ternary

	// Disallow the option
	Disallow Ternary

	// Unknown is the unknown state
	Unknown Ternary
}{
	Allow:    "ALLOW",
	Disallow: "DISALLOW",
	Unknown:  "NOT_ASKED",
}

// Valid returns true for a valid value, false otherwise
func (t Ternary) Valid() error {
	val := string(t)

	if strings.EqualFold(val, string(TernaryValues.Allow)) ||
		strings.EqualFold(val, string(TernaryValues.Disallow)) ||
		strings.EqualFold(val, string(TernaryValues.Unknown)) {
		return nil
	}

	return fmt.Errorf("\"%s\" is not a valid value; Please use one of: %s", val, TernaryValues)
}

// Bool returns true if the ternary is set and contains the true value
func (t Ternary) Bool() bool {
	return strings.EqualFold(t.String(), TernaryValues.Allow.String())
}

// String returns the string value of the ternary
func (t Ternary) String() string {
	return string(t)
}
