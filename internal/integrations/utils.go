package integrations

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	metricArg      = "metrics"
	inventoryArg   = "inventory"
	eventsArg      = "events"
	exeSuffix      = ".exe"
	prefixArg      = "--"
	prefixArgShort = "-"
)

// MigrateV3toV4Result represents the result of a migration
type MigrateV3toV4Result struct {
	MigrateV3toV4Result string `json:"migrateV3toV4Result"`
}

func migrateV3toV4(pathConfiguration string, pathDefinition string, pathOutput string) string {
	// Reading old Definition file
	v3Definition := Plugin{}
	v3DefinitionBytes, err := readAndUnmarshallConfig(pathDefinition, &v3Definition)
	if err != nil {
		log.Fatal(fmt.Errorf("error reading old config definition: %w", err))
	}

	// Reading old Configuration file
	v3Configuration := PluginInstanceWrapper{}
	v3ConfigurationBytes, err := readAndUnmarshallConfig(pathConfiguration, &v3Configuration)
	if err != nil {
		log.Fatal(fmt.Errorf("error reading old config configuration: %w", err))
	}

	// Populating new config
	v4config, err := populateV4Config(v3Definition, v3Configuration)
	if err != nil {
		log.Fatal(fmt.Errorf("error populating new config: %w", err))
	}

	// Writing output
	err = writeOutput(pathOutput, v4config, v3DefinitionBytes, v3ConfigurationBytes)
	if err != nil {
		log.Fatal(fmt.Errorf("error writing output: %w", err))
	}

	return fmt.Sprintf("The config has been migrated and placed in: %s", pathOutput)
}

func readAndUnmarshallConfig(path string, out interface{}) ([]byte, error) {
	// Reading old definition file
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s, %w", path, err)
	}

	err = yaml.Unmarshal(bytes, out)
	if err != nil {
		return nil, fmt.Errorf("unmashalling %s, %w", path, err)
	}

	return bytes, nil
}

func populateV4Config(v3Definition Plugin, v3Configuration PluginInstanceWrapper) (*v4, error) {
	if v3Configuration.IntegrationName != v3Definition.Name {
		return nil, fmt.Errorf("IntegrationName != Name: %s!=%s", v3Configuration.IntegrationName, v3Definition.Name)
	}

	// The field os does not have currently a simple way to be migrated
	if v3Definition.OS != "" {
		log.Debugf("The old definitions had a os directive, %s. Usually it is not needed, use `when` field otherwhise", v3Definition.OS)
	}

	v4Config := &v4{}
	for commandName, pluginV1Command := range v3Definition.Commands {
		for _, pluginV1Instance := range v3Configuration.Instances {
			if commandName == pluginV1Instance.Command {
				integrationInstance := populateConfigEntry(pluginV1Instance, pluginV1Command)
				v4Config.Integrations = append(v4Config.Integrations, integrationInstance)
			}
		}
	}

	return v4Config, nil
}

func populateConfigEntry(pluginV1Instance *PluginV1Instance, pluginV1Command *PluginV1Command) ConfigEntry {
	configEntry := ConfigEntry{}
	if len(pluginV1Command.Command) == 0 {
		return configEntry
	}

	executable := pluginV1Command.Command[0]
	binaryName := filepath.Base(executable)
	configEntry.InstanceName = strings.TrimSuffix(binaryName, exeSuffix)
	configEntry.Interval = fmt.Sprintf("%ds", pluginV1Command.Interval)
	configEntry.Labels = pluginV1Instance.Labels
	configEntry.User = pluginV1Instance.IntegrationUser
	configEntry.InventorySource = pluginV1Command.Prefix
	configEntry.Env = map[string]string{}
	for k, v := range pluginV1Instance.Arguments {
		configEntry.Env[strings.ToUpper(k)] = v
	}

	// Please notice that this is a simplification. If it is an absolute path we are adding it to the exec
	// if is a relative path or a integration name, we are assuming it is a standard integration included into the path
	if filepath.IsAbs(executable) {
		configEntry.Exec = pluginV1Command.Command
	} else {
		buildCLIArgs(pluginV1Command, &configEntry)
	}
	return configEntry
}

func buildCLIArgs(pluginV1Command *PluginV1Command, configEntry *ConfigEntry) {
	for index, arg := range pluginV1Command.Command {
		if index == 0 {
			// the first arg in command is the binary name
			continue
		}

		sanitized := strings.TrimPrefix(arg, prefixArg)
		sanitized = strings.TrimPrefix(sanitized, prefixArgShort)
		if sanitized == metricArg || sanitized == inventoryArg || sanitized == eventsArg {
			configEntry.Env[strings.ToUpper(sanitized)] = "true"
		} else {
			configEntry.CLIArgs = append(configEntry.CLIArgs, arg)
		}
	}
}

func writeOutput(pathOutput string, v4Config *v4, definitionBytes []byte, configurationBytes []byte) error {
	if v4Config == nil {
		return fmt.Errorf("v4Config pointer is nil")
	}

	file, err := os.OpenFile(pathOutput, os.O_RDWR|os.O_CREATE|syscall.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("opening File %s, %w", pathOutput, err)
	}
	defer file.Close()

	err = writeV4Config(v4Config, file)
	if err != nil {
		return fmt.Errorf("writing v4 config, %w", err)
	}

	err = writeTextAsComment(file, string(definitionBytes))
	if err != nil {
		return fmt.Errorf("adding old definition as comment, %w", err)
	}

	err = writeTextAsComment(file, string(configurationBytes))
	if err != nil {
		return fmt.Errorf("adding old configuration as comment, %w", err)
	}

	return nil
}

func writeV4Config(v4Config *v4, file *os.File) error {
	// see https://github.com/go-yaml/yaml/commit/7649d4548cb53a614db133b2a8ac1f31859dda8c
	yaml.FutureLineWrap()
	v4ConfigBytes, err := yaml.Marshal(*v4Config)
	if err != nil {
		return fmt.Errorf("marshallig v4Config, %w", err)
	}

	_, err = file.Write(v4ConfigBytes)
	if err != nil {
		return fmt.Errorf("writing v4ConfigBytes, %w", err)
	}
	return nil
}

func writeTextAsComment(file *os.File, text string) error {
	fileCommented := strings.ReplaceAll(text, "\n", "\n## ")
	fileCommented = "\n\n## " + fileCommented

	_, err := file.Write([]byte(fileCommented))
	if err != nil {
		return fmt.Errorf("writing text as comment: %w", err)
	}
	return nil
}
