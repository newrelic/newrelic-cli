package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
)

const (
	defaultProfileString = " (default)"
	hiddenKeyString      = "<hidden>"
)

// Profile contains data of a single profile
type Profile struct {
	AdminAPIKey    string `mapstructure:"adminAPIKey" json:"adminAPIKey"` // For accessing New Relic (Rest v2)
	PersonalAPIKey string `mapstructure:"apiKey" json:"apiKey"`           // For accessing New Relic GraphQL resources
	Region         string `mapstructure:"region" json:"region"`           // Region to use for New Relic resources
}

// Credentials is the metadata around all configured profiles
type Credentials struct {
	DefaultProfile  string
	Profiles        map[string]Profile
	ConfigDirectory string
}

// LoadCredentials loads the list of profiles
func LoadCredentials(configDirectory string, envPrefix string) (*Credentials, error) {
	var err error

	// Load profiles
	creds, err := Load(configDirectory, envPrefix)
	if err != nil {
		return &Credentials{}, err
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
func (c *Credentials) SetDefaultProfile(name string) error {

	profileExists := func(n string) bool {
		for k := range c.Profiles {
			if k == n {
				return true
			}
		}

		return false
	}(name)

	if !profileExists {
		return fmt.Errorf("the specified profile does not exist: %s", name)
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

	content := fmt.Sprintf("\"%s\"", name)

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
