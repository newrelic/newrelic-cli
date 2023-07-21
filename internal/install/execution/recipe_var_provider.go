package execution

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var (
	downloadURLAccessListRegex = []string{
		`(.)?download\.newrelic\.com$`,
		`nr-downloads-ohai-(staging|testing)\.s3\.amazonaws\.com$`,
	}
)

const (
	EnvNriaCustomAttributes       = "NRIA_CUSTOM_ATTRIBUTES"
	EnvNriaPassthroughEnvironment = "NRIA_PASSTHROUGH_ENVIRONMENT"
	EnvInstallCustomAttributes    = "INSTALL_CUSTOM_ATTRIBUTES"
)

type RecipeVarProvider struct{}

func NewRecipeVarProvider() *RecipeVarProvider {
	return &RecipeVarProvider{}
}

func (re *RecipeVarProvider) Prepare(m types.DiscoveryManifest, r types.OpenInstallationRecipe, assumeYes bool) (types.RecipeVars, error) {
	log.WithFields(log.Fields{
		"name": r.Name,
	}).Debug("preparing recipe")

	vars := types.RecipeVars{}

	results := []types.RecipeVars{}

	systemInfoResult := varsFromSystemInfo(m)

	profileResult, err := varsFromProfile()
	if err != nil {
		return types.RecipeVars{}, err
	}

	inputVarsResult, err := varsFromInput(r.InputVars, assumeYes)
	if err != nil {
		return types.RecipeVars{}, err
	}

	envVarsResult := varFromEnv()

	results = append(results, systemInfoResult)
	results = append(results, profileResult)
	results = append(results, types.RecipeVariables)
	results = append(results, inputVarsResult)
	results = append(results, envVarsResult)

	for _, result := range results {
		for k, v := range result {
			vars[k] = v
		}
	}

	return vars, nil
}

func varsFromProfile() (types.RecipeVars, error) {
	accountID := configAPI.GetActiveProfileString(config.AccountID)
	apiKey := configAPI.GetActiveProfileString(config.APIKey)
	region := configAPI.GetActiveProfileString(config.Region)

	vars := make(types.RecipeVars)

	vars["NEW_RELIC_LICENSE_KEY"] = os.Getenv("NEW_RELIC_LICENSE_KEY")
	vars["NEW_RELIC_ACCOUNT_ID"] = accountID
	vars["NEW_RELIC_API_KEY"] = apiKey
	vars["NEW_RELIC_REGION"] = region

	return vars, nil
}

func varsFromSystemInfo(m types.DiscoveryManifest) types.RecipeVars {
	vars := make(types.RecipeVars)

	vars["HOSTNAME"] = m.Hostname
	vars["OS"] = m.OS
	vars["PLATFORM"] = m.Platform
	vars["PLATFORM_FAMILY"] = m.PlatformFamily
	vars["PLATFORM_VERSION"] = m.PlatformVersion
	vars["KERNEL_ARCH"] = m.KernelArch
	vars["KERNEL_VERSION"] = m.KernelVersion

	return vars
}

func varsFromInput(inputVars []types.OpenInstallationRecipeInputVariable, assumeYes bool) (types.RecipeVars, error) {
	vars := make(types.RecipeVars)

	vars["NEW_RELIC_ASSUME_YES"] = fmt.Sprintf("%t", assumeYes)

	for _, envConfig := range inputVars {
		var err error
		envValue := os.Getenv(envConfig.Name)

		if envValue != "" {
			vars[envConfig.Name] = envValue
			continue
		}

		if assumeYes {
			if envConfig.Default == "" {
				return types.RecipeVars{}, fmt.Errorf("no default value for environment variable %s and none provided", envConfig.Name)
			}

			log.WithFields(log.Fields{
				"name":    envConfig.Name,
				"default": envConfig.Default,
			}).Debug("required env var not found, using default")

			envValue = envConfig.Default
		} else {
			log.WithFields(log.Fields{
				"name": envConfig.Name,
			}).Debug("required environment variable not found")

			envValue, err = varFromPrompt(envConfig)
			if err != nil {
				if err == terminal.InterruptErr {
					return types.RecipeVars{}, types.ErrInterrupt
				}

				return types.RecipeVars{}, fmt.Errorf("prompt failed: %s", err)
			}
		}

		vars[envConfig.Name] = envValue
	}

	return vars, nil
}

func varFromPrompt(envConfig types.OpenInstallationRecipeInputVariable) (string, error) {
	msg := fmt.Sprintf("value for %s required", envConfig.Name)

	if envConfig.Prompt != "" {
		msg = envConfig.Prompt
	}

	value := ""
	var prompt survey.Prompt

	if envConfig.Secret {
		prompt = &survey.Password{
			Message: msg,
		}
	} else {
		defaultValue := ""

		if envConfig.Default != "" {
			defaultValue = envConfig.Default
		}

		prompt = &survey.Input{
			Message: msg,
			Default: defaultValue,
		}
	}

	err := survey.AskOne(prompt, &value)
	if err != nil {
		return "", err
	}

	return value, nil

}

func varFromEnv() types.RecipeVars {
	vars := make(types.RecipeVars)

	downloadURL := "https://download.newrelic.com/"
	envDownloadURL := os.Getenv("NEW_RELIC_DOWNLOAD_URL")
	if envDownloadURL != "" {
		URL, err := url.Parse(envDownloadURL)
		if err == nil {
			if URL.Scheme == "https" {
				for _, regexString := range downloadURLAccessListRegex {
					var regex = regexp.MustCompile(regexString)
					if regex.MatchString(URL.Host) {
						downloadURL = envDownloadURL
						break
					}
				}
			}
		} else {
			log.Warnf("Could not parse download URL: %s, detail: %s", envDownloadURL, err.Error())
		}
	}
	vars["NEW_RELIC_DOWNLOAD_URL"] = downloadURL
	vars["NEW_RELIC_CLI_LOG_FILE_PATH"] = config.GetDefaultLogFilePath()
	vars["NR_CLI_CLUSTERNAME"] = os.Getenv("NR_CLI_CLUSTERNAME")

	customAttributes := os.Getenv(EnvNriaCustomAttributes)
	installCustomAttributes := os.Getenv(EnvInstallCustomAttributes)
	if len(customAttributes) > 0 || len(installCustomAttributes) > 0 {
		vars[EnvNriaCustomAttributes] = yamlFromJSON(EnvNriaCustomAttributes, customAttributes, strings.Split(installCustomAttributes, ","))
	}

	if len(installCustomAttributes) > 0 {
		vars["NEW_RELIC_CLI_TAGS"] = strings.Join(strings.Split(installCustomAttributes, ","), ";")
	}

	passthroughEnvironment := os.Getenv(EnvNriaPassthroughEnvironment)
	if len(passthroughEnvironment) > 0 {
		vars[EnvNriaPassthroughEnvironment] = yamlFromCommaDelimitedString(EnvNriaPassthroughEnvironment, passthroughEnvironment)
	}

	return vars
}

func yamlFromJSON(key string, jsonVal string, tags []string) string {
	if (!json.Valid([]byte(jsonVal)) || len(jsonVal) == 0) && len(tags) == 0 {
		log.Debugf("Invalid json passed in %s: %s", key, jsonVal)
		return ""
	}

	customerAttributeMap := map[string]string{}
	if json.Valid([]byte(jsonVal)) && len(jsonVal) > 0 {
		err := json.Unmarshal([]byte(jsonVal), &customerAttributeMap)
		if err != nil {
			log.Debugf("Could transform custom attributes to a map: %e", err)
		}
	}

	// get tag and override custom attributes
	tagMap := makeTagMap(tags)
	for k, v := range tagMap {
		customerAttributeMap[k] = v
	}

	mergedAttributesJSON, err := json.Marshal(customerAttributeMap)

	if err != nil {
		log.Debugf("Could not transform custom attribute map to json: %e", err)
		return ""
	}

	customAttributesJSON := fmt.Sprintf("{\"custom_attributes\": %s }", mergedAttributesJSON)
	yaml, err := yaml.JSONToYAML([]byte(customAttributesJSON))
	if err != nil {
		log.Debugf("Could not transform %s json value to yaml: %e", key, err)
		return ""
	}

	return string(yaml)
}

func yamlFromCommaDelimitedString(key string, commaDelimited string) string {
	if len(commaDelimited) == 0 {
		return ""
	}

	passthroughEnvironmentJSON := fmt.Sprintf("{\"passthrough_environment\": { \"env\": [%s] }}", commaDelimited)
	yaml, err := yaml.JSONToYAML([]byte(passthroughEnvironmentJSON))
	if err != nil {
		log.Debugf("Could not transform %s comma-delimited value to yaml: %e", key, err)
		return ""
	}
	// forcing indentation of list items without a key
	trimmed := strings.ReplaceAll(string(yaml), "\n  env:", "")
	return trimmed
}

func makeTagMap(tags []string) map[string]string {
	tagMap := make(map[string]string)
	for _, x := range tags {
		if !strings.Contains(x, ":") {
			continue
		}
		v := strings.SplitN(x, ":", 2)
		tagMap[v[0]] = v[1]
	}
	return tagMap
}
