package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// DefaultCredentialsFile is the default place to load profiles from
const DefaultCredentialsFile = "credentials"

// DefaultProfileFile is the configuration file containing the default profile name
const DefaultProfileFile = "default-profile"

const defaultConfigType = "json"

// Load reads all profile information from configuration files
func Load(configDir, configEnvPrefix string) (*Credentials, error) {
	credViper := viper.New()
	credViper.SetConfigName(DefaultCredentialsFile)
	credViper.SetConfigType(defaultConfigType)
	credViper.SetEnvPrefix(configEnvPrefix)
	credViper.AddConfigPath(configDir) // adding home directory as first search path
	credViper.AddConfigPath(".")       // current directory to search path
	credViper.AutomaticEnv()           // read in environment variables that match

	// Read in config
	err := credViper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("no profile configuration found")
		} else if e, ok := err.(viper.ConfigParseError); ok {
			return nil, fmt.Errorf("error parsing profile config file: %v", e)
		}
	}

	log.Debugf("loaded profiles from: %v", credViper.ConfigFileUsed())

	creds := Credentials{
		ConfigDirectory: configDir,
	}

	// Read the profiles
	err = credViper.Unmarshal(&creds.Profiles)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal profile config with error: %v", err)
	}

	// Read the default profile configuration (just a string, no key/value)
	defViper := viper.New()
	defViper.SetConfigName(DefaultProfileFile)
	defViper.AddConfigPath(configDir) // adding home directory as first search path
	defViper.AddConfigPath(".")       // current directory to search path

	err = defViper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("no profile configuration found")
		} else if _, ok := err.(viper.ConfigParseError); ok {
			// Since Viper requires key:value, we manually read it again and unmarshal the JSON...
			byteValue, err := ioutil.ReadFile(defViper.ConfigFileUsed())
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(byteValue, &creds.DefaultProfile)
			if err != nil {
				return nil, err
			}
		}
	}

	log.Debugf("default profile: '%s'", creds.DefaultProfile)

	return &creds, nil
}
