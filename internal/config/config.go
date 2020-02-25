package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/imdario/mergo"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// DefaultConfigName is the default name of the global configuration file
	DefaultConfigName = "config"

	// DefaultConfigType to read, though any file type supported by viper is allowed
	DefaultConfigType = "json"

	// DefaultEnvPrefix is used when reading environment variables
	DefaultEnvPrefix = "newrelic"

	// DefaultLogLevel is the default log level
	DefaultLogLevel = "INFO"

	// DefaultSendUsageData is the default value for sendUsageData
	DefaultSendUsageData = "NOT_ASKED"

	globalScopeIdentifier = "*"
)

var (
	// DefaultConfigDirectory is the default location for the CLI config files
	DefaultConfigDirectory string

	// DefaultPluginDirectory is the default sub-directory containing the plugins
	DefaultPluginDirectory string

	renderer      = TableRenderer{}
	defaultConfig *Config
)

// Config contains the main CLI configuration
type Config struct {
	LogLevel      string `mapstructure:"logLevel"`      // LogLevel for verbose output
	PluginDir     string `mapstructure:"pluginDir"`     // PluginDir is the directory where plugins will be installed
	SendUsageData string `mapstructure:"sendUsageData"` // SendUsageData enables sending usage statistics to New Relic
}

// Value represents an instance of a configuration field.
type Value struct {
	Name    string
	Value   interface{}
	Default interface{}
}

// IsDefault returns tru if the field's value is the default value.
func (c *Value) IsDefault() bool {
	return c.Value == c.Default
}

func init() {
	defaultConfig = &Config{
		LogLevel:      DefaultLogLevel,
		PluginDir:     DefaultPluginDirectory,
		SendUsageData: DefaultSendUsageData,
	}

	cfgDir, err := getDefaultConfigDirectory()
	if err != nil {
		log.Fatalf("error building default config directory")
	}

	DefaultConfigDirectory = cfgDir
	DefaultPluginDirectory = DefaultConfigDirectory + "/plugins"
}

// LoadConfig loads the configuration from disk, or initializes a new file
// if one doesn't currently exist.
func LoadConfig() (*Config, error) {
	cfg, err := load()
	if err != nil {
		return nil, err
	}

	cfg.setLogger()

	return cfg, nil
}

func getDefaultConfigDirectory() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/.newrelic", home), nil
}

func (c *Config) setLogger() {
	switch level := strings.ToUpper(c.LogLevel); level {
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
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
	renderer.List(c)
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

	renderer.Delete(key)
	return nil
}

// Get retrieves a config value.
func (c *Config) Get(key string) {
	renderer.Get(c, key)
}

// Set sets a config value.
func (c *Config) Set(key string, value string) error {
	if !stringInStrings(key, validConfigKeys()) {
		return fmt.Errorf("\"%s\" is not a valid key; Please use one of: %s", key, validConfigKeys())
	}

	switch k := strings.ToLower(key); k {
	case "loglevel":
		validValues := []string{"Info", "Debug", "Trace", "Warn", "Error"}
		if !stringInStrings(value, validValues) {
			return fmt.Errorf("\"%s\" is not a valid %s value; Please use one of: %s", value, key, validValues)
		}
	case "sendusagedata":
		validValues := []string{"NOT_ASKED", "DISALLOW", "ALLOW"}
		if !stringInStrings(value, validValues) {
			return fmt.Errorf("\"%s\" is not a valid %s value; Please use one of: %s", value, key, validValues)
		}
	}

	k := strings.ToLower(key)
	v := strings.ToUpper(value)

	err := c.set(k, v)
	if err != nil {
		return err
	}

	renderer.Set(k, v)
	return nil
}

func load() (*Config, error) {
	log.Debug("loading config file")

	cfgViper, err := readConfig()
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

	return &config, nil
}

func (c *Config) createFile(cfgViper *viper.Viper) error {
	c.visitAllConfigFields(func(v *Value) {
		cfgViper.Set(globalScopeIdentifier+"."+v.Name, v.Value.(string))
	})

	err := os.MkdirAll(DefaultConfigDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s.%s", DefaultConfigDirectory, DefaultConfigName, DefaultConfigType)
	err = cfgViper.WriteConfigAs(path)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) get(key string) []Value {
	return c.getAll(key)
}

func (c *Config) getAll(key string) []Value {
	values := []Value{}

	c.visitAllConfigFields(func(v *Value) {
		// Return early if name was supplied and doesn't match
		if key != "" && key != v.Name {
			return
		}

		values = append(values, *v)
	})

	return values
}

func (c *Config) set(key string, value interface{}) error {
	cfgViper, err := readConfig()

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

	err = config.validate()

	if err != nil {
		return err
	}

	err = c.createFile(cfgViper)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) getDefaultValue(key string) (interface{}, error) {
	var dv interface{}
	var found bool
	c.visitAllConfigFields(func(v *Value) {
		if key == v.Name {
			dv = v.Default
			found = true
			return
		}
	})

	if found {
		return dv, nil
	}

	return nil, fmt.Errorf("failed to locate default value for %s", key)
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
	// TODO: implement this
	return nil
}

func (c *Config) visitAllConfigFields(f func(*Value)) {
	cfgType := reflect.TypeOf(*c)
	cfgValue := reflect.ValueOf(*c)
	defaultCfgValue := reflect.ValueOf(*defaultConfig)

	// Iterate through the fields in the struct
	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		name := field.Tag.Get("mapstructure")
		value := cfgValue.Field(i).Interface()
		defaultValue := defaultCfgValue.Field(i).Interface()

		f(&Value{
			Name:    name,
			Value:   value,
			Default: defaultValue,
		})
	}
}

func unmarshalAllScopes(cfgViper *viper.Viper) (*map[string]Config, error) {
	cfgMap := map[string]Config{}
	err := cfgViper.Unmarshal(&cfgMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config with error: %v", err)
	}

	log.Debugf("loaded config from: %v", cfgViper.ConfigFileUsed())

	return &cfgMap, nil
}

func readConfig() (*viper.Viper, error) {
	cfgViper := viper.New()
	cfgViper.SetEnvPrefix(DefaultEnvPrefix)
	cfgViper.SetConfigName(DefaultConfigName)
	cfgViper.SetConfigType(DefaultConfigType)
	cfgViper.AddConfigPath(DefaultConfigDirectory) // adding home directory as first search path
	cfgViper.AddConfigPath(".")                    // current directory to search path
	cfgViper.AutomaticEnv()                        // read in environment variables that match

	err := cfgViper.ReadInConfig()
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
		if strings.EqualFold(v, s) {
			return true
		}
	}

	return false
}
