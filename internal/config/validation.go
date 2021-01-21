package config

import (
	"fmt"
	"strings"

	"github.com/newrelic/newrelic-client-go/pkg/region"
)

func RequireInsightsInsertKey() (string, error) {
	return RequireActiveProfileFieldString(InsightsInsertKey)
}

func RequireUserKey() (string, error) {
	return RequireActiveProfileFieldString(UserKey)
}

func RequireAccountID() (int, error) {
	return RequireActiveProfileFieldInt(AccountID)
}

func RequireActiveProfileFieldString(key ProfileFieldKey) (string, error) {
	v := GetActiveProfileValueString(key)
	if v == "" {
		f := findProfileField(key)
		return "", fmt.Errorf("%s is required, set it in your default profile or use the %s environment variable", AccountID, f.EnvOverride)
	}

	return v, nil
}

func RequireActiveProfileFieldInt(key ProfileFieldKey) (int, error) {
	v := GetActiveProfileValueInt(key)
	if v == 0 {
		f := findProfileField(key)
		return 0, fmt.Errorf("%s is required, set it in your default profile or use the %s environment variable", key, f.EnvOverride)
	}

	return v, nil
}

func isValidConfigKey(key CfgFieldKey) bool {
	return findConfigField(key) != nil
}

func isValidProfileKey(key ProfileFieldKey) bool {
	return findProfileField(key) != nil
}

func validConfigFieldKeys() []string {
	valid := make([]string, len(ConfigFields))

	for _, v := range ConfigFields {
		valid = append(valid, string(v.Key))
	}

	return valid
}

func validProfileFieldKeys() []string {
	valid := make([]string, len(ProfileFields))

	for _, v := range ProfileFields {
		valid = append(valid, string(v.Key))
	}

	return valid
}

func stringInSlice(validVals []string, caseSensitive bool) func(interface{}) error {
	return func(val interface{}) error {
		for _, v := range validVals {

			if !caseSensitive && strings.EqualFold(v, val.(string)) {
				return nil
			}

			if v == val {
				return nil
			}
		}

		return fmt.Errorf("valid values are %s", validVals)
	}
}

func validRegions() []string {
	validRegions := []string{}
	for k := range region.Regions {
		validRegions = append(validRegions, string(k))
	}

	return validRegions
}

func isNumber() func(interface{}) error {
	return func(val interface{}) error {
		if _, ok := val.(int); ok {
			return nil
		}

		return fmt.Errorf("value is required to be numeric")
	}
}
