package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
)

const (
	// DefaultCredentialsFile is the default place to load profiles from
	DefaultCredentialsFile = "credentials"

	defaultConfigType = "json"
	//defaultProfileString = " (default)"
	//hiddenKeyString      = "<hidden>"
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
func (c *Credentials) AddProfile(profileName string, p Profile) error {
	var err error

	if c.profileExists(profileName) {
		return fmt.Errorf("profile with name %s already exists", profileName)
	}

	// Case fold the region
	p.Region = strings.ToUpper(p.Region)

	c.Profiles[profileName] = p

	file, _ := json.MarshalIndent(c.Profiles, "", "  ")
	defaultCredentialsFile := os.ExpandEnv(fmt.Sprintf("%s/%s.json", c.ConfigDirectory, DefaultCredentialsFile))

	if _, err = os.Stat(c.ConfigDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(c.ConfigDirectory, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(defaultCredentialsFile, file, 0600)
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
}

// The order of fields in this struct dictates the ordering of the output table.
// type profileList struct {
// 	Name              string
// 	AccountID         int
// 	Region            string
// 	APIKey            string
// 	LicenseKey        string
// 	InsightsInsertKey string
// }
