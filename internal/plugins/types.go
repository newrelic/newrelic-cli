package plugins

// CommandFlag represents a command flag.
type CommandFlag struct {
	Name      string   `yaml:"Name,omitempty"`
	Options   []string `yaml:"Options,omitempty"`
	Prompt    string   `yaml:"Prompt,omitempty"`
	Required  bool     `yaml:"Required,omitempty"`
	Shorthand string   `yaml:"Shorthand,omitempty"`
	Type      string   `yaml:"Type,omitempty"`
	Usage     string   `yaml:"Usage,omitempty"`
}

// CommandDefinition represents the definition of a CLI subcommand.
type CommandDefinition struct {
	PluginCommand string         `yaml:"PluginCommand,omitempty"`
	PluginArgs    []string       `yaml:"PluginArgs,omitempty"`
	Use           string         `yaml:"Use,omitempty"`
	Short         string         `yaml:"Short,omitempty"`
	Long          string         `yaml:"Long,omitempty"`
	Flags         []*CommandFlag `yaml:"Flags,omitempty"`
	Interactive   bool           `yaml:"Interactive,omitempty"`
}

// Flag retrieves a specific CommandFlag by name.
func (cd *CommandDefinition) Flag(name string) *CommandFlag {
	for _, f := range cd.Flags {
		if f.Name == name {
			return f
		}
	}

	return nil
}
