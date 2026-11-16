#!/usr/bin/env bash

# Install the New Relic CLI.
# https://github.com/newrelic/newrelic-cli
#
# Dependencies: curl, cut, tar, gzip
#
# The version to install and the binary location can be passed in via VERSION and DESTDIR respectively.
#

set -o errexit

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
elif [ "$(uname -m)" == "aarch64" ] || [ "$(uname -m)" == "arm64" ]; then
    MACHINE="arm64"
elif [ "$(uname -m)" == "armv7l" ]; then
    MACHINE="armv7"
else
    echo "This machine architecture is not supported. The supported architectures are x86_64, aarch64, armv7."
    exit 1
fi

for x in cut tar gzip sudo; do
    which $x > /dev/null || (echo "Unable to continue.  Please install $x before proceeding."; exit 1)
done

DISTRO=$(cat /etc/issue /etc/system-release /etc/redhat-release /etc/os-release 2>/dev/null | grep -m 1 -Eo "(Ubuntu|Amazon|CentOS|Debian|Red Hat|SUSE)" || true)

IS_CURL_INSTALLED=$(which curl | wc -l)
if [ $IS_CURL_INSTALLED -eq 0 ]; then
    echo "curl is required to install, please confirm Y/N to install (default Y): "
    read -r CONFIRM_CURL
    if [ "$CONFIRM_CURL" == "Y" ] || [ "$CONFIRM_CURL" == "y" ] || [ "$CONFIRM_CURL" == "" ]; then
        if [ "$DISTRO" == "Ubuntu" ] || [ "$DISTRO" == "Debian" ]; then
            sudo apt-get update
            sudo apt-get install curl -y
        elif [ "$DISTRO" == "Amazon" ] || [ "$DISTRO" == "CentOS" ] || [ "$DISTRO" == "Red Hat" ]; then
            sudo yum install curl -y
        elif [ "$DISTRO" == "SUSE" ]; then
            sudo zypper -n install curl
        else
            echo "Unable to continue. Please install curl manually before proceeding."; exit 131
        fi
    else
        echo "Unable to continue without curl. Please install curl before proceeding."; exit 131
    fi
fi

# GitHub's URL for the latest release, will redirect.
BASE_URL="https://download.newrelic.com"
LATEST_URL="$BASE_URL/install/newrelic-cli/currentVersion.txt"
DESTDIR="${DESTDIR:-/usr/local/bin}"

# Check for connectivity to https://download.newrelic.com
curl --connect-timeout 10 -IsL "$BASE_URL" > /dev/null || ( echo "Cannot connect to $BASE_URL to download the New Relic CLI. Check your firewall settings. If you are using a proxy, make sure you have set the HTTPS_PROXY environment variable." && exit 132 )

# Create DESTDIR if it does not exist.
if [ ! -d "$DESTDIR" ]; then
    mkdir -m 755 -p "$DESTDIR"
fi

if [ -z "$VERSION" ]; then
    VERSION=$(curl -sL $LATEST_URL | cut -d "v" -f 2)
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

TARBALL_NAME="newrelic-cli_${VERSION}_${OS}_${MACHINE}.tar.gz"
CHECKSUMS_NAME="newrelic-cli_${VERSION}_checksums.txt"
RELEASE_URL="https://download.newrelic.com/install/newrelic-cli/v${VERSION}/${TARBALL_NAME}"
CHECKSUM_URL="https://download.newrelic.com/install/newrelic-cli/v${VERSION}/${CHECKSUMS_NAME}"
CHECKSUM_SIG_URL="${CHECKSUM_URL}.sig"

# NR GPG public key fingerprint: 4F9A9B5B96EC30B9 (expires 2027-04-27)
# Published on keys.openpgp.org — fetch via: gpg --keyserver hkps://keys.openpgp.org --recv-keys 4F9A9B5B96EC30B9
NR_GPG_FINGERPRINT="4F9A9B5B96EC30B9"
NR_GPG_KEYSERVER="hkps://keys.openpgp.org"

# sha256sum is GNU coreutils (Linux); macOS ships shasum instead.
sha256() {
    if command -v sha256sum > /dev/null 2>&1; then
        sha256sum "$1" | awk '{print $1}'
    elif command -v shasum > /dev/null 2>&1; then
        shasum -a 256 "$1" | awk '{print $1}'
    else
        echo ""
    fi
}

# Download tarball and checksums file.
curl -sL --retry 3 "${RELEASE_URL}" -o newrelic-cli.tar.gz
curl -sL --retry 3 "${CHECKSUM_URL}" -o "${CHECKSUMS_NAME}"

EXPECTED_HASH=$(grep "${TARBALL_NAME}" "${CHECKSUMS_NAME}" | awk '{print $1}')

if [ -z "$EXPECTED_HASH" ]; then
    echo "❌ Could not find checksum for ${TARBALL_NAME} in ${CHECKSUMS_NAME}"
    echo "   The release may be incomplete or the checksum file may be corrupt."
    exit 1
fi

ACTUAL_HASH=$(sha256 newrelic-cli.tar.gz)

if [ "$ACTUAL_HASH" != "$EXPECTED_HASH" ]; then
    echo "❌ Checksum mismatch — aborting installation."
    echo "   Expected : ${EXPECTED_HASH}"
    echo "   Actual   : ${ACTUAL_HASH}"
    echo "   The downloaded file may have been tampered with."
    exit 1
fi

echo "✅ SHA-256 verified: ${ACTUAL_HASH}"

# GPG signature verification of the checksums file.
# This confirms the checksums file was produced by New Relic's release pipeline
# and has not been tampered with since signing.
#
# Architectural note: this script is itself served from S3, so GPG verification
# here is a defence-in-depth measure against partial S3 compromise (e.g. only
# the tarball or checksums file being replaced). A full S3 compromise that also
# replaces this script would bypass it. For the strongest supply-chain guarantee,
# install via Homebrew or Chocolatey, which enforce package-manager-level signing.
if command -v gpg > /dev/null 2>&1; then
    if curl -sLf --retry 3 "${CHECKSUM_SIG_URL}" -o "${CHECKSUMS_NAME}.sig" 2>/dev/null; then
        # Import the NR public key from keys.openpgp.org if not already in the keyring.
        if ! gpg --list-keys "${NR_GPG_FINGERPRINT}" > /dev/null 2>&1; then
            gpg --keyserver "${NR_GPG_KEYSERVER}" --recv-keys "${NR_GPG_FINGERPRINT}" > /dev/null 2>&1 || true
        fi
        if gpg --verify "${CHECKSUMS_NAME}.sig" "${CHECKSUMS_NAME}" > /dev/null 2>&1; then
            echo "✅ GPG signature verified (key ${NR_GPG_FINGERPRINT})"
        else
            echo "⚠️  GPG signature verification failed for ${CHECKSUMS_NAME}"
            echo "   The checksums file could not be confirmed as genuine."
            echo "   SHA-256 integrity of the binary has still been checked."
            echo "   For a tamper-evident install, verify manually at https://github.com/newrelic/newrelic-cli/releases"
        fi
    else
        echo "⚠️  GPG signature file not available; skipping signature verification."
    fi
else
    echo "⚠️  gpg not found; skipping signature verification."
    echo "   Install gnupg for stronger supply-chain assurance."
fi

tar -xzf newrelic-cli.tar.gz

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
