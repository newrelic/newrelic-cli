package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const defaultConfigDirectory = "$HOME/.newrelic"
const defaultPluginDirectory = defaultConfigDirectory + "/plugins"
const defaultConfigName = "config"
const defaultEnvPrefix = "newrelic"

// Config contains the main CLI configuration
type Config struct {
	SendUsageData string `mapstructure:"sendusagedata"` // SendUsageData enables sending usage statistics to New Relic
	PluginDir     string `mapstructure:"plugindir"`     // PluginDir is the directory where plugins will be installed
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

	cfgViper := viper.New()
	cfgViper.SetEnvPrefix(defaultEnvPrefix)
	cfgViper.SetConfigName(defaultConfigName)
	cfgViper.AddConfigPath(defaultConfigDirectory) // adding home directory as first search path
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
	} else {
		log.Debugf("using config file: %v", cfgViper.ConfigFileUsed())
	}

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

	err = config.Validate()

	return &config, err
}

// Validate the configuration, set defaults if needed
func (c *Config) Validate() error {
	if c == nil {
		return nil
	}

	switch c.SendUsageData {
	case "ALLOW", "DISALLOW", "NOT_ASKED":
		break
	default:
		c.SendUsageData = "NOT_ASKED"
	}

	if c.PluginDir == "" {
		c.PluginDir = defaultPluginDirectory
	} else {
		// TODO: Validate the dir exists
	}

	return nil
}
