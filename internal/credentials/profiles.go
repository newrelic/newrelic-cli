package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/newrelic/newrelic-client-go/pkg/region"

	"github.com/newrelic/newrelic-cli/internal/config"
)

// DefaultProfileFile is the configuration file containing the default profile name
const DefaultProfileFile = "default-profile"

// Profile contains data of a single profile
type Profile struct {
	APIKey            string `mapstructure:"apiKey" json:"apiKey,omitempty"`                       // For accessing New Relic GraphQL resources
	InsightsInsertKey string `mapstructure:"insightsInsertKey" json:"insightsInsertKey,omitempty"` // For posting custom events
	Region            string `mapstructure:"region" json:"region,omitempty"`                       // Region to use for New Relic resources
	AccountID         int    `mapstructure:"accountID" json:"accountID,omitempty"`                 // AccountID to use for New Relic resources
	LicenseKey        string `mapstructure:"licenseKey" json:"licenseKey,omitempty"`               // License key to use for agent config and ingest
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
	var p *Profile
	if c.DefaultProfile != "" {
		if val, ok := c.Profiles[c.DefaultProfile]; ok {
			p = &val
		}
	}

	p = applyOverrides(p)
	return p
}

// applyOverrides reads Profile info out of the Environment to override config
func applyOverrides(p *Profile) *Profile {
	envAPIKey := os.Getenv("NEW_RELIC_API_KEY")
	envInsightsInsertKey := os.Getenv("NEW_RELIC_INSIGHTS_INSERT_KEY")
	envRegion := os.Getenv("NEW_RELIC_REGION")
	envAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")

	if envAPIKey == "" && envRegion == "" && envInsightsInsertKey == "" && envAccountID == "" {
		return p
	}

	out := Profile{}
	if p != nil {
		out = *p
	}

	if envAPIKey != "" {
		out.APIKey = envAPIKey
	}

	if envInsightsInsertKey != "" {
		out.InsightsInsertKey = envInsightsInsertKey
	}

	if envRegion != "" {
		out.Region = strings.ToUpper(envRegion)
	}

	if envAccountID != "" {
		accountID, err := strconv.Atoi(envAccountID)
		if err != nil {
			log.Warnf("Invalid account ID: %s", envAccountID)
			return &out
		}

		out.AccountID = accountID
	}

	return &out
}

func readDefaultProfile(configDir string) (string, error) {
	var defaultProfile string

	_, err := os.Stat(configDir)
	if err != nil {
		return "", fmt.Errorf("unable to read default-profile from %s: %s", configDir, err)
	}

	configPath := os.ExpandEnv(fmt.Sprintf("%s/%s.%s", configDir, DefaultProfileFile, defaultConfigType))

	// The default-profile.json is a quoted string of the name for the default profile.
	byteValue, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("error while reading default profile file %s: %s", configPath, err)
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

	// Have to pass in the default hooks to add one...
	err := cfgViper.Unmarshal(&cfgMap,
		viper.DecodeHook(
			mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToSliceHookFunc(","),
				StringToRegionHookFunc(), // Custom parsing of Region on unmarshal
			),
		))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials with error: %v", err)
	}

	log.Debugf("loaded credentials from: %v", cfgViper.ConfigFileUsed())

	return &cfgMap, nil
}

// MarshalJSON allows us to override the default behavior on marshal
// and lowercase the region string for backwards compatibility
func (p Profile) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		APIKey            string `json:"apiKey,omitempty"`
		InsightsInsertKey string `json:"insightsInsertKey,omitempty"`
		Region            string `json:"region,omitempty"`
		AccountID         int    `json:"accountID,omitempty"`
		LicenseKey        string `json:"licenseKey,omitempty"`
	}{
		APIKey:            p.APIKey,
		InsightsInsertKey: p.InsightsInsertKey,
		AccountID:         p.AccountID,
		LicenseKey:        p.LicenseKey,
		Region:            strings.ToLower(p.Region),
	})
}

// StringToRegionHookFunc takes a string and runs it through the region
// parser to create a valid region (or error)
func StringToRegionHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		var n region.Name

		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(n) {
			return data, nil
		}

		// Convert it by parsing
		reg, err := region.Parse(data.(string))
		return reg, err
	}
}
