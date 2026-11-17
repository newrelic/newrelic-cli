# Inventory Published and Ownership

All assets published from this repository and from `newrelic/open-install-library` are
owned by Team Virtuoso (`@newrelic/virtuoso`), which is accountable for every published
artifact, its integrity, and its publish path.

---

## New Relic CLI

| Asset | Destinations |
|-------|-------------|
| `newrelic-linux-amd64.tar.gz` and other Linux tarballs | S3 (`nr-downloads-main/install/newrelic-cli/`), GitHub Releases |
| `newrelic-darwin-amd64.tar.gz` and other macOS tarballs | S3, GitHub Releases, Homebrew (`newrelic/newrelic-cli`) |
| `NewRelicCLIInstaller.msi` (Windows installer) | S3, GitHub Releases, Chocolatey (`newrelic-cli`) |
| `newrelic_windows_amd64.zip` and other Windows zips | S3, GitHub Releases |
| `newrelic-cli_linux_amd64.deb` and other DEB packages | S3, GitHub Releases |
| `newrelic-cli_linux_amd64.rpm` and other RPM packages | S3, GitHub Releases |
| Docker image `newrelic/newrelic-cli` | Docker Hub |
| `install.sh` (Linux/macOS installer script) | S3 (`nr-downloads-main/install/newrelic-cli/install.sh`) |
| `install.ps1` (Windows installer script) | S3 (`nr-downloads-main/install/newrelic-cli/install.ps1`) |
| `checksums.txt` | S3, GitHub Releases |
| `checksums.txt.sig` (GPG signature, key `4F9A9B5B96EC30B9`) | S3, GitHub Releases |
| `NewRelicCLIInstaller.msi.sha256` | S3, GitHub Releases |

---

## Open Install Library

| Asset | Destinations |
|-------|-------------|
| `recipes.zip` | S3 (`nr-downloads-main/install/open-install-library/`), GitHub Releases |
| `recipes.zip.sha256` | S3, GitHub Releases |

---

## Access Control

Write access to all S3 paths and GitHub Release assets is restricted to the release
workflow via GitHub OIDC (short-lived role assumption — no long-lived AWS credentials).
The release workflow is gated to the `production` GitHub Environment, which requires
approval from a named reviewer before any release job proceeds.

`CODEOWNERS` (`.github/CODEOWNERS`) ensures that any PR touching release-critical paths
requires review from `@newrelic/virtuoso`.
