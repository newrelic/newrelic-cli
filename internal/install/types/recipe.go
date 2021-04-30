package types

import "strings"

const (
	InfraAgentRecipeName = "infrastructure-agent-installer"
	LoggingRecipeName    = "logs-integration"
)

type Recipe struct {
	ID                string                                   `json:"id" yaml:"id"`
	Description       string                                   `json:"description" yaml:"description"`
	DisplayName       string                                   `json:"displayName" yaml:"displayName"`
	Dependencies      []string                                 `json:"dependencies" yaml:"dependencies"`
	Stability         OpenInstallationStability                `json:"stability" yaml:"stability"`
	Quickstarts       OpenInstallationQuickstartsFilter        `json:"quickstarts,omitempty" yaml:"quickstarts"`
	File              string                                   `json:"file" yaml:"file"`
	InstallTargets    []OpenInstallationRecipeInstallTarget    `json:"installTargets" yaml:"installTargets"`
	Keywords          []string                                 `json:"keywords" yaml:"keywords"`
	LogMatch          []LogMatch                               `json:"logMatch" yaml:"logMatch"`
	Name              string                                   `json:"name" yaml:"name"`
	PreInstall        OpenInstallationPreInstallConfiguration  `json:"preInstall" yaml:"preInstall"`
	PostInstall       OpenInstallationPostInstallConfiguration `json:"postInstall" yaml:"postInstall"`
	ProcessMatch      []string                                 `json:"processMatch" yaml:"processMatch"`
	Repository        string                                   `json:"repository" yaml:"repository"`
	SuccessLinkConfig OpenInstallationSuccessLinkConfig        `json:"successLinkConfig" yaml:"successLinkConfig"`
	ValidationNRQL    string                                   `json:"validationNrql" yaml:"validationNrql"`
	Vars              map[string]interface{}
}

func (r *Recipe) PostInstallMessage() string {
	if r.PostInstall.Info != "" {
		return r.PostInstall.Info
	}

	return ""
}

func (r *Recipe) PreInstallMessage() string {
	if r.PreInstall.Info != "" {
		return r.PreInstall.Info
	}

	return ""
}

// LogMatch represents a pattern that may match one or more logs on the underlying host.
type LogMatch struct {
	Name       string             `yaml:"name"`
	File       string             `yaml:"file"`
	Attributes LogMatchAttributes `yaml:"attributes,omitempty"`
	Pattern    string             `yaml:"pattern,omitempty"`
	Systemd    string             `yaml:"systemd,omitempty"`
}

// LogMatchAttributes contains metadata about its parent LogMatch.
type LogMatchAttributes struct {
	LogType string `yaml:"logtype"`
}

type RecipeVars map[string]string

// AddVar is responsible for including a new variable on the recipe Vars
// struct, which is used by go-task executor.
func (r *Recipe) AddVar(key string, value interface{}) {
	if len(r.Vars) == 0 {
		r.Vars = make(map[string]interface{})
	}

	r.Vars[key] = value
}

func (r *Recipe) IsApm() bool {
	return r.HasKeyword("apm")
}

func (r *Recipe) HasHostTargetType() bool {
	return r.HasTargetType(OpenInstallationTargetTypeTypes.HOST)
}

func (r *Recipe) HasApplicationTargetType() bool {
	return r.HasTargetType(OpenInstallationTargetTypeTypes.APPLICATION)
}

func (r *Recipe) HasKeyword(keyword string) bool {
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

func (r *Recipe) HasTargetType(t OpenInstallationTargetType) bool {
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
