package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
)

const (
	// DefaultCredentialsFile is the default place to load profiles from
	DefaultCredentialsFile = "credentials"

	defaultConfigType    = "json"
	defaultProfileString = " (default)"
	hiddenKeyString      = "<hidden>"
)

// Credentials is the metadata around all configured profiles
type Credentials struct {
	DefaultProfile  string
	Profiles        map[string]Profile
	ConfigDirectory string
}

// LoadCredentials loads the current CLI credentials from disk.
func LoadCredentials(configDir string) (*Credentials, error) {
	log.Debug("loading credentials file")

	if configDir == "" {
		configDir = config.DefaultConfigDirectory
	} else {
		configDir = os.ExpandEnv(configDir)
	}

	creds := &Credentials{
		ConfigDirectory: configDir,
	}

	profiles, err := LoadProfiles(configDir)
	if err != nil {
		log.Debugf("no credential profiles: see newrelic profiles --help")
	}

	defaultProfile, err := LoadDefaultProfile(configDir)
	if err != nil {
		log.Debugf("no default profile set: see newrelic profiles --help")
	}

	creds.Profiles = *profiles
	creds.DefaultProfile = defaultProfile

	return creds, nil
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
func (c *Credentials) AddProfile(profileName, region, apiKey string) error {

	if c.profileExists(profileName) {
		return fmt.Errorf("profile with name %s already exists", profileName)
	}

	p := Profile{
		Region: strings.ToLower(region),
		APIKey: apiKey,
	}

	c.Profiles[profileName] = p

	file, _ := json.MarshalIndent(c.Profiles, "", "  ")
	defaultCredentialsFile := os.ExpandEnv(fmt.Sprintf("%s/%s.json", c.ConfigDirectory, DefaultCredentialsFile))

	if _, err := os.Stat(c.ConfigDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(c.ConfigDirectory, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err := ioutil.WriteFile(defaultCredentialsFile, file, 0600)
	if err != nil {
		return err
	}

	return nil
}

// RemoveProfile removes an existing profile from the credentials file.
func (c *Credentials) RemoveProfile(profileName string) error {
	if !c.profileExists(profileName) {
		return fmt.Errorf("profile with name %s not found", profileName)
	}

	delete(c.Profiles, profileName)

	file, _ := json.MarshalIndent(c.Profiles, "", "  ")
	defaultCredentialsFile := os.ExpandEnv(fmt.Sprintf("%s/%s.json", c.ConfigDirectory, DefaultCredentialsFile))

	err := ioutil.WriteFile(defaultCredentialsFile, file, 0600)
	if err != nil {
		return err
	}

	if profileName == c.DefaultProfile {
		c.DefaultProfile = ""
		defaultProfileFileName := os.ExpandEnv(fmt.Sprintf("%s/%s.json", c.ConfigDirectory, DefaultProfileFile))

		err := os.Remove(defaultProfileFileName)
		if err != nil {
			return err
		}
	}

	log.Infof("profile %s has been removed", profileName)

	return nil
}

// SetDefaultProfile modifies the profile name to use by default.
func (c *Credentials) SetDefaultProfile(profileName string) error {
	if !c.profileExists(profileName) {
		return fmt.Errorf("profile with name %s not found", profileName)
	}

	if c.ConfigDirectory == "" {
		return fmt.Errorf("credential ConfigDirectory is empty: %s", c.ConfigDirectory)
	}

	c.DefaultProfile = profileName

	defaultProfileFileName := os.ExpandEnv(fmt.Sprintf("%s/%s.json", c.ConfigDirectory, DefaultProfileFile))
	content := fmt.Sprintf("\"%s\"", profileName)

	err := ioutil.WriteFile(defaultProfileFileName, []byte(content), 0600)
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

	nameLen := 4   // Name
	keyLen := 8    // <hidden>
	regionLen := 6 // Region

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
			y := len(v.APIKey)
			if y > keyLen {
				keyLen = y
			}
		}
	}

	nameLen += len(defaultProfileString)

	format := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds\n", nameLen, regionLen, keyLen)
	fmt.Printf(format, "Name", "Region", "API key")
	fmt.Printf(format, strings.Repeat("-", nameLen), strings.Repeat("-", regionLen), strings.Repeat("-", keyLen))
	// Print them out
	for k, v := range c.Profiles {
		name := k
		if k == c.DefaultProfile {
			name += colorMuted(defaultProfileString)
		}
		key := colorMuted(hiddenKeyString)
		if showKeys {
			key = v.APIKey
		}

		fmt.Printf(format, name, v.Region, key)
	}
	fmt.Println("")
}
