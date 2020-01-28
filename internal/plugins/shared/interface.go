package shared

// CommandFlag represents a command flag.
type CommandFlag struct {
	Name      string
	Options   []string
	Prompt    string
	Required  bool
	Shorthand string
	Type      string
	Usage     string
}

// CommandDefinition represents the definition of a CLI subcommand.
type CommandDefinition struct {
	PluginCommand string
	PluginArgs    []string
	Use           string
	Short         string
	Long          string
	Flags         []*CommandFlag
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
