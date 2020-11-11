package config

var CLIConfig *Configuration

type Configuration struct {
	LogLevel string
	Profiles map[string]Profile
}

type Profile struct {
	APIKey            string `json:"apiKey,omitempty"`            // For accessing New Relic GraphQL resources
	InsightsInsertKey string `json:"insightsInsertKey,omitempty"` // For posting custom events
	Region            string `json:"region,omitempty"`            // Region to use for New Relic resources
	AccountID         int    `json:"accountID"`                   // AccountID to use for New Relic resources
	LicenseKey        string `json:"licenseKey"`                  // License key to use for agent config and ingest
}

// Just an idea for now
func (c *Configuration) GetCurrentProfile() *Profile {
	return nil
}

func (c *Configuration) GetDefaultProfile() *Profile {
	return nil
}

func (c *Configuration) GetProfile(name string) *Profile {
	return nil
}

// func init() {
// 	cfgViper := viper.New()

// 	configDir, err := getDefaultConfigDirectory()
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	cfgViper.SetEnvPrefix(DefaultEnvPrefix)
// 	cfgViper.SetConfigName(DefaultConfigName)
// 	cfgViper.SetConfigType(DefaultConfigType)

// 	// Set the config file path
// 	cfgViper.AddConfigPath(configDir)

// 	// Read in environment variables that
// 	// match the environment prefix
// 	cfgViper.AutomaticEnv()

// 	if err := viper.ReadInConfig(); err != nil {

// 		fmt.Printf("\n\n *** init - AllSettings: %v \n\n", cfgViper.AllSettings())

// 		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
// 			log.Infof("configuration file not found: %s", configDir)
// 		} else {
// 			// Config file was found but another error was produced
// 			log.Fatal(err.Error())
// 		}
// 	}
// }
