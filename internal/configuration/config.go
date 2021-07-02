package configuration

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
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
	basePath            string = configBasePath()
)

func init() {
	initializeConfigProvider()
	initializeCredentialsProvider()
}

func initializeCredentialsProvider() {
	p, err := NewConfigProvider(
		WithFilePersistence(filepath.Join(basePath, credentialsFileName)),
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
				ValidationFunc: StringInStrings(false,
					region.Staging.String(),
					region.US.String(),
					region.EU.String(),
				),
			},
			FieldDefinition{
				Key:            AccountID,
				EnvVar:         "NEW_RELIC_ACCOUNT_ID",
				ValidationFunc: IntGreaterThan(0),
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
		WithFilePersistence(filepath.Join(basePath, configFileName)),
		WithScope("*"),
		WithFieldDefinitions(
			FieldDefinition{
				Key:     LogLevel,
				EnvVar:  "NEW_RELIC_CLI_LOG_LEVEL",
				Default: "debug",
			},
			FieldDefinition{
				Key:     PluginDir,
				EnvVar:  "NEW_RELIC_CLI_PLUGIN_DIR",
				Default: filepath.Join(configBasePath(), pluginDir),
			},
			FieldDefinition{
				Key:     PreReleaseFeatures,
				EnvVar:  "NEW_RELIC_CLI_PRERELEASEFEATURES",
				Default: config.TernaryValues.Unknown,
			},
			FieldDefinition{
				Key:     SendUsageData,
				EnvVar:  "NEW_RELIC_CLI_SENDUSAGEDATA",
				Default: config.TernaryValues.Unknown,
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
	return GetProfileString(activeProfileName, key)
}

func GetProfileString(profileName string, key ConfigKey) string {
	v, err := credentialsProvider.GetStringWithScope(activeProfileName, key)
	if err != nil {
		log.Fatalf("could not load value %s from active profile %s: %s", key, profileName, err)
	}

	return v
}

func GetActiveProfileInt(key ConfigKey) int64 {
	v, err := credentialsProvider.GetIntWithScope(activeProfileName, key)
	if err != nil {
		log.Fatalf("could not load value %s from active profile %s: %s", key, activeProfileName, err)
	}

	return v
}

func GetConfigString(key ConfigKey) string {
	v, err := configProvider.GetString(key)
	if err != nil {
		log.Fatalf("could not load value %s from config: %s", key, err)
	}

	return v
}

func SetDefaultProfile(profileName string) error {
	defaultProfileFilePath := filepath.Join(basePath, defaultProfileFileName)
	return ioutil.WriteFile(defaultProfileFilePath, []byte("\""+profileName+"\""), 0644)
}

func VisitAllProfileFields(profileName string, fn func(d FieldDefinition)) {
	credentialsProvider.VisitAllFieldsWithScope(profileName, fn)
}

func GetProfileNames() []string {
	return credentialsProvider.GetScopes()
}

func RemoveProfile(profileName string) error {
	return credentialsProvider.RemoveScope(profileName)
}

func configBasePath() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("cannot locate user's home directory: %s", err)
	}

	return fmt.Sprintf("%s/.newrelic", home)
}
