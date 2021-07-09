package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-client-go/pkg/region"
)

const (
	APIKey             FieldKey = "apiKey"
	InsightsInsertKey  FieldKey = "insightsInsertKey"
	Region             FieldKey = "region"
	AccountID          FieldKey = "accountID"
	LicenseKey         FieldKey = "licenseKey"
	LogLevel           FieldKey = "loglevel"
	PluginDir          FieldKey = "plugindir"
	PreReleaseFeatures FieldKey = "prereleasefeatures"
	SendUsageData      FieldKey = "sendusagedata"

	DefaultProfileName = "default"

	defaultProfileFileName = "default-profile.json"
	configFileName         = "config.json"
	credentialsFileName    = "credentials.json"
	pluginDir              = "plugins"
)

var (
	store               *Store
	credentialsProvider *Store
	BasePath            string = configBasePath()

	FlagProfileName string
	FlagDebug       bool
	FlagTrace       bool
	FlagAccountID   int
)

func init() {
	Init(configBasePath())
}

func Init(basePath string) {
	BasePath = basePath
	initializeStore()
	initializeCredentialsProvider()
}

func initializeCredentialsProvider() {
	p, err := NewStore(
		PersistToFile(filepath.Join(BasePath, credentialsFileName)),
		EnforceStrictFields(),
		ConfigureFields(
			FieldDefinition{
				Key:       APIKey,
				EnvVar:    "NEW_RELIC_API_KEY",
				Sensitive: true,
			},
			FieldDefinition{
				Key:       InsightsInsertKey,
				EnvVar:    "NEW_RELIC_INSIGHTS_INSERT_KEY",
				Sensitive: true,
			},
			FieldDefinition{
				Key:    Region,
				EnvVar: "NEW_RELIC_REGION",
				SetValidationFunc: StringInStrings(false,
					region.Staging.String(),
					region.US.String(),
					region.EU.String(),
				),
				Default: region.US.String(),
			},
			FieldDefinition{
				Key:               AccountID,
				EnvVar:            "NEW_RELIC_ACCOUNT_ID",
				SetValidationFunc: IntGreaterThan(0),
			},
			FieldDefinition{
				Key:       LicenseKey,
				EnvVar:    "NEW_RELIC_LICENSE_KEY",
				Sensitive: true,
			},
		),
	)

	if err != nil {
		log.Fatalf("could not create credentials provider: %s", err)
	}

	credentialsProvider = p
}

func initializeStore() {
	p, err := NewStore(
		PersistToFile(filepath.Join(BasePath, configFileName)),
		UseGlobalScope("*"),
		EnforceStrictFields(),
		ConfigureFields(
			FieldDefinition{
				Key:               LogLevel,
				EnvVar:            "NEW_RELIC_CLI_LOG_LEVEL",
				Default:           "info",
				SetValidationFunc: StringInStrings(false, "Info", "Debug", "Trace", "Warn", "Error"),
			},
			FieldDefinition{
				Key:     PluginDir,
				EnvVar:  "NEW_RELIC_CLI_PLUGIN_DIR",
				Default: filepath.Join(BasePath, pluginDir),
			},
			FieldDefinition{
				Key:               PreReleaseFeatures,
				EnvVar:            "NEW_RELIC_CLI_PRERELEASEFEATURES",
				SetValidationFunc: IsTernary(),
				Default:           TernaryValues.Unknown,
			},
			FieldDefinition{
				Key:               SendUsageData,
				EnvVar:            "NEW_RELIC_CLI_SENDUSAGEDATA",
				SetValidationFunc: IsTernary(),
				Default:           TernaryValues.Unknown,
			},
		),
	)

	if err != nil {
		log.Fatalf("could not create configuration provider: %s", err)
	}

	store = p
}

func GetActiveProfileName() string {
	if FlagProfileName != "" {
		return FlagProfileName
	}

	profileName, err := GetDefaultProfileName()
	if err != nil || profileName == "" {
		return DefaultProfileName
	}

	return profileName
}

func GetActiveProfileString(key FieldKey) string {
	return GetProfileString(GetActiveProfileName(), key)
}

func RequireActiveProfileString(key FieldKey) string {
	v := GetProfileString(GetActiveProfileName(), key)
	if v == "" {
		log.Fatalf("%s is required", key)
	}

	return v
}

func GetActiveProfileValue(profileName string, key FieldKey) interface{} {
	return GetProfileValue("", key)
}
func GetProfileValue(profileName string, key FieldKey) interface{} {
	v, err := credentialsProvider.GetWithScope(profileName, key)
	if err != nil {
		return nil
	}

	return v
}

func GetProfileString(profileName string, key FieldKey) string {
	v, err := credentialsProvider.GetStringWithScope(profileName, key)
	if err != nil {
		return ""
	}

	return v
}

func GetLogLevelWithFlagOverride() string {
	var override string
	if FlagDebug {
		override = "debug"
	}

	if FlagTrace {
		override = "trace"
	}

	return GetConfigStringWithOverride(LogLevel, override)
}

func RequireActiveProfileAccountIDWithFlagOverride() int {
	v := GetActiveProfileAccountIDWithFlagOverride()
	if v == 0 {
		log.Fatalf("%s is required", AccountID)
	}

	return v
}

func GetActiveProfileAccountIDWithFlagOverride() int {
	return GetActiveProfileIntWithOverride(AccountID, FlagAccountID)
}

func RequireActiveProfileIntWithOverride(key FieldKey, override int) int {
	v := GetProfileIntWithOverride(GetActiveProfileName(), key, override)
	if v == 0 {
		log.Fatalf("%s is required", key)
	}

	return v
}

func RequireActiveProfileInt(key FieldKey) int {
	v := GetProfileInt(GetActiveProfileName(), key)
	if v == 0 {
		log.Fatalf("%s is required", key)
	}

	return v
}

func GetActiveProfileInt(key FieldKey) int {
	return GetProfileInt(GetActiveProfileName(), key)
}

func GetActiveProfileIntWithOverride(key FieldKey, override int) int {
	return GetProfileIntWithOverride(GetActiveProfileName(), key, override)
}

func GetProfileInt(profileName string, key FieldKey) int {
	v, err := credentialsProvider.GetIntWithScope(profileName, key)
	if err != nil {
		return 0
	}

	return int(v)
}

func GetProfileIntWithOverride(profileName string, key FieldKey, override int) int {
	o := int64(override)
	v, err := credentialsProvider.GetIntWithScopeAndOverride(profileName, key, &o)
	if err != nil {
		return 0
	}

	return int(v)
}

func GetConfigString(key FieldKey) string {
	return GetConfigStringWithOverride(key, "")
}

func GetConfigStringWithOverride(key FieldKey, override string) string {
	v, err := store.GetStringWithOverride(key, &override)
	if err != nil {
		return ""
	}

	return v
}

func GetConfigTernary(key FieldKey) Ternary {
	v, err := store.GetTernary(key)
	if err != nil {
		return Ternary("")
	}

	return v
}

func GetDefaultProfileName() (string, error) {
	defaultProfileFilePath := filepath.Join(BasePath, defaultProfileFileName)
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
	defaultProfileFilePath := filepath.Join(BasePath, defaultProfileFileName)
	return ioutil.WriteFile(defaultProfileFilePath, []byte("\""+profileName+"\""), 0644)
}

func RemoveDefaultProfile() error {
	defaultProfileFilePath := filepath.Join(BasePath, defaultProfileFileName)
	return os.Remove(defaultProfileFilePath)
}

func GetProfileFieldDefinition(key FieldKey) *FieldDefinition {
	return credentialsProvider.getFieldDefinition(key)
}

func GetConfigFieldDefinition(key FieldKey) *FieldDefinition {
	return store.getFieldDefinition(key)
}

func VisitAllProfileFields(profileName string, fn func(d FieldDefinition)) {
	credentialsProvider.VisitAllFieldsWithScope(profileName, fn)
}

func VisitAllConfigFields(fn func(d FieldDefinition)) {
	store.VisitAllFields(fn)
}

func GetValidFieldKeys() (fieldKeys []FieldKey) {
	store.VisitAllFields(func(fd FieldDefinition) {
		fieldKeys = append(fieldKeys, fd.Key)
	})

	return fieldKeys
}

func GetProfileNames() []string {
	return credentialsProvider.GetScopes()
}

func RemoveProfile(profileName string) error {
	if err := credentialsProvider.RemoveScope(profileName); err != nil {
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

func SetConfigString(key FieldKey, value string) error {
	return SetConfigValue(key, value)
}

func SetConfigValue(key FieldKey, value interface{}) error {
	return store.Set(key, value)
}

func DeleteConfigValue(key FieldKey, value interface{}) error {
	return store.DeleteKey(key)
}

func SetActiveProfileString(key FieldKey, value string) error {
	return SetProfileString(GetActiveProfileName(), key, value)
}

func SetProfileString(profileName string, key FieldKey, value string) error {
	return credentialsProvider.SetWithScope(profileName, key, value)
}

func SetActiveProfileInt(key FieldKey, value int) error {
	return SetProfileInt(GetActiveProfileName(), key, value)
}

func SetProfileInt(profileName string, key FieldKey, value int) error {
	return credentialsProvider.SetWithScope(profileName, key, value)
}

func configBasePath() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("cannot locate user's home directory: %s", err)
	}

	return fmt.Sprintf("%s/.newrelic", home)
}
