package profile

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
func (p *Credentials) Default() *Profile {
	if p.DefaultProfile != "" {
		if val, ok := p.Profiles[p.DefaultProfile]; ok {
			return &val
		}
	}

	return nil
}
