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
