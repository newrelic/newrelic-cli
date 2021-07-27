package recipes

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

const (
	embeddedRecipesPath = "files"
)

var (
	//go:embed files/*
	EmbeddedFS embed.FS
)

type EmbeddedRecipeFetcher struct{}

func NewEmbeddedRecipeFetcher() *EmbeddedRecipeFetcher {
	return &EmbeddedRecipeFetcher{}
}

func (f *EmbeddedRecipeFetcher) FetchRecipes(context.Context) (out []types.OpenInstallationRecipe, err error) {
	files, err := f.getYAMLFiles(embeddedRecipesPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		b, err := EmbeddedFS.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("could not read embedded file %s: %w", f, err)
		}

		var r types.OpenInstallationRecipe
		if err := yaml.Unmarshal(b, &r); err != nil {
			return nil, fmt.Errorf("could not unmarshal embedded file %s: %w", f, err)
		}

		out = append(out, r)
	}

	return out, nil
}

func (f *EmbeddedRecipeFetcher) FetchLibraryVersion(ctx context.Context) string {
	versionFilename := "version.txt"
	data, err := EmbeddedFS.ReadFile(embeddedRecipesPath + "/" + versionFilename)
	if err == nil {
		return string(data)
	}
	return ""
}

func (f *EmbeddedRecipeFetcher) getYAMLFiles(path string) (out []string, err error) {
	return f.getFiles(path, isYAMLFile)
}

func isYAMLFile(path string) bool {
	e := filepath.Ext(path)
	return strings.EqualFold(e, ".yml") || strings.EqualFold(e, ".yaml")
}

func (f *EmbeddedRecipeFetcher) getFiles(path string, filterFunc func(string) bool) (out []string, err error) {
	dirs, err := EmbeddedFS.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, d := range dirs {
		pathname := path + "/" + d.Name()
		if d.IsDir() {
			files, err := f.getFiles(pathname, filterFunc)
			if err != nil {
				return nil, err
			}

			out = append(out, files...)
		} else {
			if filterFunc(d.Name()) {
				out = append(out, pathname)
			}
		}
	}

	return out, nil
}
