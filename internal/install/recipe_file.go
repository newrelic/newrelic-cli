package install

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type recipeFile struct {
	Description    string                 `yaml:"description"`
	InputVars      []variableConfig       `yaml:"inputVars"`
	Install        map[string]interface{} `yaml:"install"`
	InstallTargets []recipeInstallTarget  `yaml:"installTargets"`
	Keywords       []string               `yaml:"keywords"`
	LogMatch       []logMatch             `yaml:"logMatch"`
	Name           string                 `yaml:"name"`
	ProcessMatch   []string               `yaml:"processMatch"`
	Repository     string                 `yaml:"repository"`
	ValidationNRQL string                 `yaml:"validationNrql"`
}

type variableConfig struct {
	Name    string `yaml:"name"`
	Prompt  string `yaml:"prompt"`
	Secret  bool   `secret:"prompt"`
	Default string `yaml:"default"`
}

type recipeInstallTarget struct {
	Type            string `yaml:"type"`
	OS              string `yaml:"os"`
	Platform        string `yaml:"platform"`
	PlatformFamily  string `yaml:"platformFamily"`
	PlatformVersion string `yaml:"platformVersion"`
	KernelVersion   string `yaml:"kernelVersion"`
	KernelArch      string `yaml:"kernelArch"`
}

type logMatch struct {
	Name       string             `yaml:"name"`
	File       string             `yaml:"file"`
	Attributes logMatchAttributes `yaml:"attributes"`
	Pattern    string             `yaml:"pattern"`
	Systemd    string             `yaml:"systemd"`
}

type logMatchAttributes struct {
	LogType string `yaml:"logtype"`
}

func loadRecipeFile(filename string) (*recipeFile, error) {
	out, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	f, err := newRecipeFile(string(out))
	if err != nil {
		return nil, err
	}

	return f, nil
}

func newRecipeFile(recipeFileString string) (*recipeFile, error) {
	var f recipeFile
	err := yaml.Unmarshal([]byte(recipeFileString), &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (f *recipeFile) String() (string, error) {
	out, err := yaml.Marshal(f)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (f *recipeFile) ToRecipe() (*recipe, error) {
	fileStr, err := f.String()
	if err != nil {
		return nil, err
	}

	r := recipe{
		File:           fileStr,
		Name:           f.Name,
		Description:    f.Description,
		Repository:     f.Repository,
		Keywords:       f.Keywords,
		ProcessMatch:   f.ProcessMatch,
		LogMatch:       f.LogMatch,
		ValidationNRQL: f.ValidationNRQL,
	}

	return &r, nil
}
