package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// DefaultConfigDirectory is the default location for the CLI config files
const DefaultConfigDirectory = "$HOME/.newrelic"

// DefaultPluginDirectory is the default sub-directory containing the plugings
const DefaultPluginDirectory = DefaultConfigDirectory + "/plugins"

// DefaultConfigName is the default name of the global configuration file
const DefaultConfigName = "config"

// DefaultEnvPrefix is used when reading environment variables
const DefaultEnvPrefix = "newrelic"

// Config contains the main CLI configuration
type Config struct {
	LogLevel      string `mapstructure:"loglevel"`      // LogLevel for verbose output
	PluginDir     string `mapstructure:"plugindir"`     // PluginDir is the directory where plugins will be installed
	SendUsageData string `mapstructure:"sendusagedata"` // SendUsageData enables sending usage statistics to New Relic
	ProfileName   string // ProfileName is the configured profile to use
}

// Load initializes the cli configuration
func Load(cfgFile string, logLevel string) (*Config, error) {
	if logLevel != "" {
		lvl, err := log.ParseLevel(logLevel)
		if err != nil {
			return nil, err
		}

		log.SetLevel(lvl)
	}

	log.Debug("loading config file")

	cfgViper := viper.New()
	cfgViper.SetEnvPrefix(DefaultEnvPrefix)
	cfgViper.SetConfigName(DefaultConfigName)
	cfgViper.AddConfigPath(DefaultConfigDirectory) // adding home directory as first search path
	cfgViper.AddConfigPath(".")                    // current directory to search path
	cfgViper.AutomaticEnv()                        // read in environment variables that match

	if cfgFile != "" {
		cfgViper.SetConfigFile(cfgFile)
	}

	// Read in config
	err := cfgViper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug("no config file found, using defaults")
			cfg := Config{}
			cfg.Validate()
			return &cfg, nil
		} else if e, ok := err.(viper.ConfigParseError); ok {
			return nil, fmt.Errorf("error parsing config file: %v", e)
		}
	}

	log.Debugf("loaded config from: %v", cfgViper.ConfigFileUsed())

	// For legacy reasons the config has a scope level, default scope is '*'
	cfgMap := map[string]Config{}
	err = cfgViper.Unmarshal(&cfgMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config with error: %v", err)
	}
	config, ok := cfgMap["*"]
	if !ok {
		return nil, fmt.Errorf("failed to locate global config")
	}

	if config.LogLevel == "" {
		config.LogLevel = logLevel
	}

	err = config.Validate()

	return &config, err
}

// Validate the configuration, set defaults if needed
func (c *Config) Validate() error {
	log.Debug("validating config")

	if c == nil {
		return nil
	}

	if c.LogLevel == "" {
		c.LogLevel = "info"
	}

	switch c.SendUsageData {
	case "ALLOW", "DISALLOW", "NOT_ASKED":
		break
	default:
		c.SendUsageData = "NOT_ASKED"
	}

	if c.PluginDir == "" {
		c.PluginDir = DefaultPluginDirectory
	} else {
		// TODO: Validate the dir exists
	}

	return nil
}
