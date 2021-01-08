package configuration

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	configType            = "json"
	configEnvPrefix       = "NEW_RELIC_CLI"
	globalScopeIdentifier = "*"
)

// Config keys
const (
	logLevelKey           = "loglevel"
	pluginDirectoryKey    = "plugindir"
	prereleaseFeaturesKey = "prereleasefeatures"
	sendDataUsageKey      = "sendusagedata"
)

// Credential keys
const (
	apiKeyKey     = "apiKey"
	regionKey     = "region"
	accountIDKey  = "accountID"
	licenseKeyKey = "licenseKey"
)

var (
	configDir              string
	configFileName         = "config.json"
	credsFileName          = "credentials.json"
	defaultProfileFileName = "default-profile.json"
	defaultProfileValue    string
	viperConfig            *viper.Viper
	viperCreds             *viper.Viper
	viperDefaultProfile    *viper.Viper
)

// TODO: SetDefaultProfile(profileName string) {}

func GetConfigValue(key string) interface{} {
	return viperConfig.Get(keyGlobalScope(key))
}

func SetLogLevel(logLevel string) error {
	return setConfigValue(logLevelKey, logLevel)
}

func SetPluginDirectory(directory string) error {
	return setConfigValue(pluginDirectoryKey, directory)
}

func SetPreleaseFeatures(prereleaseFeatures string) error {
	return setConfigValue(prereleaseFeaturesKey, prereleaseFeatures)
}

func SetSendUsageData(sendUsageData string) error {
	return setConfigValue(sendDataUsageKey, sendUsageData)
}

func setConfigValue(key string, value string) error {
	viperConfig.Set(keyGlobalScope(key), value)

	err := viperConfig.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func GetCredentialValue(key string) interface{} {
	return viperCreds.Get(keyDefaultProfile(key))
}

func SetAPIKey(profileName string, apiKey string) error {
	return setCredentialValue(profileName, apiKeyKey, apiKey)
}

func SetRegion(profileName string, region string) error {
	return setCredentialValue(profileName, regionKey, region)
}

func SetAccountID(profileName string, accountID string) error {
	return setCredentialValue(profileName, accountIDKey, accountID)
}

func SetLicenseKey(profileName string, licenseKey string) error {
	return setCredentialValue(profileName, licenseKeyKey, licenseKey)
}

func setCredentialValue(profileName string, key string, value string) error {
	keyPath := fmt.Sprintf("%s.%s", profileName, key)
	viperCreds.Set(keyPath, value)

	err := viperCreds.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func load() error {
	// if configDirectory == "" {
	// 	configDir, err := getDefaultConfigDirectory()
	// 	if err != nil {
	// 		return err
	// 	}
	// } else {
	// 	configDir = configDirectory
	// }

	if err := loadConfigFile(); err != nil {
		return err
	}

	if err := loadCredsFile(); err != nil {
		return err
	}

	if err := loadDefaultProfileFile(); err != nil {
		return err
	}

	return nil
}

func loadConfigFile() error {
	viperConfig = viper.New()
	viperConfig.SetEnvPrefix(configEnvPrefix)
	viperConfig.SetConfigName(configFileName)
	viperConfig.SetConfigType(configType)
	viperConfig.AddConfigPath(configDir)
	viperConfig.AutomaticEnv()

	if err := loadFile(viperConfig); err != nil {
		return fmt.Errorf("error loading config file: %s", err)
	}

	return nil
}

func loadCredsFile() error {
	viperCreds = viper.New()
	viperCreds.SetEnvPrefix(configEnvPrefix)
	viperCreds.SetConfigName(credsFileName)
	viperCreds.SetConfigType(configType)
	viperCreds.AddConfigPath(configDir)
	viperCreds.AutomaticEnv()

	if err := loadFile(viperCreds); err != nil {
		return fmt.Errorf("error loading credentials file: %s", err)
	}

	return nil
}

func loadDefaultProfileFile() error {
	defaultProfileFilePath := filepath.Join(configDir, defaultProfileFileName)
	defaultProfileBytes, err := ioutil.ReadFile(defaultProfileFilePath)
	if err != nil {
		return err
	}

	defaultProfileValue = strings.Trim(string(defaultProfileBytes), "\"")

	return nil
}

func loadFile(v *viper.Viper) error {
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug("file not found, using defaults")
		} else if e, ok := err.(viper.ConfigParseError); ok {
			return e
		}
	}

	return nil
}

func keyGlobalScope(key string) string {
	return fmt.Sprintf("%s.%s", globalScopeIdentifier, key)
}

func keyDefaultProfile(key string) string {
	return fmt.Sprintf("%s.%s", defaultProfile(), key)
}

func defaultProfile() string {
	return "default"
}

func getDefaultConfigDirectory() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/.newrelic", home), nil
}
