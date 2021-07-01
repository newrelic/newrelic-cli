package configuration

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-client-go/pkg/region"
	log "github.com/sirupsen/logrus"
)

const (
	configFileName      = "config.json"
	credentialsFileName = "credentials.json"
	pluginDir           = "plugins"
	activeProfileName   = "default"
)

var (
	configProvider      *ConfigProvider
	credentialsProvider *ConfigProvider
)

func init() {
	basePath := configBasePath()
	initializeConfigProvider(basePath)
	initializeCredentialsProvider(basePath)
}

func initializeCredentialsProvider(basePath string) {
	p, err := NewConfigProvider(
		WithFilePersistence(filepath.Join(basePath, credentialsFileName)),
		WithFieldDefinitions(
			FieldDefinition{
				Key:    "apiKey",
				EnvVar: "NEW_RELIC_API_KEY",
			},
			FieldDefinition{
				Key:    "insightsInsertKey",
				EnvVar: "NEW_RELIC_INSIGHTS_INSERT_KEY",
			},
			FieldDefinition{
				Key:            "region",
				EnvVar:         "NEW_RELIC_REGION",
				ValidationFunc: StringInStrings(false, region.Staging.String(), region.US.String(), region.EU.String()),
			},
			FieldDefinition{
				Key:            "accountID",
				EnvVar:         "NEW_RELIC_ACCOUNT_ID",
				ValidationFunc: IntGreaterThan(0),
			},
			FieldDefinition{
				Key:    "licenseKey",
				EnvVar: "NEW_RELIC_LICENSE_KEY",
			},
		),
	)

	if err != nil {
		log.Fatalf("could not create credentials provider: %s", err)
	}

	credentialsProvider = p
}

func initializeConfigProvider(basePath string) {
	p, err := NewConfigProvider(
		WithFilePersistence(filepath.Join(basePath, configFileName)),
		WithFieldDefinitions(
			FieldDefinition{
				Key:     "loglevel",
				EnvVar:  "NEW_RELIC_CLI_LOG_LEVEL",
				Default: "debug",
			},
			FieldDefinition{
				Key:     "plugindir",
				EnvVar:  "NEW_RELIC_CLI_PLUGIN_DIR",
				Default: filepath.Join(configBasePath(), pluginDir),
			},
			FieldDefinition{
				Key:     "prereleasefeatures",
				EnvVar:  "NEW_RELIC_CLI_PRERELEASEFEATURES",
				Default: config.TernaryValues.Unknown,
			},
			FieldDefinition{
				Key:     "sendusagedata",
				EnvVar:  "NEW_RELIC_CLI_SENDUSAGEDATA",
				Default: config.TernaryValues.Unknown,
			},
		),
		WithScope("*"),
	)

	if err != nil {
		log.Fatalf("could not create configuration provider: %s", err)
	}

	configProvider = p
}

func GetActiveProfileString(key string) string {
	v, err := credentialsProvider.GetStringWithScope(activeProfileName, key)
	if err != nil {
		log.Fatalf("could not load value %s from active profile %s: %s", key, activeProfileName, err)
	}

	return v

}

func GetActiveProfileInt(key string) int64 {
	v, err := credentialsProvider.GetIntWithScope(activeProfileName, key)
	if err != nil {
		log.Fatalf("could not load value %s from active profile %s: %s", key, activeProfileName, err)
	}

	return v

}

func GetConfigString(key string) string {
	v, err := configProvider.GetString(key)
	if err != nil {
		log.Fatalf("could not load value %s from config: %s", key, err)
	}

	return v

}

func configBasePath() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("cannot locate user's home directory: %s", err)
	}

	return fmt.Sprintf("%s/.newrelic", home)
}
