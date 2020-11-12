package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// DefaultConfigName is the default name of the global configuration file
	DefaultConfigName = "config"

	// DefaultConfigType to read, though any file type supported by viper is allowed
	DefaultConfigType = "json"

	// DefaultEnvPrefix is used when reading environment variables
	DefaultEnvPrefix = "NEW_RELIC_CLI"

	// DefaultLogLevel is the default log level
	DefaultLogLevel = "INFO"

	globalScopeIdentifier = "*"
)

var (
	// Should we export configuration?
	configuration          *Configuration
	once                   sync.Once
	DefaultConfigDirectory string // Not sure if we need this
)

type Configuration struct {
	LogLevel      string
	Profiles      map[string]Profile
	ViperInstance *viper.Viper

	// Used with `Get` and `Set`
	configFilePath string
}

type Profile struct {
	APIKey            string `json:"apiKey,omitempty"`            // For accessing New Relic GraphQL resources
	InsightsInsertKey string `json:"insightsInsertKey,omitempty"` // For posting custom events
	Region            string `json:"region,omitempty"`            // Region to use for New Relic resources
	AccountID         int    `json:"accountID"`                   // AccountID to use for New Relic resources
	LicenseKey        string `json:"licenseKey"`                  // License key to use for agent config and ingest
}

func Configure(configDirectory string) (*Configuration, error) {
	configDir, err := getDefaultConfigDirectory()
	if err != nil {
		log.Fatal(err.Error())
	}

	if configDirectory != "" {
		configDir = configDirectory
	}

	configFilePath := fmt.Sprintf("%s/%s.%s", configDir, DefaultConfigName, DefaultConfigType)

	// Initialize the Viper config. The Viper instance can be considered
	// the central source  of truth for configuration settings.
	cfgViper := initViperConfig(configDir)

	// Create a singleton instance of the configuration
	once.Do(func() {
		logLevel := cfgViper.Get(keyGlobalScope("loglevel")).(string)

		configuration = new(Configuration)
		configuration.ViperInstance = cfgViper
		configuration.configFilePath = configFilePath

		configuration.setLogger(logLevel)
	})

	if err = createConfigFile(configDir, configFilePath, cfgViper); err != nil {
		return nil, err
	}

	// TODO: Attempt to add a default profile.
	// TODO: Attempt to add default profile credentials as well.

	return configuration, nil
}

func initViperConfig(configDir string) *viper.Viper {
	cfgViper := viper.New()

	cfgViper.Set(globalScopeIdentifier, map[string]interface{}{
		"loglevel":           DefaultLogLevel,
		"sendusagedata":      TernaryValues.Unknown,
		"prereleasefeatures": TernaryValues.Unknown,
	})

	// Mapping our config file and env to Viper
	cfgViper.SetEnvPrefix(DefaultEnvPrefix)
	cfgViper.SetConfigName(DefaultConfigName)
	cfgViper.SetConfigType(DefaultConfigType)

	// Set the config file path
	cfgViper.AddConfigPath(configDir)

	// Read environment variables that
	// match the environment prefix
	cfgViper.AutomaticEnv()

	return cfgViper
}

func hasConfigFile(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func createConfigFile(configDir string, path string, cfgViper *viper.Viper) error {
	if !hasConfigFile(path) {
		if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
			return err
		}

		if err := cfgViper.WriteConfigAs(path); err != nil {
			return err
		}
	} else {
		if err := cfgViper.WriteConfig(); err != nil {
			return err
		}
	}

	return nil
}

// This was just a convience thing for testing. Do we need this?
func (c *Configuration) GetConfigFilePath() string {
	return c.configFilePath
}

// Just an idea for now
func (c *Configuration) GetCurrentProfile() *Profile {
	return nil
}

func (c *Configuration) GetDefaultProfile() *Profile {
	return nil
}

func (c *Configuration) GetProfile(name string) *Profile {
	return nil
}

// Get returns the value for a given key at the global scope level.
// TODO: Support other scopes?
func (c *Configuration) Get(key string) interface{} {
	return c.ViperInstance.Get(keyGlobalScope(key))
}

// Set sets a specified key with the provided value at the global scope level.
// TODO: Support other scopes?
func (c *Configuration) Set(key string, value interface{}) {
	c.ViperInstance.Set(keyGlobalScope(key), value)
}

func (c *Configuration) setLogger(level string) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        time.RFC3339,
		DisableLevelTruncation: true,
	})

	c.setLogLevel(level)
}

// TODO: Come up with a better/safer way to handle this
// And does this "need" to be a method on the Configuration struct?
func (c *Configuration) setLogLevel(level string) {
	logLevelMap := map[string]logrus.Level{
		"info":  log.InfoLevel, // default
		"debug": log.DebugLevel,
		"warn":  log.WarnLevel,
		"error": log.ErrorLevel,
		"trace": log.TraceLevel,
	}

	if _, ok := logLevelMap[strings.ToLower(level)]; !ok {
		log.Warnf("error setting log level '%s', using default %s", level, DefaultLogLevel)
		// Set to default 'info'
		log.SetLevel(logLevelMap[strings.ToLower(DefaultLogLevel)])
	} else {
		log.SetLevel(logLevelMap[strings.ToLower(DefaultLogLevel)])
	}
}

// Set default configuration directory on init
func init() {
	cfgDir, err := getDefaultConfigDirectory()
	if err != nil {
		log.Fatalf("error building default config directory: %s", err)
	}

	DefaultConfigDirectory = cfgDir
}
