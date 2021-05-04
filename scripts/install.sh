#!/usr/bin/env bash

# Install the New Relic CLI.
# https://github.com/newrelic/newrelic-cli
#
# Dependencies: curl, cut, tar, gzip
#
# The version to install and the binary location can be passed in via VERSION and DESTDIR respectively.
#

set -o errexit

echo "Starting installation."

# Determine release filename. This can be expanded with CPU arch in the future.
if [ "$(uname)" == "Linux" ]; then
    OS="Linux"
elif [ "$(uname)" == "Darwin" ]; then
    OS="Darwin"
else
    echo "This operating system is not supported. The supported operating systems are Linux and Darwin"
    exit 1
fi

if [ "$(uname -m)" == "x86_64" ]; then
    MACHINE="x86_64"
elif [ "$(uname)" == "aarch64" ]; then
    MACHINE="aarch64"
else
    echo "This machine architecture is not supported. The supported architectures are x86_64 and aarch64."
    exit 1
fi

for x in curl cut tar gzip sudo; do
    which $x > /dev/null || (echo "Unable to continue.  Please install $x before proceeding."; exit 1)
done

# GitHub's URL for the latest release, will redirect.
LATEST_URL="https://github.com/newrelic/newrelic-cli/releases/latest"
DESTDIR="${DESTDIR:-/usr/local/bin}"

if [ -z "$VERSION" ]; then
    VERSION=$(curl -sLI -o /dev/null -w '%{url_effective}' $LATEST_URL | cut -d "v" -f 2)
fi

echo "Installing New Relic CLI v${VERSION}"

# Run the script in a temporary directory that we know is empty.
SCRATCH=$(mktemp -d || mktemp -d -t 'tmp')
cd "$SCRATCH"

function error {
  echo "An error occurred installing the tool."
  echo "The contents of the directory $SCRATCH have been left in place to help to debug the issue."
}

trap error ERR

RELEASE_URL="https://github.com/newrelic/newrelic-cli/releases/download/v${VERSION}/newrelic-cli_${VERSION}_${OS}_${MACHINE}.tar.gz"

# Download & unpack the release tarball.
curl -sL --retry 3 "${RELEASE_URL}" | tar -xz

if [ "$UID" != "0" ]; then
    echo "Installing to $DESTDIR using sudo"
    sudo mv newrelic "$DESTDIR"
    sudo chmod +x "$DESTDIR/newrelic"
    sudo chown root:0 "$DESTDIR/newrelic"
else
    echo "Installing to $DESTDIR"
    mv newrelic "$DESTDIR"
    chmod +x "$DESTDIR/newrelic"
    chown root:0 "$DESTDIR/newrelic"
fi

# Delete the working directory when the install was successful.
rm -r "$SCRATCH"
