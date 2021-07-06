package configuration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-client-go/pkg/region"
)

const (
	APIKey             ConfigKey = "apiKey"
	InsightsInsertKey  ConfigKey = "insightsInsertKey"
	Region             ConfigKey = "region"
	AccountID          ConfigKey = "accountID"
	LicenseKey         ConfigKey = "licenseKey"
	LogLevel           ConfigKey = "loglevel"
	PluginDir          ConfigKey = "plugindir"
	PreReleaseFeatures ConfigKey = "prereleasefeatures"
	SendUsageData      ConfigKey = "sendusagedata"

	configFileName         = "config.json"
	credentialsFileName    = "credentials.json"
	defaultProfileFileName = "default-profile.json"
	pluginDir              = "plugins"
	activeProfileName      = "default"
)

var (
	configProvider      *ConfigProvider
	credentialsProvider *ConfigProvider
	BasePath            string = configBasePath()
)

func init() {
	initializeConfigProvider()
	initializeCredentialsProvider()
}

func initializeCredentialsProvider() {
	p, err := NewConfigProvider(
		WithFilePersistence(filepath.Join(BasePath, credentialsFileName)),
		WithFieldDefinitions(
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

func initializeConfigProvider() {
	p, err := NewConfigProvider(
		WithFilePersistence(filepath.Join(BasePath, configFileName)),
		WithScope("*"),
		WithFieldDefinitions(
			FieldDefinition{
				Key:               LogLevel,
				EnvVar:            "NEW_RELIC_CLI_LOG_LEVEL",
				Default:           "info",
				SetValidationFunc: StringInStrings(false, "Info", "Debug", "Trace", "Warn", "Error"),
			},
			FieldDefinition{
				Key:     PluginDir,
				EnvVar:  "NEW_RELIC_CLI_PLUGIN_DIR",
				Default: filepath.Join(configBasePath(), pluginDir),
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

	configProvider = p
}

func GetActiveProfileName() string {
	return activeProfileName
}

func GetActiveProfileString(key ConfigKey) string {
	return GetProfileString(GetActiveProfileName(), key)
}

func GetProfileString(profileName string, key ConfigKey) string {
	v, err := credentialsProvider.GetStringWithScope(GetActiveProfileName(), key)
	if err != nil {
		return ""
	}

	return v
}

func GetActiveProfileInt(key ConfigKey) int {
	return GetProfileInt(GetActiveProfileName(), key)
}

func GetProfileInt(profileName string, key ConfigKey) int {
	v, err := credentialsProvider.GetIntWithScope(GetActiveProfileName(), key)
	if err != nil {
		return 0
	}

	return int(v)
}

func GetConfigString(key ConfigKey) string {
	v, err := configProvider.GetString(key)
	if err != nil {
		log.Fatalf("could not load value %s from config: %s", key, err)
	}

	return v
}

func GetConfigTernary(key ConfigKey) Ternary {
	v, err := configProvider.GetTernary(key)
	if err != nil {
		log.Fatalf("could not load value %s from config: %s", key, err)
	}

	return v
}

func GetDefaultProfile() (string, error) {
	defaultProfileFilePath := filepath.Join(BasePath, defaultProfileFileName)
	data, err := ioutil.ReadFile(defaultProfileFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	return string(data), nil
}

func SetDefaultProfile(profileName string) error {
	defaultProfileFilePath := filepath.Join(BasePath, defaultProfileFileName)
	return ioutil.WriteFile(defaultProfileFilePath, []byte("\""+profileName+"\""), 0644)
}

func GetProfileFieldDefinition(key ConfigKey) *FieldDefinition {
	return credentialsProvider.getFieldDefinition(key)
}

func GetConfigFieldDefinition(key ConfigKey) *FieldDefinition {
	return configProvider.getFieldDefinition(key)
}

func VisitAllProfileFields(profileName string, fn func(d FieldDefinition)) {
	credentialsProvider.VisitAllFieldsWithScope(profileName, fn)
}

func VisitAllConfigFields(fn func(d FieldDefinition)) {
	configProvider.VisitAllFields(fn)
}

func GetValidConfigKeys() (configKeys []ConfigKey) {
	configProvider.VisitAllFields(func(fd FieldDefinition) {
		configKeys = append(configKeys, fd.Key)
	})

	return configKeys
}

func GetProfileNames() []string {
	return credentialsProvider.GetScopes()
}

func RemoveProfile(profileName string) error {
	return credentialsProvider.RemoveScope(profileName)
}

func SetConfigString(key ConfigKey, value string) error {
	return SetConfigValue(key, value)
}

func SetConfigValue(key ConfigKey, value interface{}) error {
	return configProvider.Set(key, value)
}

func SetActiveProfileString(key ConfigKey, value string) error {
	return SetProfileString(GetActiveProfileName(), key, value)
}

func SetProfileString(profileName string, key ConfigKey, value string) error {
	return credentialsProvider.SetWithScope(profileName, key, value)
}

func SetActiveProfileInt(key ConfigKey, value int) error {
	return SetProfileInt(GetActiveProfileName(), key, value)
}

func SetProfileInt(profileName string, key ConfigKey, value int) error {
	return credentialsProvider.SetWithScope(profileName, key, value)
}

func configBasePath() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("cannot locate user's home directory: %s", err)
	}

	return fmt.Sprintf("%s/.newrelic", home)
}
