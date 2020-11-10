package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/newrelic/newrelic-cli/internal/utils"

	"github.com/imdario/mergo"
	"github.com/jedib0t/go-pretty/v6/text"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/newrelic/newrelic-cli/internal/output"
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
	// DefaultConfigDirectory is the default location for the CLI config files
	DefaultConfigDirectory string

	defaultConfig *Config
)

// Config contains the main CLI configuration
type Config struct {
	LogLevel           string  `mapstructure:"logLevel"`           // LogLevel for verbose output
	PluginDir          string  `mapstructure:"pluginDir"`          // PluginDir is the directory where plugins will be installed
	SendUsageData      Ternary `mapstructure:"sendUsageData"`      // SendUsageData enables sending usage statistics to New Relic
	PreReleaseFeatures Ternary `mapstructure:"preReleaseFeatures"` // PreReleaseFeatures enables display on features within the CLI that are announced but not generally available to customers

	configDir string
}

// Value represents an instance of a configuration field.
type Value struct {
	Name    string
	Value   interface{}
	Default interface{}
}

// IsDefault returns true if the field's value is the default value.
func (c *Value) IsDefault() bool {
	if v, ok := c.Value.(string); ok {
		return strings.EqualFold(v, c.Default.(string))
	}

	return c.Value == c.Default
}

func init() {
	defaultConfig = &Config{
		LogLevel:           DefaultLogLevel,
		SendUsageData:      TernaryValues.Unknown,
		PreReleaseFeatures: TernaryValues.Unknown,
	}

	cfgDir, err := utils.GetDefaultConfigDirectory()
	if err != nil {
		log.Fatalf("error building default config directory: %s", err)
	}

	DefaultConfigDirectory = cfgDir
	defaultConfig.PluginDir = DefaultConfigDirectory + "/plugins"
}

// LoadConfig loads the configuration from disk, substituting defaults
// if the file does not exist.
func LoadConfig(configDir string) (*Config, error) {
	log.Debugf("loading config file from %s", configDir)

	if configDir == "" {
		configDir = DefaultConfigDirectory
	} else {
		configDir = os.ExpandEnv(configDir)
	}

	cfg, err := load(configDir)
	if err != nil {
		return nil, err
	}

	cfg.setLogger()
	cfg.configDir = configDir

	return cfg, nil
}

func (c *Config) setLogger() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        time.RFC3339,
		DisableLevelTruncation: true,
	})

	switch level := strings.ToUpper(c.LogLevel); level {
	case "TRACE":
		log.SetLevel(log.TraceLevel)
		log.SetReportCaller(true)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

// List outputs a list of all the configuration values
func (c *Config) List() {
	output.Text(c.getAll(""))
}

// Delete deletes a config value.
// This has the effect of reverting the value back to its default.
func (c *Config) Delete(key string) error {
	defaultValue, err := c.getDefaultValue(key)
	if err != nil {
		return err
	}

	err = c.set(key, defaultValue)
	if err != nil {
		return err
	}

	output.Printf("%s %s removed successfully\n", text.FgGreen.Sprint("✔"), text.Bold.Sprint(key))

	return nil
}

// Get retrieves a config value.
func (c *Config) Get(key string) {
	output.Text(c.getAll(key))
}

// Set is used to update a config value.
func (c *Config) Set(key string, value interface{}) error {
	if !stringInStrings(key, validConfigKeys()) {
		return fmt.Errorf("\"%s\" is not a valid key; Please use one of: %s", key, validConfigKeys())
	}

	err := c.set(key, value)
	if err != nil {
		return err
	}

	output.Printf("%s set to %s\n", text.Bold.Sprint(key), text.FgCyan.Sprint(value))

	return nil
}

func load(configDir string) (*Config, error) {
	cfgViper, err := readConfig(configDir)
	if err != nil {
		return nil, err
	}

	allScopes, err := unmarshalAllScopes(cfgViper)

	if err != nil {
		return nil, err
	}

	config, ok := (*allScopes)[globalScopeIdentifier]
	if !ok {
		config = Config{}
	}

	err = config.setDefaults()
	if err != nil {
		return nil, err
	}

	err = config.applyOverrides()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) createFile(path string, cfgViper *viper.Viper) error {
	err := c.visitAllConfigFields(func(v *Value) error {
		cfgViper.Set(globalScopeIdentifier+"."+v.Name, v.Value)
		return nil
	})
	if err != nil {
		return err
	}

	err = os.MkdirAll(c.configDir, os.ModePerm)
	if err != nil {
		return err
	}

	log.Debugf("creating config file at %s: %+v", path, cfgViper.AllSettings())

	err = cfgViper.WriteConfigAs(path)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) getAll(key string) []Value {
	values := []Value{}

	err := c.visitAllConfigFields(func(v *Value) error {
		// Return early if name was supplied and doesn't match
		if key != "" && key != v.Name {
			return nil
		}

		values = append(values, *v)

		return nil
	})
	if err != nil {
		log.Error(err)
	}

	return values
}

func (c *Config) set(key string, value interface{}) error {
	cfgViper, err := readConfig(c.configDir)
	if err != nil {
		return err
	}

	cfgViper.Set(globalScopeIdentifier+"."+key, value)

	allScopes, err := unmarshalAllScopes(cfgViper)
	if err != nil {
		return err
	}

	config, ok := (*allScopes)[globalScopeIdentifier]
	if !ok {
		return fmt.Errorf("failed to locate global scope")
	}

	err = config.setDefaults()
	if err != nil {
		return err
	}

	err = config.validate()
	if err != nil {
		return err
	}

	// Update our instance of the config with what was taken from cfgViper.  This
	// is required for the createFile below to function properly, as it relies on
	// the instance values.
	if err := mergo.Merge(c, config, mergo.WithOverride); err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s.%s", c.configDir, DefaultConfigName, DefaultConfigType)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		createErr := c.createFile(path, cfgViper)
		if createErr != nil {
			return createErr
		}
	} else {
		log.Debugf("writing config file at %s", path)
		err = cfgViper.WriteConfigAs(path)
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

func (c *Config) getDefaultValue(key string) (interface{}, error) {
	var dv interface{}
	var found bool

	err := c.visitAllConfigFields(func(v *Value) error {
		if key == v.Name {
			dv = v.Default
			found = true
		}

		return nil
	})

	if err != nil {
		return dv, err
	}

	if found {
		return dv, nil
	}

	return nil, fmt.Errorf("failed to locate default value for %s", key)
}

func (c *Config) applyOverrides() error {
	log.Debug("setting config overrides")

	if v := os.Getenv("NEW_RELIC_CLI_PRERELEASEFEATURES"); v != "" {
		c.PreReleaseFeatures = Ternary(v)
	}

	return nil
}

func (c *Config) setDefaults() error {
	log.Debug("setting config default")

	if c == nil {
		return nil
	}

	if err := mergo.Merge(c, defaultConfig); err != nil {
		return err
	}

	return nil
}

func (c *Config) validate() error {
	err := c.visitAllConfigFields(func(v *Value) error {
		switch k := strings.ToLower(v.Name); k {
		case "loglevel":
			validValues := []string{"Info", "Debug", "Trace", "Warn", "Error"}
			if !stringInStringsIgnoreCase(v.Value.(string), validValues) {
				return fmt.Errorf("\"%s\" is not a valid %s value; Please use one of: %s", v.Value, v.Name, validValues)
			}
		case "sendusagedata", "prereleasefeatures":
			err := (v.Value.(Ternary)).Valid()
			if err != nil {
				return fmt.Errorf("invalid value for '%s': %s", v.Name, err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *Config) visitAllConfigFields(f func(*Value) error) error {
	cfgType := reflect.TypeOf(*c)
	cfgValue := reflect.ValueOf(*c)
	defaultCfgValue := reflect.ValueOf(*defaultConfig)

	// Iterate through the fields in the struct
	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		name := field.Tag.Get("mapstructure")
		value := cfgValue.Field(i).Interface()
		defaultValue := defaultCfgValue.Field(i).Interface()

		err := f(&Value{
			Name:    name,
			Value:   value,
			Default: defaultValue,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func unmarshalAllScopes(cfgViper *viper.Viper) (*map[string]Config, error) {
	cfgMap := map[string]Config{}
	err := cfgViper.Unmarshal(&cfgMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config with error: %v", err)
	}

	return &cfgMap, nil
}

func readConfig(configDir string) (*viper.Viper, error) {
	cfgViper := viper.New()
	cfgViper.SetEnvPrefix(DefaultEnvPrefix)
	cfgViper.SetConfigName(DefaultConfigName)
	cfgViper.SetConfigType(DefaultConfigType)
	cfgViper.AddConfigPath(configDir) // adding provided directory as search path
	cfgViper.AutomaticEnv()           // read in environment variables that match

	err := cfgViper.ReadInConfig()
	// nolint
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug("no config file found, using defaults")
		} else if e, ok := err.(viper.ConfigParseError); ok {
			return nil, fmt.Errorf("error parsing config file: %v", e)
		}
	}

	return cfgViper, nil
}

func validConfigKeys() []string {
	var keys []string

	cfgType := reflect.TypeOf(Config{})
	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		name := field.Tag.Get("mapstructure")
		keys = append(keys, name)
	}

	return keys
}

func stringInStrings(s string, ss []string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}

	return false
}

// Function ignores the case
func stringInStringsIgnoreCase(s string, ss []string) bool {
	for _, v := range ss {
		if strings.EqualFold(v, s) {
			return true
		}
	}

	return false
}
