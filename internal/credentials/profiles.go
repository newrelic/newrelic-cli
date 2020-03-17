package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/newrelic/newrelic-cli/internal/config"
)

// DefaultProfileFile is the configuration file containing the default profile name
const DefaultProfileFile = "default-profile"

// Profile contains data of a single profile
type Profile struct {
	APIKey string `mapstructure:"apiKey" json:"apiKey,omitempty"` // For accessing New Relic GraphQL resources
	Region string `mapstructure:"region" json:"region,omitempty"` // Region to use for New Relic resources
}

// LoadProfiles reads the credential profiles from the default path.
func LoadProfiles(configDir string) (*map[string]Profile, error) {
	cfgViper, err := readCredentials(configDir)

	if err != nil {
		return &map[string]Profile{}, fmt.Errorf("error while reading credentials: %s", err)
	}

	creds, err := unmarshalProfiles(cfgViper)
	if err != nil {
		return &map[string]Profile{}, fmt.Errorf("error unmarshaling profiles: %s", err)
	}

	return creds, nil
}

// LoadDefaultProfile reads the profile name from the default profile file.
func LoadDefaultProfile(configDir string) (string, error) {
	defProfile, err := readDefaultProfile(configDir)
	if err != nil {
		return "", err
	}

	return defProfile, nil
}

// Default returns the default profile
func (c *Credentials) Default() *Profile {
	if c.DefaultProfile != "" {
		if val, ok := c.Profiles[c.DefaultProfile]; ok {
			return &val
		}
	}

	return nil
}

func readDefaultProfile(configDir string) (string, error) {
	var defaultProfile string

	cfgViper := viper.New()
	cfgViper.SetConfigName(DefaultProfileFile)
	cfgViper.SetConfigType(defaultConfigType)
	cfgViper.AddConfigPath(configDir)

	// ReadInConfig must be called here, even though we receive an error back,
	// ConfigFileUsed() does not return the value without this call here.
	err := cfgViper.ReadInConfig()
	if err != nil {
		log.Debug(err)
	}

	// Since Viper requires key:value, we manually read it again and unmarshal the JSON...
	byteValue, err := ioutil.ReadFile(cfgViper.ConfigFileUsed())
	if err != nil {
		return "", fmt.Errorf("error while reading default profile file %s: %s", cfgViper.ConfigFileUsed(), err)
	}
	err = json.Unmarshal(byteValue, &defaultProfile)
	if err != nil {
		return "", fmt.Errorf("error while unmarshaling default profile: %s", err)
	}

	return defaultProfile, nil
}

func readCredentials(configDir string) (*viper.Viper, error) {
	credViper := viper.New()
	credViper.SetConfigName(DefaultCredentialsFile)
	credViper.SetConfigType(defaultConfigType)
	credViper.SetEnvPrefix(config.DefaultEnvPrefix)
	credViper.AddConfigPath(configDir) // adding home directory as first search path
	credViper.AutomaticEnv()           // read in environment variables that match

	// Read in config
	err := credViper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

			filePath := os.ExpandEnv(fmt.Sprintf("%s/%s.json", configDir, DefaultCredentialsFile))

			err = credViper.WriteConfigAs(filePath)
			if err != nil {
				return nil, fmt.Errorf("error initializing new configuration directory %s: %s", filePath, err)
			}
		}

		if e, ok := err.(viper.ConfigParseError); ok {
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
