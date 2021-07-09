package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/newrelic/newrelic-cli/internal/config"
)

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

func GetActiveProfileString(key config.FieldKey) string {
	return GetProfileString(GetActiveProfileName(), key)
}

func RequireActiveProfileString(key config.FieldKey) string {
	v := GetProfileString(GetActiveProfileName(), key)
	if v == "" {
		log.Fatalf("%s is required", key)
	}

	return v
}

func GetActiveProfileValue(profileName string, key config.FieldKey) interface{} {
	return GetProfileValue("", key)
}
func GetProfileValue(profileName string, key config.FieldKey) interface{} {
	v, err := config.CredentialsProvider.GetWithScope(profileName, key)
	if err != nil {
		return nil
	}

	return v
}

func GetProfileString(profileName string, key config.FieldKey) string {
	v, err := config.CredentialsProvider.GetStringWithScope(profileName, key)
	if err != nil {
		return ""
	}

	return v
}

func GetLogLevelWithFlagOverride() string {
	var override string
	if config.FlagDebug {
		override = "debug"
	}

	if config.FlagTrace {
		override = "trace"
	}

	return GetConfigStringWithOverride(config.LogLevel, override)
}

func RequireActiveProfileAccountIDWithFlagOverride() int {
	v := GetActiveProfileAccountIDWithFlagOverride()
	if v == 0 {
		log.Fatalf("%s is required", config.AccountID)
	}

	return v
}

func GetActiveProfileAccountIDWithFlagOverride() int {
	return GetActiveProfileIntWithOverride(config.AccountID, config.FlagAccountID)
}

func RequireActiveProfileIntWithOverride(key config.FieldKey, override int) int {
	v := GetProfileIntWithOverride(GetActiveProfileName(), key, override)
	if v == 0 {
		log.Fatalf("%s is required", key)
	}

	return v
}

func RequireActiveProfileInt(key config.FieldKey) int {
	v := GetProfileInt(GetActiveProfileName(), key)
	if v == 0 {
		log.Fatalf("%s is required", key)
	}

	return v
}

func GetActiveProfileInt(key config.FieldKey) int {
	return GetProfileInt(GetActiveProfileName(), key)
}

func GetActiveProfileIntWithOverride(key config.FieldKey, override int) int {
	return GetProfileIntWithOverride(GetActiveProfileName(), key, override)
}

func GetProfileInt(profileName string, key config.FieldKey) int {
	v, err := config.CredentialsProvider.GetIntWithScope(profileName, key)
	if err != nil {
		return 0
	}

	return int(v)
}

func GetProfileIntWithOverride(profileName string, key config.FieldKey, override int) int {
	o := int64(override)
	v, err := config.CredentialsProvider.GetIntWithScopeAndOverride(profileName, key, &o)
	if err != nil {
		return 0
	}

	return int(v)
}

func GetConfigString(key config.FieldKey) string {
	return GetConfigStringWithOverride(key, "")
}

func GetConfigStringWithOverride(key config.FieldKey, override string) string {
	v, err := config.ConfigStore.GetStringWithOverride(key, &override)
	if err != nil {
		return ""
	}

	return v
}

func GetConfigTernary(key config.FieldKey) config.Ternary {
	v, err := config.ConfigStore.GetTernary(key)
	if err != nil {
		return config.Ternary("")
	}

	return v
}

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

func SetDefaultProfile(profileName string) error {
	defaultProfileFilePath := filepath.Join(config.BasePath, config.DefaultProfileFileName)
	return ioutil.WriteFile(defaultProfileFilePath, []byte("\""+profileName+"\""), 0644)
}

func RemoveDefaultProfile() error {
	defaultProfileFilePath := filepath.Join(config.BasePath, config.DefaultProfileFileName)
	return os.Remove(defaultProfileFilePath)
}

func GetProfileFieldDefinition(key config.FieldKey) *config.FieldDefinition {
	return config.CredentialsProvider.GetFieldDefinition(key)
}

func GetConfigFieldDefinition(key config.FieldKey) *config.FieldDefinition {
	return config.ConfigStore.GetFieldDefinition(key)
}

func VisitAllProfileFields(profileName string, fn func(d config.FieldDefinition)) {
	config.CredentialsProvider.VisitAllFieldsWithScope(profileName, fn)
}

func VisitAllConfigFields(fn func(d config.FieldDefinition)) {
	config.ConfigStore.VisitAllFields(fn)
}

func GetValidFieldKeys() (fieldKeys []config.FieldKey) {
	config.ConfigStore.VisitAllFields(func(fd config.FieldDefinition) {
		fieldKeys = append(fieldKeys, fd.Key)
	})

	return fieldKeys
}

func GetProfileNames() []string {
	return config.CredentialsProvider.GetScopes()
}

func RemoveProfile(profileName string) error {
	if err := config.CredentialsProvider.RemoveScope(profileName); err != nil {
		return err
	}

	// Set a new default profile, or delete if there are no others
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
		} else {
			if err := RemoveDefaultProfile(); err != nil {
				return fmt.Errorf("could not delete default profile")
			}
		}
	}

	return nil
}

func SetConfigString(key config.FieldKey, value string) error {
	return SetConfigValue(key, value)
}

func SetConfigValue(key config.FieldKey, value interface{}) error {
	return config.ConfigStore.Set(key, value)
}

func DeleteConfigValue(key config.FieldKey, value interface{}) error {
	return config.ConfigStore.DeleteKey(key)
}

func SetActiveProfileString(key config.FieldKey, value string) error {
	return SetProfileString(GetActiveProfileName(), key, value)
}

func SetProfileString(profileName string, key config.FieldKey, value string) error {
	return config.CredentialsProvider.SetWithScope(profileName, key, value)
}

func SetActiveProfileInt(key config.FieldKey, value int) error {
	return SetProfileInt(GetActiveProfileName(), key, value)
}

func SetProfileInt(profileName string, key config.FieldKey, value int) error {
	return config.CredentialsProvider.SetWithScope(profileName, key, value)
}
