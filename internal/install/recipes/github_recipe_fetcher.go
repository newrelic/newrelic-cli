package recipes

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/google/go-github/v35/github"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type GithubRecipeFetcher struct{}

const githubOrg = "newrelic"
const githubRepo = "open-install-library"

func (f *GithubRecipeFetcher) FetchRecipe(ctx context.Context, manifest *types.DiscoveryManifest, friendlyName string) (*types.Recipe, error) {
	recipes, err := f.FetchRecommendations(ctx, manifest)
	if err != nil {
		return nil, err
	}

	for _, recipe := range recipes {
		if recipe.Name == friendlyName {
			return &recipe, nil
		}
	}

	return nil, fmt.Errorf("%s: %w", friendlyName, ErrRecipeNotFound)
}

func (f *GithubRecipeFetcher) FetchRecommendations(ctx context.Context, manifest *types.DiscoveryManifest) ([]types.Recipe, error) {
	recipes, err := f.FetchRecipes(ctx, manifest)
	if err != nil {
		return nil, err
	}

	return manifest.ConstrainRecipes(recipes), nil
}

func (f *GithubRecipeFetcher) FetchRecipes(ctx context.Context, manifest *types.DiscoveryManifest) ([]types.Recipe, error) {
	var recipes []types.Recipe
	var err error

	err = cacheRecipes(ctx)
	if err != nil {
		return nil, err
	}

	recipes, err = loadRecipesFromCache(ctx)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func loadRecipesFromCache(ctx context.Context) ([]types.Recipe, error) {
	cacheDir := filepath.Join(config.DefaultConfigDirectory, "recipes")

	return loadRecipesFromDir(ctx, cacheDir)
}

func cacheRecipes(ctx context.Context) error {
	return cacheLatestRelease(ctx)
}

func cacheLatestRelease(ctx context.Context) error {
	release, url, err := getLatestRelease(ctx)
	if err != nil {
		return err
	}

	cacheDir := filepath.Join(config.DefaultConfigDirectory, "cache")
	filepath.Dir(cacheDir)

	path := filepath.Join(
		cacheDir,
		fmt.Sprintf("recipes-%s.zip", release),
	)

	return cacheRecipeArchive(ctx, url, path)
}

func cacheRecipeArchive(ctx context.Context, u *url.URL, path string) error {
	err := fetchRecipeArchive(ctx, u, path)
	if err != nil {
		return err
	}

	recipeDir := filepath.Join(config.DefaultConfigDirectory, "recipes")

	err = unzipRecipes(ctx, path, recipeDir)
	if err != nil {
		return err
	}

	return nil
}

func unzipRecipes(ctx context.Context, archivePath string, destDir string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	err = os.MkdirAll(destDir, 0750)
	if err != nil {
		return err
	}

	extractAndWriteFileToPath := func(f *zip.File, path string) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if fileErr := rc.Close(); err != nil {
				log.Error(fileErr)
			}
		}()

		log.WithFields(log.Fields{
			"archive": archivePath,
			"dest":    destDir,
			"name":    f.Name,
		}).Trace("extracting zip file")

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				return err
			}
		} else {
			err = os.MkdirAll(filepath.Dir(path), 0750)
			if err != nil {
				return err
			}

			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if fileErr := f.Close(); err != nil {
					log.Error(fileErr)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}

		return nil
	}

	re := regexp.MustCompile(`newrelic-open-install-library-[a-f0-9]{7}/recipes/(.*\.ya?ml)`)
	for _, f := range r.File {
		matches := re.FindAllStringSubmatch(f.Name, -1)
		if matches == nil {
			continue
		}

		path := filepath.Join(
			config.DefaultConfigDirectory,
			"recipes",
			matches[0][1],
		)

		err := extractAndWriteFileToPath(f, path)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	return nil
}

func fetchRecipeArchive(ctx context.Context, u *url.URL, path string) error {
	log.WithFields(log.Fields{
		"url":  u.String(),
		"path": path,
	}).Debug("caching recipe archive")

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err = os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		err = os.Mkdir(filepath.Dir(path), 0700)
		if err != nil {
			return err
		}
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func getLatestRelease(ctx context.Context) (string, *url.URL, error) {
	var err error
	var latestRelease string
	var latestReleaseURL *url.URL

	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases(ctx, githubOrg, githubRepo, nil)
	if err != nil {
		log.Error(err)
	}

	if len(releases) > 0 {
		release := releases[0]

		latestReleaseURL, err = url.Parse(release.GetZipballURL())
		if err != nil {
			return "", nil, err
		}

		if release.TagName != nil {
			latestRelease = *release.TagName
		}
	}

	if latestRelease != "" && latestReleaseURL != nil {
		return latestRelease, latestReleaseURL, nil
	}

	return "", nil, fmt.Errorf("unable to determine latest libray release")
}
