package credentials

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const (
	defaultProfileString = " (default)"
	hiddenKeyString      = "<hidden>"
)

// Profile contains data of a single profile
type Profile struct {
	PersonalAPIKey string `mapstructure:"apiKey"` // PersonalAPIKey for accessing New Relic
	Region         string `mapstructure:"Region"` // Region to use when accessing New Relic
}

// Credentials is the metadata around all configured profiles
type Credentials struct {
	DefaultProfile string
	Profiles       map[string]Profile
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
			y := len(v.PersonalAPIKey)
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
			key = v.PersonalAPIKey
		}
		fmt.Printf(format, name, v.Region, key)
	}
	fmt.Println("")
}
