package configuration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	configType            = "json"
	globalScopeIdentifier = "*"
)

type ConfigFieldKey string
type ProfileFieldKey string

const (
	LogLevel           ConfigFieldKey = "loglevel"
	PluginDir          ConfigFieldKey = "plugindir"
	PrereleaseFeatures ConfigFieldKey = "prereleasefeatures"
	SendUsageData      ConfigFieldKey = "sendusagedata"

	APIKey            ProfileFieldKey = "apiKey"
	Region            ProfileFieldKey = "region"
	AccountID         ProfileFieldKey = "accountID"
	InsightsInsertKey ProfileFieldKey = "insightsInsertKey"
	LicenseKey        ProfileFieldKey = "licenseKey"
)

var (
	ConfigFields = []ConfigField{
		{
			Name:    "LogLevel",
			Key:     LogLevel,
			Default: "info",
		},
		{
			Name:    "SendUsageData",
			Key:     SendUsageData,
			Default: string(TernaryValues.Unknown),
		},
		{
			Name:    "PluginDir",
			Key:     PluginDir,
			Default: "",
		},
		{
			Name:    "PrereleaseFeatures",
			Key:     PrereleaseFeatures,
			Default: string(TernaryValues.Unknown),
		},
	}
	ProfileFields = []ProfileField{
		{
			Name:        "APIKey",
			Key:         APIKey,
			EnvOverride: "NEW_RELIC_API_KEY",
		},
		{
			Name:        "Region",
			Key:         Region,
			EnvOverride: "NEW_RELIC_REGION",
		},
		{
			Name:        "AccountID",
			Key:         AccountID,
			EnvOverride: "NEW_RELIC_ACCOUNT_ID",
		},
		{
			Name:        "LicenseKey",
			Key:         LicenseKey,
			EnvOverride: "NEW_RELIC_LICENSE_KEY",
		},
		{
			Name:        "InsightsInsertKey",
			Key:         InsightsInsertKey,
			EnvOverride: "NEW_RELIC_INSIGHTS_INSERT_KEY",
		},
	}
	ConfigDir                 string
	EnvVarResolver            envResolver = &osEnvResolver{}
	configFilename                        = "config.json"
	credsFilename                         = "credentials.json"
	defaultProfileFilename                = "default-profile.json"
	defaultDefaultProfileName             = "default"
)

type ConfigField struct {
	Name    string
	Key     ConfigFieldKey
	Default string
}

type ProfileField struct {
	Name        string
	Key         ProfileFieldKey
	EnvOverride string
}

type ConfigValue struct {
	Name    string
	Value   interface{}
	Default interface{}
}

// IsDefault returns true if the field's value is the default value.
func (c *ConfigValue) IsDefault() bool {
	if v, ok := c.Value.(string); ok {
		return strings.EqualFold(v, c.Default.(string))
	}

	return c.Value == c.Default
}

func init() {
	var err error
	ConfigDir, err = getDefaultConfigDirectory()
	if err != nil {
		log.Debug(err)
	}
}

func GetConfigValue(key ConfigFieldKey) (interface{}, error) {
	if ok := isValidConfigKey(key); !ok {
		return nil, fmt.Errorf("config key %s is not valid.  valid keys are %s", key, validConfigFieldKeys())
	}

	return config().Get(keyGlobalScope(string(key))), nil
}

func GetProfileValue(profileName string, key ProfileFieldKey) (interface{}, error) {
	if ok := isValidCredentialKey(key); !ok {
		return nil, fmt.Errorf("credential key %s is not valid.  valid keys are %s", key, validProfileFieldKeys())
	}

	if o := getProfileValueEnvOverride(key); o != "" {
		return o, nil
	}

	return profiles().Get(keyDefaultProfile(string(key))), nil
}

func GetActiveProfileValue(key ProfileFieldKey) (interface{}, error) {
	return GetProfileValue(GetActiveProfileName(), key)
}

func GetActiveProfileName() string {
	return defaultProfileName()
}

func GetDefaultProfileName() string {
	return defaultProfileName()
}

func SetDefaultProfileName(profileName string) error {
	return saveDefaultProfileName(profileName)
}

func SetConfigValue(key ConfigFieldKey, value string) error {
	c := config()
	c.Set(keyGlobalScope(string(key)), value)

	cfgFilePath := path.Join(ConfigDir, configFilename)
	if err := c.WriteConfigAs(cfgFilePath); err != nil {
		return err
	}

	return nil
}

func SetActiveProfileValue(key ProfileFieldKey, value interface{}) error {
	return SetProfileValue(GetActiveProfileName(), key, value)
}

func SetProfileValue(profileName string, key ProfileFieldKey, value interface{}) error {
	p := profiles()
	keyPath := fmt.Sprintf("%s.%s", profileName, key)
	p.Set(keyPath, value)

	credsFilePath := path.Join(ConfigDir, credsFilename)
	if err := p.WriteConfigAs(credsFilePath); err != nil {
		return err
	}

	if defaultProfileName() == "" {
		log.Debugf("setting %s as default profile", defaultDefaultProfileName)
		if err := SetDefaultProfileName(profileName); err != nil {
			return err
		}
	}

	return nil
}

func RemoveProfile(profileName string) error {
	p := profiles()
	configMap := p.AllSettings()
	delete(configMap, profileName)

	encodedConfig, _ := json.MarshalIndent(configMap, "", " ")
	err := viper.ReadConfig(bytes.NewReader(encodedConfig))
	if err != nil {
		return err
	}

	credsFilePath := path.Join(ConfigDir, credsFilename)
	if err := p.WriteConfigAs(credsFilePath); err != nil {
		return err
	}

	return nil
}

func GetProfileNames() []string {
	profileMap := map[string]interface{}{}
	if err := profiles().Unmarshal(&profileMap); err != nil {
		log.Debug(err)
		return []string{}
	}

	n := []string{}
	for k := range profileMap {
		n = append(n, k)
	}

	return n
}

func getProfileValueEnvOverride(key ProfileFieldKey) string {
	for _, p := range ProfileFields {
		if !strings.EqualFold(string(p.Key), string(key)) {
			continue
		}

		e := EnvVarResolver.Getenv(p.EnvOverride)

		if e != "" {
			return e
		}
	}

	return ""
}

func config() *viper.Viper {
	v, err := loadConfigFile()
	if err != nil {
		if err == os.ErrNotExist {
			log.Debug("config file not found, writing defaults")
			err = writeConfigDefaults(v)
			if err != nil {
				log.Fatal("could not write config defaults")
			}
		}

		log.Debug(err)
	}

	return v
}

func profiles() *viper.Viper {
	v, err := loadCredsFile()
	if err != nil {
		log.Debug(err)
	}

	return v
}

func writeConfigDefaults(v *viper.Viper) error {
	for _, c := range ConfigFields {
		v.Set(keyGlobalScope(string(c.Key)), c.Default)
	}

	if err := v.WriteConfigAs(path.Join(ConfigDir, configFilename)); err != nil {
		return err
	}

	return nil
}

func loadConfigFile() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(configFilename)
	v.SetConfigType(configType)
	v.AddConfigPath(ConfigDir)

	for _, c := range ConfigFields {
		v.SetDefault(fmt.Sprintf("*.%s", c.Key), c.Default)
	}

	if err := loadFile(v); err != nil {
		return nil, err
	}

	return v, nil
}

func loadCredsFile() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(credsFilename)
	v.SetConfigType(configType)
	v.AddConfigPath(ConfigDir)

	if err := loadFile(v); err != nil {
		return nil, fmt.Errorf("credentials file not found: %s", path.Join(ConfigDir, credsFilename))
	}

	return v, nil
}

func defaultProfileName() string {
	p, err := loadDefaultProfileName()
	if err != nil {
		log.Debug("default profile not found")
	}

	return p
}

func loadDefaultProfileName() (string, error) {
	defaultProfileFilePath := filepath.Join(ConfigDir, defaultProfileFilename)
	defaultProfileBytes, err := ioutil.ReadFile(defaultProfileFilePath)
	if err != nil {
		return "", err
	}

	v := strings.Trim(string(defaultProfileBytes), "\"")

	return v, nil
}

func loadFile(v *viper.Viper) error {
	err := v.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Debug("file not found, using defaults")
	} else if e, ok := err.(viper.ConfigParseError); ok {
		return e
	}

	return nil
}

func saveDefaultProfileName(profileName string) error {
	defaultProfileFilePath := filepath.Join(ConfigDir, defaultProfileFilename)

	if err := ioutil.WriteFile(defaultProfileFilePath, []byte("\""+profileName+"\""), 0644); err != nil {
		return err
	}

	return nil
}

func keyGlobalScope(key string) string {
	return fmt.Sprintf("%s.%s", globalScopeIdentifier, key)
}

func keyDefaultProfile(key string) string {
	return fmt.Sprintf("%s.%s", defaultProfileName(), key)
}

func getDefaultConfigDirectory() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/.newrelic", home), nil
}

func isValidConfigKey(key ConfigFieldKey) bool {
	for _, v := range ConfigFields {
		if strings.EqualFold(string(v.Key), string(key)) {
			return true
		}
	}

	return false
}

func isValidCredentialKey(key ProfileFieldKey) bool {
	for _, v := range ProfileFields {
		if strings.EqualFold(string(v.Key), string(key)) {
			return true
		}
	}

	return false
}

func validConfigFieldKeys() []string {
	valid := make([]string, len(ConfigFields))

	for _, v := range ConfigFields {
		valid = append(valid, string(v.Key))
	}

	return valid
}

func validProfileFieldKeys() []string {
	valid := make([]string, len(ProfileFields))

	for _, v := range ProfileFields {
		valid = append(valid, string(v.Key))
	}

	return valid
}
