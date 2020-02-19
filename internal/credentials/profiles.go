package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/newrelic/newrelic-cli/internal/config"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// DefaultProfileFile is the configuration file containing the default profile name
const DefaultProfileFile = "default-profile"

// Profile contains data of a single profile
type Profile struct {
	AdminAPIKey    string `mapstructure:"adminAPIKey"` // For accessing New Relic (Rest v2)
	PersonalAPIKey string `mapstructure:"apiKey"`      // For accessing New Relic GraphQL resources
	Region         string `mapstructure:"region"`      // Region to use for New Relic resources
}

// LoadProfiles reads the credential profiles from the default path.
func LoadProfiles(configDir string) (*map[string]Profile, error) {
	err := config.InitializeConfigDirectory(configDir)
	if err != nil {
		return nil, fmt.Errorf("error initializing config directory %s: %s", configDir, err)
	}

	cfgViper, err := readCredentials(configDir)
	if err != nil {
		return nil, fmt.Errorf("error while reading credentials: %s", err)
	}

	creds, err := unmarshalProfiles(cfgViper)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling profiles: %s", err)
	}

	return creds, nil
}

// LoadDefaultProfile reads the profile name from the default profile file.
func LoadDefaultProfile(configDir string) (string, error) {
	err := config.InitializeConfigDirectory(configDir)
	if err != nil {
		return "", fmt.Errorf("error initializing config directory %s: %s", configDir, err)
	}

	defProfile, err := readDefaultProfile(configDir)
	if err != nil {
		return "", err
	}

	return defProfile, nil
}

func readDefaultProfile(configDir string) (string, error) {
	var defaultProfile string

	cfgViper := viper.New()
	cfgViper.SetConfigName(DefaultProfileFile)
	cfgViper.SetConfigType(defaultConfigType)
	cfgViper.AddConfigPath(configDir)

	// ReadInConfig must be called here, even though we receive an error back,
	// ConfigFileUsed() does not return the value without this call here.
	cfgViper.ReadInConfig()
	// if err != nil {
	// 	return nil, err
	// }

	// Since Viper requires key:value, we manually read it again and unmarshal the JSON...
	byteValue, err := ioutil.ReadFile(cfgViper.ConfigFileUsed())
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(byteValue, &defaultProfile)
	if err != nil {
		return "", err
	}

	return defaultProfile, nil
}

func readCredentials() (*viper.Viper, error) {
	credViper := viper.New()
	credViper.SetConfigName(DefaultCredentialsFile)
	credViper.SetConfigType(defaultConfigType)
	credViper.SetEnvPrefix(config.DefaultEnvPrefix)
	credViper.AddConfigPath(config.DefaultConfigDirectory) // adding home directory as first search path
	credViper.AddConfigPath(".")                           // current directory to search path
	credViper.AutomaticEnv()                               // read in environment variables that match

	// Read in config
	err := credViper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("no profile configuration found")
		} else if e, ok := err.(viper.ConfigParseError); ok {
			return nil, fmt.Errorf("error parsing profile config file: %v", e)
		}
	}

	return credViper, nil
}

func unmarshalProfiles(cfgViper *viper.Viper) (*map[string]Profile, error) {
	cfgMap := map[string]Profile{}
	err := cfgViper.Unmarshal(&cfgMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials with error: %v", err)
	}

	log.Debugf("loaded credentials from: %v", cfgViper.ConfigFileUsed())

	return &cfgMap, nil
}

func (p *Profile) validate() error {
	if p.Region == "" {
		return fmt.Errorf("Profile.Region is required")
	}

	if p.AdminAPIKey == "" || p.PersonalAPIKey == "" {
		return fmt.Errorf("Profile.AdminAPIKey or Profile.PersonalAPIKey is required")
	}

	return nil
}
