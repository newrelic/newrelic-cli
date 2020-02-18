package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// DefaultCredentialsFile is the default place to load profiles from
const DefaultCredentialsFile = "credentials"

// DefaultProfileFile is the configuration file containing the default profile name
const DefaultProfileFile = "default-profile"

const defaultConfigType = "json"

const (
	defaultProfileString = " (default)"
	hiddenKeyString      = "<hidden>"
)

// Profile contains data of a single profile
type Profile struct {
	AdminAPIKey    string `mapstructure:"adminAPIKey" json:"adminAPIKey"` // For accessing New Relic (Rest v2)
	PersonalAPIKey string `mapstructure:"apiKey"      json:"apiKey"`      // For accessing New Relic GraphQL resources
	Region         string `mapstructure:"region"      json:"region"`      // Region to use for New Relic resources
}

// Credentials is the metadata around all configured profiles
type Credentials struct {
	DefaultProfile  string
	Profiles        map[string]Profile
	ConfigDirectory string
}

// LoadCredentials loads the list of profiles
func LoadCredentials() (*Credentials, error) {
	var err error

	// Load profiles
	creds, err := Load()
	if err != nil {
		return &Credentials{}, err
	}

	return creds, nil
}

func Load() (*map[string]Credentials, error) {
	log.Debug("loading credentials file")

	cfgViper, err := readCredentails()
	if err != nil {
		return nil, err
	}

	creds, err := unmarshalCredentials(cfgViper)
	if err != nil {
		return nil, err
	}

	return &creds
}

func readCredentials() (*viper.Viper, error) {
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
}

func unmarshalCredentials(cfgViper *viper.Viper) (*map[string]Credentials, error) {
	cfgMap := map[string]Credentials{}
	err := cfgViper.Unmarshal(&cfgMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials with error: %v", err)
	}

	log.Debugf("loaded credentials from: %v", cfgViper.ConfigFileUsed())

	return &cfgMap, nil
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

func (c *Credentials) profileExists(profileName string) bool {
	for k := range c.Profiles {
		if k == profileName {
			return true
		}
	}

	return false
}

func (c *Credentials) Set(profileName string, key string, value string) error {
	err := c.set(profileName, key, value)
	if err != nil {
		return err
	}

	renderer.Set(key, value)
	return nil
}

func (c *Credentials) set(profileName string, key string, value interface{}) error {
	cfgViper, err := readCredentials()
	if err != nil {
		return err
	}

	setPath := fmt.Sprintf("%s.%s", profileName, key)
	cfgViper.Set(setPath, value)

	allCreds, err := unmarshalCredentials(cfgViper)
	if err != nil {
		return err
	}

	creds, ok := (*allCreds)[profileName]
	if !ok {
		return fmt.Errorf("failed to locate global credentials")
	}

	err = creds.validate()
	if err != nil {
		return err
	}

	cfgViper.WriteConfig()

	return nil
}

func (c *Credentials) validate() error {
	// TODO: implement this
	return nil
}

// AddProfile adds a new profile to the credentials file.
func (c *Credentials) AddProfile(profileName, region, apiKey, adminAPIKey string) error {

	if c.profileExists(profileName) {
		return fmt.Errorf("Profile with name %s already exists", profileName)
	}

	p := Profile{
		Region:         region,
		PersonalAPIKey: apiKey,
		AdminAPIKey:    adminAPIKey,
	}

	c.Profiles[profileName] = p

	file, _ := json.MarshalIndent(c.Profiles, "", "  ")

	defaultCredentialsFile := os.ExpandEnv(fmt.Sprintf("%s/%s.json", c.ConfigDirectory, DefaultCredentialsFile))

	err := ioutil.WriteFile(defaultCredentialsFile, file, 0600)
	if err != nil {
		return err
	}

	return nil
}

// RemoveProfile removes an existing profile from the credentials file.
func (c *Credentials) RemoveProfile(profileName string) error {

	if !c.profileExists(profileName) {
		return fmt.Errorf("Profile with name %s not found", profileName)
	}

	delete(c.Profiles, profileName)

	file, _ := json.MarshalIndent(c.Profiles, "", "  ")

	defaultCredentialsFile := os.ExpandEnv(fmt.Sprintf("%s/%s.json", c.ConfigDirectory, DefaultCredentialsFile))

	err := ioutil.WriteFile(defaultCredentialsFile, file, 0600)
	if err != nil {
		return err
	}

	return nil
}

// SetDefaultProfile modifies the profile name to use by default.
func (c *Credentials) SetDefaultProfile(profileName string) error {

	if !c.profileExists(profileName) {
		return fmt.Errorf("Profile with name %s not found", profileName)
	}

	if c.ConfigDirectory == "" {
		return fmt.Errorf("credential ConfigDirectory is empty: %s", c.ConfigDirectory)
	}

	defaultProfileFileName := os.ExpandEnv(fmt.Sprintf("%s/%s.json", c.ConfigDirectory, DefaultProfileFile))

	_, err := os.Stat(defaultProfileFileName)
	if err != nil {
		// TODO perhaps the create of the directory should be handled by the config package
		_, err := os.Create(defaultProfileFileName)
		if err != nil {
			return fmt.Errorf("error creating file %s: %s", defaultProfileFileName, err)
		}
	}

	content := fmt.Sprintf("\"%s\"", profileName)

	err = ioutil.WriteFile(defaultProfileFileName, []byte(content), 0)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %s", defaultProfileFileName, err)
	}

	return nil
}

// List outputs a list of all the configured Credentials
func (c *Credentials) List() {
	// Console colors
	color.Set(color.Bold)
	defer color.Unset()
	colorMuted := color.New(color.FgHiBlack).SprintFunc()

	nameLen := 4     // Name
	keyLen := 8      // <hidden>
	adminKeyLen := 8 // <hidden>
	regionLen := 6   // Region

	// Find lengths for pretty printing
	for k, v := range c.Profiles {
		x := len(k)
		if x > nameLen {
			nameLen = x
		}

		z := len(v.Region)
		if z > regionLen {
			regionLen = z
		}

		if showKeys {
			y := len(v.PersonalAPIKey)
			if y > keyLen {
				keyLen = y
			}
		}
	}

	nameLen += len(defaultProfileString)

	format := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds  %%-%ds\n", nameLen, regionLen, keyLen, adminKeyLen)
	fmt.Printf(format, "Name", "Region", "API key", "Admin API key")
	fmt.Printf(format, strings.Repeat("-", nameLen), strings.Repeat("-", regionLen), strings.Repeat("-", keyLen), strings.Repeat("-", adminKeyLen))
	// Print them out
	for k, v := range c.Profiles {
		name := k
		if k == c.DefaultProfile {
			name += colorMuted(defaultProfileString)
		}
		key := colorMuted(hiddenKeyString)
		if showKeys {
			key = v.PersonalAPIKey
		}

		adminKey := colorMuted(hiddenKeyString)
		if showKeys {
			adminKey = v.AdminAPIKey
		}
		fmt.Printf(format, name, v.Region, key, adminKey)
	}
	fmt.Println("")
}
