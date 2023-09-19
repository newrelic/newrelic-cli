package validation

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type integrationType struct {
	Name string            `yaml:"name"`
	Env  map[string]string `yaml:"env"`
}

type configType struct {
	Integrations []integrationType `yaml:"integrations"`
}

/*
  Validate an integration using its <integrationName>-config.yml
  by iterating over the defined integrations and running the
  integration command with its defined environment variables.

  <integrationName>-config.yml is determined to be valid if every
  defined integration exits without error (exits with exit code 0).

  The <integrationName>-config.yml is located in the
  default configuration directory, which may vary by GOOS.
*/
func ValidateIntegration(integrationName string) (string, error) {
	configBasename := fmt.Sprintf("%s-config.yml", integrationName)

	configPath := filepath.Join(ConfigurationsDirname, configBasename)

	config, err := readConfig(configPath)

	if err != nil {
		return "", err
	}

	for _, integration := range config.Integrations {
		integrationBasename := integration.Name

		binPath := filepath.Join(IntegrationsDirname, integrationBasename)

		cmd := exec.Command(binPath)

		for k, v := range integration.Env {
			env := fmt.Sprintf("%s=%s", k, v)

			cmd.Env = append(cmd.Env, env)
		}

		var e strings.Builder

		cmd.Stderr = &e

		if err := cmd.Run(); err != nil {
			stderr := e.String()

			if stderr != "" {
				stderr := strings.TrimSpace(stderr)

				return "", fmt.Errorf("%w: %s", err, stderr)
			}

			return "", err
		}
	}

	return "", nil
}

/*
  Reads and unmarshals an <integrationName>-config.yml from the
  given configPath, returning a configType{} containing the
  defined integrations and their respective environments.

  Returns an empty configType{} and an error if:
  - The file does not exist
  - The file can not be read
  - The file cannot be unmarshalled to a configType{}
*/
func readConfig(configPath string) (configType, error) {
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return configType{}, err
	}

	b, err := ioutil.ReadFile(configPath)

	if err != nil {
		return configType{}, err
	}

	config := configType{}

	err = yaml.Unmarshal(b, &config)

	if err != nil {
		return configType{}, err
	}

	return config, nil
}
