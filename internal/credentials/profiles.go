package credentials

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
