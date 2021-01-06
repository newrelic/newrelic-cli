package configuration

import (
	"fmt"

	"github.com/apex/log"
	"github.com/spf13/viper"
)

const (
	configType            = "json"
	configEnvPrefix       = "NEW_RELIC_CLI"
	globalScopeIdentifier = "*"
)

var (
	configDir              string
	configFileName         = "config.json"
	credsFileName          = "credentials.json"
	defaultProfileFileName = "default-profile.json"
	viperConfig            *viper.Viper
	viperCreds             *viper.Viper
	viperDefaultProfile    *viper.Viper
)

func GetConfigValue(key string) interface{} {
	return viperConfig.Get(keyGlobalScope(key))
}

func GetProfileValue(key string) interface{} {
	return viperCreds.Get(keyDefaultProfile(key))
}

func load() error {
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
	viperDefaultProfile = viper.New()
	viperDefaultProfile.SetEnvPrefix(configEnvPrefix)
	viperDefaultProfile.SetConfigName(defaultProfileFileName)
	viperDefaultProfile.SetConfigType(configType)
	viperDefaultProfile.AddConfigPath(configDir)
	viperDefaultProfile.AutomaticEnv()

	if err := loadFile(viperDefaultProfile); err != nil {
		return fmt.Errorf("error loading credentials file: %s", err)
	}

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
