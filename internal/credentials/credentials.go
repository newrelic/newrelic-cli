package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/newrelic/newrelic-cli/internal/config"
	log "github.com/sirupsen/logrus"
)

// DefaultCredentialsFile is the default place to load profiles from
const DefaultCredentialsFile = "credentials"

const defaultConfigType = "json"

const (
	defaultProfileString = " (default)"
	hiddenKeyString      = "<hidden>"
)

// Credentials is the metadata around all configured profiles
type Credentials struct {
	DefaultProfile  string
	Profiles        map[string]Profile
	ConfigDirectory string
}

// LoadCredentials loads the list of profiles
func LoadCredentials(configDir string) (*Credentials, error) {
	log.Debug("loading credentials file")

	if configDir == "" {
		configDir = os.ExpandEnv(config.DefaultConfigDirectory)
	} else {
		configDir = os.ExpandEnv(configDir)
	}

	// Load profiles
	profiles, err := LoadProfiles(configDir)
	if err != nil {
		return &Credentials{}, err
	}

	defaultProfile, err := LoadDefaultProfile(configDir)
	if err != nil {
		log.Warnf("No default configuration set.  Please set one with ... TBD")
	}

	creds := &Credentials{
		ConfigDirectory: configDir,
		Profiles:        *profiles,
		DefaultProfile:  defaultProfile,
	}

	return creds, nil
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

// Set configures a profile config value.
func (c *Credentials) Set(profileName string, key string, value string) error {
	err := c.set(profileName, key, value)
	if err != nil {
		return err
	}

	// TODO
	// renderer.Set(key, value)
	return nil
}

func (c *Credentials) set(profileName string, key string, value interface{}) error {
	cfgViper, err := readCredentials(c.ConfigDirectory)
	if err != nil {
		return err
	}

	setPath := fmt.Sprintf("%s.%s", profileName, key)
	cfgViper.Set(setPath, value)

	allCreds, err := unmarshalProfiles(cfgViper)
	if err != nil {
		return err
	}

	c.Profiles = *allCreds
	c.validate()

	cfgViper.WriteConfig()

	return nil
}

func (c *Credentials) validate() error {
	for _, p := range c.Profiles {
		p.validate()
	}

	if c.ConfigDirectory == "" {
		return fmt.Errorf("Credentials.ConfigDirectory should not be empty")
	}

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
		return fmt.Errorf("Error writing file %s: %s", file, err)
	}

	// If we only have the single profile we've just added, let's make it the default.
	if len(c.Profiles) == 1 {
		err = c.SetDefaultProfile(profileName)
		if err != nil {
			return fmt.Errorf("error setting the default profile: %s", err)
		}
	}

	// TODO Using set() here results in key strings that are lowercased by the
	// viper lib pending https://github.com/spf13/viper/issues/373

	// if region != "" {
	// 	err := c.set(profileName, "region", region)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	//
	// if apiKey != "" {
	// 	err := c.set(profileName, "apiKey", apiKey)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	//
	// if adminAPIKey != "" {
	// 	err := c.set(profileName, "adminAPIKey", adminAPIKey)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

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

	c.DefaultProfile = profileName

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
