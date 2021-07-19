package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

// GitHubRepositoryTagResponse is the data structure returned
// from a repository's `/releases/latest` endpoint.
type GitHubRepositoryTagResponse struct {
	TagName string `json:"tag_name"`
}

// NewRelicCLILatestReleaseURL is the URL used to fetch the latest release data.
const NewRelicCLILatestReleaseURL string = "https://api.github.com/repos/newrelic/newrelic-cli/releases/latest"

// NewRelicCLILatestReleaseURL is the URL used to fetch the latest release data.
const PrereleaseEnvironmentMsgFormat string = `
  It appears you're in a development environment using prerelease version %s.
  To upgrade the New Relic CLI, you must be using non-prerelease version.
`

const OldVersionMsgFormat string = `
  You are using an old version of the New Relic CLI.

    Current version: %s
    Latest version:  %s

  Upgrade to the latest version using the following command:

    newrelic version upgrade
`

// IsLatestVersion returns true if the latest remote release version matches
// the currently installed version.
func IsLatestVersion(ctx context.Context, currentVersion string) (bool, error) {
	cv, err := semver.NewVersion(currentVersion)
	if err != nil {
		log.Fatalf("error parsing current CLI version %s: %s", cv.String(), err.Error())
	}

	latestRelease, err := FetchLatestRelease(ctx)
	if err != nil {
		return false, fmt.Errorf("error fetching latest release %s: %s", latestRelease.TagName, err.Error())
	}

	latestVersion, err := semver.NewVersion(latestRelease.TagName)
	if err != nil {
		return false, fmt.Errorf("error parsing latest tag %s: %s", latestRelease.TagName, err.Error())
	}

	if cv.LessThan(latestVersion) {
		return false, nil
	}

	return cv.Equal(latestVersion), nil
}

func FetchLatestRelease(ctx context.Context) (*GitHubRepositoryTagResponse, error) {
	client := utils.NewHTTPClient()

	respBytes, err := client.Get(ctx, NewRelicCLILatestReleaseURL)
	if err != nil {
		return nil, err
	}

	repoTag := GitHubRepositoryTagResponse{}
	err = json.Unmarshal(respBytes, &repoTag)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"tag": repoTag.TagName,
	}).Debug("fetch tag success")

	return &repoTag, nil
}

// UpgradeCLIVersion handles upgrading the CLI to newer version.
// By default, the latest version is installed. If a target version
// if provided, the specified version will be installed.
func UpgradeCLIVersion(ctx context.Context, internalVersion string, targetVersion string, forceUpgrade bool) error {
	currentVersion, err := semver.NewVersion(internalVersion)
	if err != nil {
		return fmt.Errorf("error parsing current CLI version %s", internalVersion)
	}

	isDevEnvironment := currentVersion.Prerelease() != ""
	if isDevEnvironment && !forceUpgrade {
		return fmt.Errorf(PrereleaseEnvironmentMsgFormat, internalVersion)
	}

	release, err := FetchLatestRelease(ctx)
	if err != nil {
		return err
	}

	latestVersion, err := semver.NewVersion(release.TagName)
	if err != nil {
		return fmt.Errorf("error parsing latest CLI version %s", release.TagName)
	}

	if currentVersion.LessThan(latestVersion) {
		output.Printf(OldVersionMsgFormat, currentVersion.String(), latestVersion.String())
	}

	// TODO: Refactor this function now that the POC seems viable.

	return nil
}
