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

# NR GPG key — fingerprint D3EEEEE880A91DEE209A82CC4F9A9B5B96EC30B9 (expires 2027-04-28)
# Embedded inline — no keyserver dependency, works in air-gapped/firewalled environments.
NR_GPG_FINGERPRINT="4F9A9B5B96EC30B9"
NR_PUBLIC_KEY="-----BEGIN PGP PUBLIC KEY BLOCK-----

mQINBGgPxpQBEADYHbpyxZQjZibqbc53GC6BSDR0jo93FICwv6Jj4uc8ikjzcBSP
xSLcBPtq7rKdLUCB7ccjf4W5Fvzh9uLBOKrAiEEEiNYaVS1zq1yiFtoKIUC9B5Uo
kbp3diHJl6m/erpSeICg7LbV4iDCjLanTfHufTUF+HMexTIPIaufmbBBKWZ7a57u
dHyzHsnYEa9iW3vm4WLM7Ked4Vl/eFbCa7oAkfCi22Q2MreuH+dXFqtTVQ71EfQk
j7jAxTjnxvxYVhgYHELIZptY7LqudY7p24ANF0tu1NfGqGkEnEJ6OMdcMeYIy+rR
UCuYND5ite/FEBRwB7h64utgrrNXXq4tbEJuk02+v1RXHyJXIJ3W4xgkKk+U8027
SLwNJzVx9igxKMLZlYSHgpbvq7DAZzGhFbfJHh7Mv33CayxD1bzNLhF9Pk2jYX80
YQifL3QnfpahWVA0z2CHBqqaCiWY6of7j5G8JF1xMb1MxKymTjlBLcHtvnepZ8jc
f4+bp5PE27Zkbd+qd4eI1Elc0ofillldkwh6QnGWe0qgCNFyxC0U5209JV6E7dZ5
OfG99D75JFOtEduvOFXbi/782TWgFlS4yOd6zyxxckmpgAR36z0xAZu1gsKfrc/h
KGRJYyvH15VW6bWzEfbU/bkODLCckUR2S4OoGGfg1x1TsM93gA0EvTEqhQARAQAB
tGtOZXcgUmVsaWMgKE5ldyBSZWxpYyBPcGVuIFNvdXJjZSBEZXZlbG9wZXIgVG9v
bGtpdCBHUEcgS2V5KSA8bnItZGV2ZWxvcGVyLXRvb2xraXRAdXNlcnMubm9yZXBs
eS5naXRodWIuY29tPokCVwQTAQgAQRYhBNPu7uiAqR3uIJqCzE+am1uW7DC5BQJo
D8aUAhsDBQkDwmcABQsJCAcCAiICBhUKCQgLAgQWAgMBAh4HAheAAAoJEE+am1uW
7DC5nSIQALg0+q5zM+jORY/81UA23nGB2rgtOSQjlyxIy02amoUpKWZCdke7e3H1
pjPmrU4ZCVPMuhujjtXKKi6IyCmW1McxmzoGK2CAwfBURqvcRi303XLP1/J9D8cb
Be74V6bfH39wsi5RHF7bmOwlOH3bV9o+H/VmMJzgjQw57u9sZa1VIxRCNK8KXgeB
pFBnxp7XU3w2vdmEmDHnVuRmBg6myhwM0N+popK50tghAcSx5dmEe/mz3LjLzZTU
jzOaYoSsyAoU8HCbtWFa3+6w7CD+crZy3mQBVOf0tDKkqlrs3O3xR5L3u/LNAnwZ
+0wJr9uUPfLrv10TWNX+1Ldj7hLY+Y1aoLP7HBhVW1hMoBZFHc/D9gyxjIZ7t8KR
QKwu91JoqnqO+kdBF5IhHsOXwafSSPDT9oA2Pc0onuFsh9D0sIGv/n6CjjQRPfd0
E4NZQoobP9KUiES5Qglfm8jxzdt9LcCKDnZ5kYdVoPkYYNkVGRIENu6zkyZtjBFO
GToQSwVTy/RPiRRhbijm4lTDEoUufvYvJ7ThGhCjl7ejP7xHKAB1MJddfUqeRkDx
PBpDV4lTDlxIZcv9dfZl8MvdiKwwRQ0inp41UD9z2pGS7woPXrY+MaTnJW2L5EVb
lLJAb3ebyQUIFPomxv5uczqYuXej3rXPL8QHMRH93j2W81qO720Q
=rdhy
-----END PGP PUBLIC KEY BLOCK-----"

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
        GPG_TMPDIR=$(mktemp -d)
        if echo "${NR_PUBLIC_KEY}" | gpg --homedir "${GPG_TMPDIR}" --batch --import > /dev/null 2>&1; then
            if gpg --homedir "${GPG_TMPDIR}" --batch --verify "${CHECKSUMS_NAME}.sig" "${CHECKSUMS_NAME}" > /dev/null 2>&1; then
                echo "✅ GPG signature verified (key ${NR_GPG_FINGERPRINT})"
            else
                echo "⚠️  GPG signature verification failed for ${CHECKSUMS_NAME}"
                echo "   SHA-256 integrity of the binary has still been checked."
                echo "   For a tamper-evident install, verify manually at https://github.com/newrelic/newrelic-cli/releases"
            fi
        else
            echo "⚠️  GPG key import failed; skipping signature verification."
        fi
        rm -rf "${GPG_TMPDIR}"
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
