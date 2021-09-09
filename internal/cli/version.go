package cli

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Masterminds/semver/v3"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

var (
	// Represents the currently installed/compiled version of the CLI
	version string

	// Represents the latest release version per the GitHub repository's latest release
	latestVersion string
)

const installCLISnippetLinux = `curl -Ls https://download.newrelic.com/install/newrelic-cli/scripts/install.sh | bash && sudo newrelic install`
const installCLISnippetWindows = `[Net.ServicePointManager]::SecurityProtocol = 'tls12, tls'; (New-Object System.Net.WebClient).DownloadFile(\"https://download.newrelic.com/install/newrelic-cli/scripts/install.ps1\", \"$env:TEMP\\install.ps1\");", "& $env:TEMP\\install.ps1;", "& 'C:\\Program Files\\New Relic\\New Relic CLI\\newrelic.exe' install`

// NewRelicCLILatestReleaseURL is the URL used to fetch the latest release data utilizing GitHub's API.
const NewRelicCLILatestReleaseURL string = "https://download.newrelic.com/install/newrelic-cli/currentVersion.txt"

// UpdateVersionMsgFormat is the message displayed to a user when an older
// version of the CLI is installed.
const UpdateVersionMsgFormat string = `
  We need to update your New Relic CLI version to continue.

    Installed version: %s
    Latest version:    %s

  To update your CLI and continue this installation, run this command:

    %s
`

// Version returns the version of the CLI that's currently installed.
func Version() string {
	if version == "" {
		version = os.Getenv("NEW_RELIC_CLI_VERSION")
	}

	return version
}

// IsLatestVersion returns true if the provided version string matches
// the current installed version.
func IsLatestVersion(ctx context.Context, latestVersion string) (bool, error) {
	installedVersion, err := semver.NewVersion(Version())
	if err != nil {
		return false, fmt.Errorf("error parsing installed CLI version %s: %s", version, err.Error())
	}

	lv, err := semver.NewVersion(latestVersion)
	if err != nil {
		return false, fmt.Errorf("error parsing version to check %s: %s", latestVersion, err.Error())
	}

	if installedVersion.LessThan(lv) {
		return false, nil
	}

	return installedVersion.Equal(lv), nil
}

// GetLatestReleaseVersion returns the latest released tag.
// The latest tag is pulled from the newrelic-cli GitHub repository.
func GetLatestReleaseVersion(ctx context.Context) (string, error) {
	if latestVersion != "" {
		return latestVersion, nil
	}

	latestRelease, err := fetchLatestRelease(ctx)
	if err != nil {
		return "", fmt.Errorf("error fetching latest release: %w", err)
	}

	lv, err := semver.NewVersion(latestRelease)
	if err != nil {
		return "", fmt.Errorf("error parsing latest release version string '%s': %w", latestRelease, err)
	}

	// Cache the latest version string
	versionString := lv.String()
	setLatestVersion(versionString)

	return versionString, nil
}

// IsDevEnvironment is a naive implementation to determine if the CLI
// is being run in a dev environment. IsDevEnvironment returns true when
// the installed CLI version is either in a prerelease state or in a dirty state.
// The version string is generated at compile time using git. The prerelease string
// is appended to the primary semver version string.
//
// If you're doing local development on the CLI, your version may look similar to
// the examples below.
//
// Examples of versions that have a prerelease tag (i.e. the suffix):
//
//  v0.32.1-10-gbe63a24
//  v0.32.1-10-gbe63a24-dirty
//
// In this example version string, "10" represents the number of commits
// since the 0.32.1 tag was created. The "gbe63a24" is the previous commit's
// abbreviated sha. The "dirty" part means that git was in a dirty state at compile
// time, meaning an updated file was saved, but not yet committed.
func IsDevEnvironment() bool {
	v, err := semver.NewVersion(version)
	if err != nil {
		return true
	}

	prereleaseString := v.Prerelease()
	if prereleaseString == "" {
		return false
	}

	if strings.Contains(prereleaseString, "dirty") {
		return true
	}

	hasDevVersionSuffix := len(strings.Split(prereleaseString, "-")) > 1

	return hasDevVersionSuffix
}

func PrintUpdateCLIMessage(latestReleaseVersion string) {
	output.Printf("%s", FormatUpdateVersionMessage(latestReleaseVersion))
}

func FormatUpdateVersionMessage(latestReleaseVersion string) string {
	snippet := installCLISnippetLinux

	if isWindowsOS() {
		snippet = installCLISnippetWindows
	}

	return fmt.Sprintf(UpdateVersionMsgFormat, version, latestReleaseVersion, snippet)
}

func setLatestVersion(v string) {
	latestVersion = v
}

func isWindowsOS() bool {
	return runtime.GOOS == "windows"
}

func fetchLatestRelease(ctx context.Context) (string, error) {
	client := utils.NewHTTPClient()

	respBytes, err := client.Get(ctx, NewRelicCLILatestReleaseURL)
	if err != nil {
		return "", err
	}

	gitTag := string(respBytes)
	_, err = semver.NewVersion(gitTag)
	if err != nil {
		return "", err
	}

	log.WithFields(log.Fields{
		"tag": gitTag,
	}).Debug("fetch tag success")

	return gitTag, nil
}
