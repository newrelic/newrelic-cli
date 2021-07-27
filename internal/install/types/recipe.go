package types

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	InfraAgentRecipeName = "infrastructure-agent-installer"
	LoggingRecipeName    = "logs-integration"
)

var (
	RecipeVariables = map[string]string{}
)

// RecipeVars is used to pass dynamic data to recipes and go-task.
type RecipeVars map[string]string

// The API response returns OpenInstallationRecipe.Install as a string.
// When specifying a recipe path, OpenInstallationRecipe.Install is a map[interface{}]interface{}.
// For this reason we need a custom unmarshal method for YAML.
func (r *OpenInstallationRecipe) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var recipe map[string]interface{}
	err := unmarshal(&recipe)
	if err != nil {
		return err
	}

	if v, ok := recipe["dependencies"]; ok {
		r.Dependencies = interfaceSliceToStringSlice(v.([]interface{}))
	}

	r.Description = toStringByFieldName("description", recipe)
	r.DisplayName = toStringByFieldName("displayName", recipe)
	r.File = toStringByFieldName("file", recipe)
	r.ID = toStringByFieldName("id", recipe)
	r.InputVars = expandInputVars(recipe)

	installAsString, err := expandInstalllMapToString(recipe)
	if err != nil {
		return err
	}
	r.Install = installAsString

	r.InstallTargets = expandInstallTargets(recipe)

	if v, ok := recipe["keywords"]; ok {
		r.Keywords = interfaceSliceToStringSlice(v.([]interface{}))
	}

	r.LogMatch = expandLogMatch(recipe)
	r.Name = toStringByFieldName("name", recipe)
	r.PostInstall = expandPostInstall(recipe)
	r.PreInstall = expandPreInstall(recipe)

	if v, ok := recipe["processMatch"]; ok {
		r.ProcessMatch = interfaceSliceToStringSlice(v.([]interface{}))
	}

	r.Quickstarts = expandQuickStarts(recipe)
	r.ObservabilityPacks = expandObservabilityPacks(recipe)

	r.Repository = toStringByFieldName("repository", recipe)

	if v, ok := recipe["stability"]; ok {
		r.Stability = OpenInstallationStability(v.(string))
	}

	r.SuccessLinkConfig = expandSuccessLinkConfig(recipe)

	// DEPRECATED: Use `validation` parameter instead
	if v, ok := recipe["validationNrql"]; ok {
		r.ValidationNRQL = NRQL(v.(string))
	}

	if v, ok := recipe["validationUrl"]; ok {
		r.ValidationURL = v.(string)
	}

	return err
}

func (r *OpenInstallationRecipe) ToShortDisplayString() string {
	output := r.Name
	targets := ""
	for _, target := range r.InstallTargets {
		s := getInstallTargetAsString(target)
		if targets != "" {
			targets = fmt.Sprintf("%s-%s", targets, s)
		} else {
			targets = s
		}
	}
	if targets != "" {
		output = fmt.Sprintf("%s (%s)", output, targets)
	}
	return output
}

func getInstallTargetAsString(target OpenInstallationRecipeInstallTarget) string {
	output := string(target.Os)
	if target.Platform != "" {
		output = fmt.Sprintf("%s/%s", output, target.Platform)
	}
	if target.PlatformVersion != "" {
		output = fmt.Sprintf("%s/%s", output, target.PlatformVersion)
	}
	if target.KernelArch != "" {
		output = fmt.Sprintf("%s/%s", output, target.KernelArch)
	}
	return strings.ToLower(output)
}

func expandObservabilityPacks(recipe map[string]interface{}) []OpenInstallationObservabilityPackFilter {
	v, ok := recipe["observabilityPacks"]
	if !ok {
		return []OpenInstallationObservabilityPackFilter{}
	}

	dataIn := v.([]interface{})
	dataOut := make([]OpenInstallationObservabilityPackFilter, len(dataIn))
	dataz := make([]map[string]interface{}, len(dataIn))

	for i, vv := range dataIn {
		vvv := vv.(map[interface{}]interface{})
		varr := map[string]interface{}{}

		for k, v := range vvv {
			varr[k.(string)] = v
			dataz[i] = varr
		}
	}

	for i, v := range dataz {
		vOut := OpenInstallationObservabilityPackFilter{
			Name: toStringByFieldName("name", v),
		}
		if v, ok := v["level"]; ok {
			vOut.Level = OpenInstallationObservabilityPackLevel(v.(string))
		}

		dataOut[i] = vOut
	}

	return dataOut
}

func expandQuickStarts(recipe map[string]interface{}) OpenInstallationQuickstartsFilter {
	v, ok := recipe["quickstarts"]
	if !ok {
		return OpenInstallationQuickstartsFilter{}
	}

	dataIn := v.(map[interface{}]interface{})
	reData := map[string]interface{}{}
	for k, v := range dataIn {
		reData[k.(string)] = v
	}

	dataOut := OpenInstallationQuickstartsFilter{
		Name:       toStringByFieldName("name", reData),
		EntityType: expandEntityType(reData),
	}

	if v, ok := reData["category"]; ok {
		dataOut.Category = OpenInstallationCategory(v.(string))
	}

	return dataOut
}

func expandEntityType(data map[string]interface{}) OpenInstallationQuickstartEntityType {
	v, ok := data["entityType"]
	if !ok {
		return OpenInstallationQuickstartEntityType{}
	}

	dataIn := v.(map[interface{}]interface{})
	reData := map[string]interface{}{}
	for k, v := range dataIn {
		reData[k.(string)] = v
	}

	entityType := OpenInstallationQuickstartEntityType{
		Type:   toStringByFieldName("type", reData),
		Domain: toStringByFieldName("domain", reData),
	}

	return entityType
}

func expandSuccessLinkConfig(recipe map[string]interface{}) OpenInstallationSuccessLinkConfig {
	v, ok := recipe["successLinkConfig"]
	if !ok {
		return OpenInstallationSuccessLinkConfig{}
	}

	dataIn := v.(map[interface{}]interface{})
	reData := map[string]interface{}{}
	for k, v := range dataIn {
		reData[k.(string)] = v
	}

	dataOut := OpenInstallationSuccessLinkConfig{
		Filter: toStringByFieldName("filter", reData),
	}

	if v, ok := reData["type"]; ok {
		dataOut.Type = OpenInstallationSuccessLinkType(v.(string))
	}

	return dataOut
}

func expandInstallTargets(recipe map[string]interface{}) []OpenInstallationRecipeInstallTarget {
	v, ok := recipe["installTargets"]
	if !ok {
		return []OpenInstallationRecipeInstallTarget{}
	}

	dataIn := v.([]interface{})
	dataOut := make([]OpenInstallationRecipeInstallTarget, len(dataIn))
	dataz := make([]map[string]interface{}, len(dataIn))
	for i, vv := range dataIn {
		vvv := vv.(map[interface{}]interface{})
		varr := map[string]interface{}{}

		for k, v := range vvv {
			varr[k.(string)] = v
			dataz[i] = varr
		}
	}

	for i, v := range dataz {
		vOut := OpenInstallationRecipeInstallTarget{
			KernelArch:      toStringByFieldName("kernelArch", v),
			KernelVersion:   toStringByFieldName("kernelVersion", v),
			PlatformVersion: toStringByFieldName("platformVersion", v),
		}

		if v, ok := v["os"]; ok {
			vOut.Os = OpenInstallationOperatingSystem(v.(string))
		}

		if v, ok := v["platform"]; ok {
			vOut.Platform = OpenInstallationPlatform(v.(string))
		}

		if v, ok := v["platformFamily"]; ok {
			vOut.PlatformFamily = OpenInstallationPlatformFamily(v.(string))
		}

		if v, ok := v["type"]; ok {
			vOut.Type = OpenInstallationTargetType(v.(string))
		}

		dataOut[i] = vOut
	}

	return dataOut
}

func expandPreInstall(recipe map[string]interface{}) OpenInstallationPreInstallConfiguration {
	v, ok := recipe["preInstall"]
	if !ok {
		return OpenInstallationPreInstallConfiguration{}
	}

	vv := v.(map[interface{}]interface{})
	infoOut := map[string]interface{}{}
	for k, v := range vv {
		infoOut[k.(string)] = v
	}

	return OpenInstallationPreInstallConfiguration{
		Info:               toStringByFieldName("info", infoOut),
		Prompt:             toStringByFieldName("prompt", infoOut),
		RequireAtDiscovery: toStringByFieldName("requireAtDiscovery", infoOut),
	}
}

func expandPostInstall(recipe map[string]interface{}) OpenInstallationPostInstallConfiguration {
	v, ok := recipe["postInstall"]
	if !ok {
		return OpenInstallationPostInstallConfiguration{}
	}

	vv := v.(map[interface{}]interface{})
	infoOut := map[string]interface{}{}
	for k, v := range vv {
		infoOut[k.(string)] = v
	}

	return OpenInstallationPostInstallConfiguration{
		Info: toStringByFieldName("info", infoOut),
	}
}

func expandInputVars(recipe map[string]interface{}) []OpenInstallationRecipeInputVariable {
	v, ok := recipe["inputVars"]
	if !ok {
		return []OpenInstallationRecipeInputVariable{}
	}

	vars := v.([]interface{})
	varsOut := make([]OpenInstallationRecipeInputVariable, len(vars))

	varz := make([]map[string]interface{}, len(vars))
	for i, vv := range vars {
		vvv := vv.(map[interface{}]interface{})
		varr := map[string]interface{}{}

		for k, v := range vvv {
			varr[k.(string)] = v
			varz[i] = varr
		}
	}

	for i, v := range varz {
		vOut := OpenInstallationRecipeInputVariable{
			Default: toStringByFieldName("default", v),
			Name:    toStringByFieldName("name", v),
			Prompt:  toStringByFieldName("prompt", v),
			Secret:  toBoolByFieldName("secret", v),
		}

		varsOut[i] = vOut
	}

	return varsOut
}

func expandLogMatch(recipe map[string]interface{}) []OpenInstallationLogMatch {
	v, ok := recipe["logMatch"]
	if !ok {
		return []OpenInstallationLogMatch{}
	}

	dataIn := v.([]interface{})
	dataOut := make([]OpenInstallationLogMatch, len(dataIn))
	dataz := make([]map[string]interface{}, len(dataIn))
	for i, vv := range dataIn {
		vvv := vv.(map[interface{}]interface{})
		varr := map[string]interface{}{}

		for k, v := range vvv {
			varr[k.(string)] = v
			dataz[i] = varr
		}
	}

	for i, v := range dataz {
		vOut := OpenInstallationLogMatch{
			Attributes: expandLogAttributes(v),
			File:       toStringByFieldName("file", v),
			Name:       toStringByFieldName("name", v),
			Pattern:    toStringByFieldName("pattern", v),
			Systemd:    toStringByFieldName("systemd", v),
		}

		dataOut[i] = vOut
	}

	return dataOut
}

func expandLogAttributes(data map[string]interface{}) OpenInstallationAttributes {
	attributesOut := OpenInstallationAttributes{}

	v, ok := data["attributes"]
	if !ok {
		return attributesOut
	}

	attrs := v.(map[interface{}]interface{})
	attrsOut := map[string]string{}
	for k, v := range attrs {
		attrsOut[k.(string)] = v.(string)
	}

	if v, ok := attrsOut["logtype"]; ok {
		attributesOut.Logtype = v
	}

	return attributesOut
}

func toBoolByFieldName(fieldName string, data map[string]interface{}) bool {
	if v, ok := data[fieldName]; ok {
		return v.(bool)
	}

	return false
}

func toStringByFieldName(fieldName string, data map[string]interface{}) string {
	out := ""
	if in, ok := data[fieldName]; ok {
		switch v := in.(type) {
		case int:
			return strconv.Itoa(v)
		case string:
			return v
		case bool:
			return strconv.FormatBool(v)
		}
		return out
	}

	return out
}

func expandInstalllMapToString(recipeIn map[string]interface{}) (string, error) {
	installIn, ok := recipeIn["install"]
	if !ok {
		return "", nil
	}

	installOut := map[string]interface{}{}
	installMap := installIn.(map[interface{}]interface{})
	for k, v := range installMap {
		installOut[k.(string)] = v
	}

	installAsString, err := yaml.Marshal(installOut)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling recipe.install to string: %s", err)
	}

	return string(installAsString), nil
}

func interfaceSliceToStringSlice(slice []interface{}) []string {
	out := make([]string, len(slice))

	for i, v := range slice {
		out[i] = v.(string)
	}

	return out
}

func (r *OpenInstallationRecipe) PostInstallMessage() string {
	if r.PostInstall.Info != "" {
		return r.PostInstall.Info
	}

	return ""
}

func (r *OpenInstallationRecipe) PreInstallMessage() string {
	if r.PreInstall.Info != "" {
		return r.PreInstall.Info
	}

	return ""
}

// SetRecipeVar is responsible for including a new variable on the RecipeVariables
// struct, which is used by go-task executor.
func (r *OpenInstallationRecipe) SetRecipeVar(key string, value string) {
	RecipeVariables[key] = value
}

func (r *OpenInstallationRecipe) IsApm() bool {
	return r.HasKeyword("apm")
}

func (r *OpenInstallationRecipe) HasHostTargetType() bool {
	return r.HasTargetType(OpenInstallationTargetTypeTypes.HOST)
}

func (r *OpenInstallationRecipe) HasApplicationTargetType() bool {
	return r.HasTargetType(OpenInstallationTargetTypeTypes.APPLICATION)
}

func (r *OpenInstallationRecipe) HasKeyword(keyword string) bool {
	if len(r.Keywords) == 0 {
		return false
	}

	for _, single := range r.Keywords {
		if strings.EqualFold(single, keyword) {
			return true
		}
	}

	return false
}

func (r *OpenInstallationRecipe) HasTargetType(t OpenInstallationTargetType) bool {
	if len(r.InstallTargets) == 0 {
		return false
	}

	for _, target := range r.InstallTargets {
		if target.Type == t {
			return true
		}
	}

	return false
}

func (r *OpenInstallationRecipe) GetOrderKey() string {
	if r.Name == InfraAgentRecipeName {
		return fmt.Sprintf("%d-%s", 10, InfraAgentRecipeName)
	}
	if r.Name == LoggingRecipeName {
		return fmt.Sprintf("%d-%s", 20, LoggingRecipeName)
	}
	return fmt.Sprintf("%d-%s", 50, r.Name)
}
