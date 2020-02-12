package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// DefaultConfigDirectory is the default location for the CLI config files
const DefaultConfigDirectory = "$HOME/.newrelic"

// DefaultPluginDirectory is the default sub-directory containing the plugings
const DefaultPluginDirectory = DefaultConfigDirectory + "/plugins"

// DefaultConfigName is the default name of the global configuration file
const DefaultConfigName = "config"

// DefaultConfigType to read, though any file type supported by viper is allowed
const DefaultConfigType = "json"

// DefaultEnvPrefix is used when reading environment variables
const DefaultEnvPrefix = "newrelic"

// Config contains the main CLI configuration
type Config struct {
	LogLevel      string `mapstructure:"logLevel"`      // LogLevel for verbose output
	PluginDir     string `mapstructure:"pluginDir"`     // PluginDir is the directory where plugins will be installed
	SendUsageData string `mapstructure:"sendUsageData"` // SendUsageData enables sending usage statistics to New Relic
	ProfileName   string `mapstructure:"profileName"`   // ProfileName is the configured profile to use
}

// DefaultConfig represents the configuration default values.
var DefaultConfig = Config{
	LogLevel:      "DEBUG",
	PluginDir:     DefaultPluginDirectory,
	SendUsageData: "NOT ASKED",
	ProfileName:   "",
}

// LoadConfig loads the configuration
func LoadConfig(configFile string, logLevel string) (*Config, error) {
	config, err := Load(configFile, logLevel)
	if err != nil {
		return &Config{}, err
	}

	return config, nil
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
	cfgViper.SetConfigType(DefaultConfigType)
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

// List outputs a list of all the configuration values
func (c *Config) List() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Key", "Value", "Origin"})
	t.AppendRows(c.toRows())
	t.SetStyle(table.Style{
		Name: "nr-cli-table",
		Box: table.BoxStyle{
			MiddleHorizontal: "-",
			MiddleSeparator:  " ",
			MiddleVertical:   " ",
		},
		Color: table.ColorOptions{
			Header: text.Colors{text.Bold},
		},
		Options: table.Options{
			DrawBorder:      false,
			SeparateColumns: true,
			SeparateHeader:  true,
		},
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:   "Value",
			Colors: text.Colors{text.FgHiCyan},
		},
		{
			Name:   "Origin",
			Colors: text.Colors{text.FgHiBlack},
		},
	})

	t.Render()

	bold := color.New(color.Bold).SprintFunc()
	fmt.Printf("\n\nRun %s for more info.\n", bold("\"newrelic config get --key KEY\""))
}

func (c *Config) toRows() []table.Row {
	t := reflect.TypeOf(*c)
	v := reflect.ValueOf(*c)
	d := reflect.ValueOf(DefaultConfig)

	o := make([]table.Row, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if name, ok := field.Tag.Lookup("mapstructure"); ok {
			value := v.Field(i).Interface().(string)
			defaultValue := d.Field(i).Interface().(string)

			origin := "Default"
			if defaultValue != value {
				origin = "User config"
			}

			out := table.Row{name, value, origin}
			o = append(o, out)
		}
	}

	return o
}
