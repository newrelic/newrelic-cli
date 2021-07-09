package config

import (
	"fmt"
	"path/filepath"

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

	DefaultProfileFileName = "default-profile.json"
	ConfigFileName         = "config.json"
	CredentialsFileName    = "credentials.json"
	DefaultPluginDir       = "plugins"
)

var (
	ConfigStore         *Store
	CredentialsProvider *Store
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
	InitializeConfigStore()
	InitializeCredentialsStore()
}

func InitializeCredentialsStore() {
	p, err := NewStore(
		PersistToFile(filepath.Join(BasePath, CredentialsFileName)),
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

	CredentialsProvider = p
}

func InitializeConfigStore() {
	p, err := NewStore(
		PersistToFile(filepath.Join(BasePath, ConfigFileName)),
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
				Default: filepath.Join(BasePath, DefaultPluginDir),
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

	ConfigStore = p
}

func configBasePath() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("cannot locate user's home directory: %s", err)
	}

	return fmt.Sprintf("%s/.newrelic", home)
}
