package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

// GetActiveProfileName retrieves the profile in use for this command execution.
// To retrieve the active profile, the following criteria are evaluated in order,
// short circuiting and returning the described value if true:
// 1. a profile has been provided with the global `--profile` flag
// 2. a profile is set in the default profile config file
// 3. "default" is returned if none of the above are true
func GetActiveProfileName() string {
	if config.FlagProfileName != "" {
		return config.FlagProfileName
	}

	profileName, err := GetDefaultProfileName()
	if err != nil || profileName == "" {
		return config.DefaultProfileName
	}

	return profileName
}

// GetDefaultProfileName retrieves the profile set in the default profile config
// file.  If the file does not exist, an empty string will be returned.
func GetDefaultProfileName() (string, error) {
	defaultProfileFilePath := filepath.Join(config.BasePath, config.DefaultProfileFileName)
	data, err := ioutil.ReadFile(defaultProfileFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}

		return "", err
	}

	return strings.Trim(string(data), "\""), nil
}

// GetProfileNames retrieves all profile names currently configured in the credentials file.
func GetProfileNames() []string {
	return config.CredentialsProvider.GetScopes()
}

// SetDefaultProfile sets the given profile as the new default in the default profile
// config file. If the given profile does not exist, the set operation will return
// an error.
func SetDefaultProfile(profileName string) error {
	if ok := utils.StringInSlice(profileName, GetProfileNames()); !ok {
		return fmt.Errorf("profile %s does not exist", profileName)
	}

	defaultProfileFilePath := filepath.Join(config.BasePath, config.DefaultProfileFileName)
	return ioutil.WriteFile(defaultProfileFilePath, []byte("\""+profileName+"\""), 0644)
}

// RemoveProfile removes a profile from the credentials file.  If the profile being
// removed is the default, it will attempt to find another profile to set as the
// new default. If another profile cannot be found, the default profile config file
// will be deleted.
func RemoveProfile(profileName string) error {
	if err := config.CredentialsProvider.RemoveScope(profileName); err != nil {
		return err
	}

	d, err := GetDefaultProfileName()
	if err != nil {
		return err
	}

	if d == profileName {
		names := GetProfileNames()
		if len(names) > 0 {
			if err = SetDefaultProfile(names[0]); err != nil {
				return fmt.Errorf("could not set new default profile")
			}

			log.Infof("setting %s as the new default profile", names[0])
		} else {
			if err := removeDefaultProfile(); err != nil {
				return fmt.Errorf("could not delete default profile")
			}
		}
	}

	return nil
}

// RequireActiveProfileAccountID retrieves the currently configured account ID,
// returning an error if the value retrieved is the zero value.
// When returning an account ID, the following will be evaluated in order, short-circuiting
// and returning the described value if true:
// 1. An environment variable override has been set with NEW_RELIC_ACCOUNT_ID
// 2. An account ID has been provided with the `--accountId` global flag
// 3. An account ID has been set in the active profile
func RequireActiveProfileAccountID() int {
	v := GetActiveProfileAccountID()
	if v == 0 {
		log.Fatalf("%s is required", config.AccountID)
	}

	return v
}

// GetActiveProfileAccountID retrieves the currently configured account ID.
// When returning an account ID, the following will be evaluated in order, short-circuiting
// and returning the described value if true:
// 1. An environment variable override has been set with NEW_RELIC_ACCOUNT_ID
// 2. An account ID has been provided with the `--accountId` global flag
// 3. An account ID has been set in the active profile
// 4. The zero value will be returned if none of the above are true
func GetActiveProfileAccountID() int {
	return getActiveProfileIntWithOverride(config.AccountID, config.FlagAccountID)
}

// GetActiveProfileString retrieves the value set for the given key in the active
// profile, if any. Environment variable overrides will be preferred over values
// set in the active profile, and a default value will be returned if it has been
// configured and no value has been set for the key in the active profile.
// An attempt will be made to convert the underlying value to a string if is not
// already stored that way. Failing the above, the zero value wil be returned.
func GetActiveProfileString(key config.FieldKey) string {
	return GetProfileString(GetActiveProfileName(), key)
}

// GetProfileString retrieves the value set for the given key and profile, if any.
// Environment variable overrides will be preferred over values set in the given
// profile, and a default value will be returned if it has been configured and no
// value has been set for the key in the given profile.
// An attempt will be made to convert the underlying value to a string if is not
// already stored that way.  Failing the above, the zero value wil be returned.
func GetProfileString(profileName string, key config.FieldKey) string {
	v, err := config.CredentialsProvider.GetStringWithScope(profileName, key)
	if err != nil {
		return ""
	}

	return v
}

// GetProfileInt retrieves the value set for the given key and profile, if any.
// Environment variable overrides will be preferred over values set in the given
// profile, and a default value will be returned if it has been configured and no
// value has been set for the key in the given profile.
// An attempt will be made to convert the underlying value to an int if is not
// already stored that way.  Failing the above, the zero value wil be returned.
func GetProfileInt(profileName string, key config.FieldKey) int {
	v, err := config.CredentialsProvider.GetIntWithScope(profileName, key)
	if err != nil {
		return 0
	}

	return int(v)
}

// SetProfileValue sets a value for the given key and profile.
func SetProfileValue(profileName string, key config.FieldKey, value interface{}) error {
	return config.CredentialsProvider.SetWithScope(profileName, key, value)
}

// GetLogLevel retrieves the currently configured log level.
// When returning a log level, the following will be evaluated in order, short-circuiting
// and returning the described value if true:
// 1. An environment variable override has been set with NEW_RELIC_CLI_LOG_LEVEL
// 2. A log level has been provided with the `--trace` global flag
// 2. A log level has been provided with the `--debug` global flag
// 3. A log level has been set in the config file
// 4. If none of the above is true, the default log level will be returned.
func GetLogLevel() string {
	if config.FlagDebug {
		return "debug"
	}

	if config.FlagTrace {
		return "trace"
	}

	return GetConfigString(config.LogLevel)
}

// GetConfigString retrieves the config value set for the given key, if any.
// Environment variable overrides will be preferred over values set in the given
// profile, and a default value will be returned if it has been configured and no
// value has been set for the key in the config file.
// An attempt will be made to convert the underlying value to a string if is not
// already stored that way.  Failing the above, the zero value wil be returned.
func GetConfigString(key config.FieldKey) string {
	v, err := config.ConfigStore.GetString(key)
	if err != nil {
		return ""
	}

	return v
}

// GetConfigTernary retrieves the config value set for the given key, if any.
// Environment variable overrides will be preferred over values set in the given
// profile, and a default value will be returned if it has been configured and no
// value has been set for the key in the config file.
// An attempt will be made to convert the underlying value to a Ternary if is not
// already stored that way.  Failing the above, the zero value wil be returned.
func GetConfigTernary(key config.FieldKey) config.Ternary {
	v, err := config.ConfigStore.GetTernary(key)
	if err != nil {
		return config.Ternary("")
	}

	return v
}

// SetConfigValue sets a config value for the given key.
func SetConfigValue(key config.FieldKey, value interface{}) error {
	return config.ConfigStore.Set(key, value)
}

// DeleteConfigValue deletes a config value for the given key.
func DeleteConfigValue(key config.FieldKey) error {
	return config.ConfigStore.DeleteKey(key)
}

// GetConfigFieldDefinition retrieves the field definition for the given config key.
func GetConfigFieldDefinition(key config.FieldKey) *config.FieldDefinition {
	return config.ConfigStore.GetFieldDefinition(key)
}

// ForEachProfileFieldDefinition iterates the field definitions for the profile fields.
func ForEachProfileFieldDefinition(profileName string, fn func(d config.FieldDefinition)) {
	config.CredentialsProvider.ForEachFieldDefinition(fn)
}

// ForEachConfigFieldDefinition iterates the field definitions for the config fields.
func ForEachConfigFieldDefinition(fn func(d config.FieldDefinition)) {
	config.ConfigStore.ForEachFieldDefinition(fn)
}

// GetValidConfigFieldKeys returns all the config field keys that can be set.
func GetValidConfigFieldKeys() (fieldKeys []config.FieldKey) {
	config.ConfigStore.ForEachFieldDefinition(func(fd config.FieldDefinition) {
		fieldKeys = append(fieldKeys, fd.Key)
	})

	return fieldKeys
}

func getActiveProfileIntWithOverride(key config.FieldKey, override int) int {
	return getProfileIntWithOverride(GetActiveProfileName(), key, override)
}

func getProfileIntWithOverride(profileName string, key config.FieldKey, override int) int {
	o := int64(override)
	v, err := config.CredentialsProvider.GetIntWithScopeAndOverride(profileName, key, &o)
	if err != nil {
		return 0
	}

	return int(v)
}

func removeDefaultProfile() error {
	defaultProfileFilePath := filepath.Join(config.BasePath, config.DefaultProfileFileName)
	return os.Remove(defaultProfileFilePath)
}
