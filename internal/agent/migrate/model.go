package migrate

import "time"

// This code is taken directly from the infrastructure agent
// To avoid an extra dependency it has been copied/pasted

// V3 Configuration Model
// https://github.com/newrelic/infrastructure-agent/blob/e25f6a3bfb4c637f6222663d07dbb34afee956a2/pkg/integrations/legacy/types.go#L35

// Plugin represents a single plugin, with all associated metadata
type Plugin struct {
	Name            string                      `yaml:"name"`             // Name of the plugin (required)
	Description     string                      `yaml:"description"`      // A short plugin description (optional)
	Commands        map[string]*PluginV1Command `yaml:"commands"`         // Map of commands for v1 plugins
	OS              string                      `yaml:"os"`               // OS (or comma-separated list of OSes) supported for the plugin
	ProtocolVersion int                         `yaml:"protocol_version"` // Protocol version (0 == original version)
}

type PluginV1Command struct {
	Command  []string `yaml:"command"`  // Command to execute, run from the plugin's directory.
	Prefix   string   `yaml:"prefix"`   // "Plugin path" for inventory data produced by the plugin. Not applicable for event sources.
	Interval int      `yaml:"interval"` // Number of seconds to wait between invocations of the source.
}

type PluginInstanceWrapper struct {
	IntegrationName string              `yaml:"integration_name"`
	Instances       []*PluginV1Instance `yaml:"instances"`
}

type PluginV1Instance struct {
	Name            string            `yaml:"name"`
	Command         string            `yaml:"command"`
	Arguments       map[string]string `yaml:"arguments"`
	Labels          map[string]string `yaml:"labels"`
	IntegrationUser string            `yaml:"integration_user"`
}

// V4 Configuration Model
// https://github.com/newrelic/infrastructure-agent/blob/e25f6a3bfb4c637f6222663d07dbb34afee956a2/pkg/integrations/v4/config/config.go#L15

// YAML stores the information from a single V4 integrations file
type v4 struct {
	Integrations []ConfigEntry `yaml:"integrations"`
}

// ConfigEntry holds an integrations YAML configuration entry. It may define multiple types of tasks
type ConfigEntry struct {
	InstanceName    string            `yaml:"name,omitempty" json:"name"`         // integration instance name
	CLIArgs         []string          `yaml:"cli_args,omitempty" json:"cli_args"` // optional when executable is deduced by "name" instead of "exec"
	Exec            ShlexOpt          `yaml:"exec,omitempty" json:"exec"`         // it may be a CLI string or a YAML array
	Env             map[string]string `yaml:"env,omitempty" json:"env"`           // User-defined environment variables
	Interval        string            `yaml:"interval,omitempty" json:"interval"` // User-defined interval string (duration notation)
	Timeout         *time.Duration    `yaml:"timeout,omitempty" json:"timeout"`
	User            string            `yaml:"integration_user,omitempty" json:"integration_user"`
	WorkDir         string            `yaml:"working_dir,omitempty" json:"working_dir"`
	Labels          map[string]string `yaml:"labels,omitempty" json:"labels"`
	When            EnableConditions  `yaml:"when,omitempty" json:"when"`
	InventorySource string            `yaml:"inventory_source,omitempty" json:"inventory_source"`
}

// ShlexOpt is a wrapper around []string so we can use go-shlex for shell tokenizing
type ShlexOpt []string

// EnableConditions condition the execution of an integration to the trueness of ALL the conditions
type EnableConditions struct {
	// Feature allows enabling/disabling the OHI via agent cfg "feature" or cmd-channel Feature Flag
	Feature string `yaml:"feature,omitempty"`
	// FileExists conditions the execution of the OHI only if the given file path exists
	FileExists string `yaml:"file_exists,omitempty"`
	// EnvExists conditions the execution of the OHI only if the given
	// environment variables exists and match the value.
	EnvExists map[string]string `yaml:"env_exists,omitempty"`
}
