package fleetcontrol

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/newrelic/newrelic-cli/internal/utils"
)

// Embed all YAML configuration files
// This allows the configurations to be bundled into the binary at compile time
// All configs are in the configs/ directory with names matching their handler files
// Example: fleet_management_create.yaml matches fleet_management_create.go
//
//go:embed configs/*.yaml
var configFS embed.FS

// CommandConfig represents the entire command configuration
// It contains all command definitions loaded from YAML files
type CommandConfig struct {
	Commands []CommandDefinition `yaml:"commands"`
}

// CommandDefinition represents a single command definition from YAML
// Each command has metadata (name, description, examples) and flag definitions
type CommandDefinition struct {
	Name    string           `yaml:"name"`    // Command name (e.g., "create", "update")
	Short   string           `yaml:"short"`   // Short one-line description
	Long    string           `yaml:"long"`    // Detailed multi-line description
	Example string           `yaml:"example"` // Usage examples
	Flags   []FlagDefinition `yaml:"flags"`   // Command flags
}

// FlagDefinition represents a single flag definition from YAML
// Flags define the command-line arguments a command accepts
type FlagDefinition struct {
	Name        string           `yaml:"name"`        // Flag name (e.g., "name", "id")
	Type        string           `yaml:"type"`        // Data type: string, int, bool, stringSlice, file
	Required    bool             `yaml:"required"`    // Whether this flag is required
	Default     interface{}      `yaml:"default"`     // Default value if not provided
	Description string           `yaml:"description"` // Help text for this flag
	Validation  *ValidationRules `yaml:"validation,omitempty"` // Optional validation rules
}

// ValidationRules defines validation constraints for a flag
// These rules are enforced by the framework before the command handler runs
type ValidationRules struct {
	AllowedValues   []string `yaml:"allowed_values,omitempty"`   // List of acceptable values
	CaseInsensitive bool     `yaml:"case_insensitive,omitempty"` // Whether to ignore case when validating
}

// FlagValues provides dynamic, validated access to flag values
// This is passed to command handlers and provides type-safe flag access
type FlagValues struct {
	cmd        *cobra.Command              // The cobra command instance
	flagDefs   map[string]*FlagDefinition  // Map of flag name to definition
	cachedVals map[string]interface{}      // Cache of retrieved values
}

// newFlagValues creates a new FlagValues instance for a command
// This is called internally by the framework when executing a command
//
// Parameters:
//   - cmd: The cobra command being executed
//   - flagDefs: The flag definitions for this command
//
// Returns:
//   - A new FlagValues instance with caching enabled
func newFlagValues(cmd *cobra.Command, flagDefs []FlagDefinition) *FlagValues {
	defsMap := make(map[string]*FlagDefinition)
	for i := range flagDefs {
		defsMap[flagDefs[i].Name] = &flagDefs[i]
	}
	return &FlagValues{
		cmd:        cmd,
		flagDefs:   defsMap,
		cachedVals: make(map[string]interface{}),
	}
}

// GetString retrieves a validated string flag value by name
// Values are cached after first retrieval for performance
//
// Parameters:
//   - name: The flag name
//
// Returns:
//   - The string value of the flag
func (fv *FlagValues) GetString(name string) string {
	if cached, ok := fv.cachedVals[name]; ok {
		return cached.(string)
	}
	val, _ := fv.cmd.Flags().GetString(name)
	fv.cachedVals[name] = val
	return val
}

// GetStringSlice retrieves a validated string slice flag value by name
// Values are cached after first retrieval for performance
//
// Parameters:
//   - name: The flag name
//
// Returns:
//   - The string slice value of the flag
func (fv *FlagValues) GetStringSlice(name string) []string {
	if cached, ok := fv.cachedVals[name]; ok {
		return cached.([]string)
	}
	val, _ := fv.cmd.Flags().GetStringSlice(name)
	fv.cachedVals[name] = val
	return val
}

// GetInt retrieves a validated int flag value by name
// Values are cached after first retrieval for performance
//
// Parameters:
//   - name: The flag name
//
// Returns:
//   - The int value of the flag
func (fv *FlagValues) GetInt(name string) int {
	if cached, ok := fv.cachedVals[name]; ok {
		return cached.(int)
	}
	val, _ := fv.cmd.Flags().GetInt(name)
	fv.cachedVals[name] = val
	return val
}

// GetBool retrieves a validated bool flag value by name
// Values are cached after first retrieval for performance
//
// Parameters:
//   - name: The flag name
//
// Returns:
//   - The bool value of the flag
func (fv *FlagValues) GetBool(name string) bool {
	if cached, ok := fv.cachedVals[name]; ok {
		return cached.(bool)
	}
	val, _ := fv.cmd.Flags().GetBool(name)
	fv.cachedVals[name] = val
	return val
}

// GetFile retrieves file content from a file path (strict file reading only)
// This method requires a valid file path and errors if the file doesn't exist.
// For inline content, use separate -content flags with GetString instead.
// Values are cached after first retrieval for performance
//
// Parameters:
//   - name: The flag name
//
// Returns:
//   - The file content
//   - Error if file doesn't exist or reading fails
//
// Example:
//   // User must provide a valid file path:
//   // --configuration-file-path ./config.json  (reads file)
//   // For inline content, use --configuration-content flag instead
func (fv *FlagValues) GetFile(name string) (string, error) {
	if cached, ok := fv.cachedVals[name]; ok {
		return cached.(string), nil
	}

	val, _ := fv.cmd.Flags().GetString(name)

	// If the flag is empty (not provided), return empty string without error
	if val == "" {
		return "", nil
	}

	// Read the file - error if it doesn't exist
	content, err := os.ReadFile(val)
	if err != nil {
		return "", fmt.Errorf("failed to read file '%s': %w", val, err)
	}

	result := string(content)
	fv.cachedVals[name] = result
	return result, nil
}

// Has checks if a flag was explicitly provided by the user
// This is useful for distinguishing between a flag being omitted vs set to default value
//
// Parameters:
//   - name: The flag name
//
// Returns:
//   - true if the flag was provided, false otherwise
func (fv *FlagValues) Has(name string) bool {
	return fv.cmd.Flags().Changed(name)
}

// validateFlags validates all flags against their YAML-defined validation rules
// This is called automatically by the framework before the command handler runs
//
// Returns:
//   - Error if any flag fails validation
func (fv *FlagValues) validateFlags() error {
	for name, def := range fv.flagDefs {
		// Skip validation for flags that weren't provided
		if !fv.cmd.Flags().Changed(name) {
			continue
		}

		// Get raw value based on type
		var rawVal string
		switch def.Type {
		case "string", "file":
			rawVal, _ = fv.cmd.Flags().GetString(name)
		case "stringSlice":
			// For slices, validate each element
			vals, _ := fv.cmd.Flags().GetStringSlice(name)
			for _, v := range vals {
				if err := validateValue(name, v, def.Validation); err != nil {
					return err
				}
			}
			continue
		default:
			// Skip validation for non-string types (int, bool)
			continue
		}

		if err := validateValue(name, rawVal, def.Validation); err != nil {
			return err
		}
	}
	return nil
}

// validateValue validates a single value against validation rules from YAML
// This enforces the allowed_values constraint defined in the YAML configuration
//
// Parameters:
//   - flagName: The flag name (for error messages)
//   - value: The value to validate
//   - rules: The validation rules from YAML
//
// Returns:
//   - Error if validation fails
func validateValue(flagName, value string, rules *ValidationRules) error {
	if rules == nil {
		return nil
	}

	// Validate allowed values
	if len(rules.AllowedValues) > 0 {
		found := false
		compareValue := value
		if rules.CaseInsensitive {
			compareValue = strings.ToUpper(value)
		}

		for _, allowed := range rules.AllowedValues {
			compareAllowed := allowed
			if rules.CaseInsensitive {
				compareAllowed = strings.ToUpper(allowed)
			}
			if compareValue == compareAllowed {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("invalid value '%s' for flag --%s: must be one of [%s]",
				value, flagName, strings.Join(rules.AllowedValues, ", "))
		}
	}

	return nil
}

// CommandHandler is a function that handles command execution
// All command handlers must match this signature
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Positional arguments (rarely used in this CLI)
//   - flags: Validated flag values accessible through typed accessors
//
// Returns:
//   - Error if command execution fails
type CommandHandler func(cmd *cobra.Command, args []string, flags *FlagValues) error

// LoadCommandConfig loads and parses all command configuration YAML files
// This reads from the embedded configs directory
// Config file names match their handler file names (e.g., fleet_management_create.yaml)
//
// Returns:
//   - CommandConfig containing all command definitions
//   - Error if loading or parsing fails
func LoadCommandConfig() (*CommandConfig, error) {
	var allCommands []CommandDefinition

	// Read all YAML files from the configs directory
	entries, err := configFS.ReadDir("configs")
	if err != nil {
		return nil, fmt.Errorf("failed to read configs directory: %w", err)
	}

	for _, entry := range entries {
		// Skip directories and non-YAML files
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		// Read the YAML file
		configPath := "configs/" + entry.Name()
		data, err := configFS.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", entry.Name(), err)
		}

		// Parse the YAML into a CommandDefinition
		var cmdDef CommandDefinition
		if err := yaml.Unmarshal(data, &cmdDef); err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", entry.Name(), err)
		}

		allCommands = append(allCommands, cmdDef)
	}

	return &CommandConfig{Commands: allCommands}, nil
}

// BuildCommand creates a cobra command from a CommandDefinition
// This sets up all flags, validation, and wires up the handler
//
// Parameters:
//   - def: The command definition from YAML
//   - handler: The function that implements the command logic
//
// Returns:
//   - A fully configured cobra.Command
func BuildCommand(def CommandDefinition, handler CommandHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:     def.Name,
		Short:   def.Short,
		Long:    def.Long,
		Example: def.Example,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create flag values accessor
			flags := newFlagValues(cmd, def.Flags)

			// Validate all flags before executing handler
			// This enforces YAML-defined validation rules
			if err := flags.validateFlags(); err != nil {
				return err
			}

			// Call the command handler with validated flags
			return handler(cmd, args, flags)
		},
	}

	// Register all flags with cobra
	for _, flag := range def.Flags {
		switch flag.Type {
		case "string", "file":
			// Both string and file types are registered as string flags
			// The difference is in how FlagValues.GetFile() retrieves them
			defaultVal := ""
			if flag.Default != nil {
				defaultVal = flag.Default.(string)
			}
			cmd.Flags().String(flag.Name, defaultVal, flag.Description)

		case "stringSlice":
			var defaultVal []string
			if flag.Default != nil {
				defaultVal = flag.Default.([]string)
			}
			cmd.Flags().StringSlice(flag.Name, defaultVal, flag.Description)

		case "int":
			defaultVal := 0
			if flag.Default != nil {
				switch v := flag.Default.(type) {
				case int:
					defaultVal = v
				case float64:
					// YAML numbers are parsed as float64
					defaultVal = int(v)
				}
			}
			cmd.Flags().Int(flag.Name, defaultVal, flag.Description)

		case "bool":
			defaultVal := false
			if flag.Default != nil {
				defaultVal = flag.Default.(bool)
			}
			cmd.Flags().Bool(flag.Name, defaultVal, flag.Description)
		}

		// Mark as required if specified in YAML
		if flag.Required {
			utils.LogIfError(cmd.MarkFlagRequired(flag.Name))
		}
	}

	return cmd
}

// GetCommandDefinition finds a command definition by name
// This is used to look up specific commands from the loaded configuration
//
// Parameters:
//   - config: The loaded command configuration
//   - name: The command name to find
//
// Returns:
//   - The command definition, or nil if not found
func GetCommandDefinition(config *CommandConfig, name string) *CommandDefinition {
	for _, cmd := range config.Commands {
		if cmd.Name == name {
			return &cmd
		}
	}
	return nil
}
