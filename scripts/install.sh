#!/usr/bin/env bash

# Install the New Relic CLI.
# https://github.com/newrelic/newrelic-cli
#
# Dependencies: curl, cut, tar, gzip
#
# The version to install and the binary location can be passed in via NEW_RELIC_CLI_VERSION and DESTDIR respectively.
#

set -o errexit

echo "Starting installation."

# Determine release filename. This can be expanded with CPU arch in the future.
if [ "$(uname)" == "Linux" ]; then
	OS="Linux"
elif [ "$(uname)" == "Darwin" ]; then
	OS="Darwin"
else
	echo "This operating system is not supported."
	exit 1
fi

for x in curl cut tar gzip; do
    which $x > /dev/null || (echo "Unable to continue.  Please install $x before proceed."; exit 1)
done

# GitHub's URL for the latest release, will redirect.
LATEST_URL="https://github.com/newrelic/newrelic-cli/releases/latest"
DESTDIR="${DESTDIR:-/usr/local/bin}"

if [ -z "$NEW_RELIC_CLI_VERSION" ]; then
	echo "CHECK PASSED v${NEW_RELIC_CLI_VERSION}"

	NEW_RELIC_CLI_VERSION=$(curl -LI -o /dev/null -w '%{url_effective}' $LATEST_URL | cut -d "v" -f 2)
fi

echo "Installing New Relic CLI v${NEW_RELIC_CLI_VERSION} - $1"

# Run the script in a temporary directory that we know is empty.
SCRATCH=$(mktemp -d || mktemp -d -t 'tmp')
cd "$SCRATCH"

function error {
  echo "An error occurred installing the tool."
  echo "The contents of the directory $SCRATCH have been left in place to help to debug the issue."
}

trap error ERR

RELEASE_URL="https://github.com/newrelic/newrelic-cli/releases/download/v${NEW_RELIC_CLI_VERSION}/newrelic-cli_${NEW_RELIC_CLI_VERSION}_${OS}_x86_64.tar.gz"

# Download & unpack the release tarball.
curl -sL --retry 3 "${RELEASE_URL}" | tar -xz

echo "Installing to $DESTDIR"

mv newrelic "$DESTDIR"
chmod +x "$DESTDIR/newrelic"

# Delete the working directory when the install was successful.
rm -r "$SCRATCH"
