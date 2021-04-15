#!/usr/bin/env bash

if [ -z "$VERSION" ]; then
    VERSION=$(curl -sH "Authorization: token $GITHUB_TOKEN" https://api.github.com/repos/newrelic/open-install-library/releases | jq -r '.[0].tag_name')
fi

# Run the script in a temporary directory that we know is empty.
SCRATCH="tmp"

function error {
  echo "An error occurred fetching the recipes."
  echo "The contents of the directory $SCRATCH have been left in place to help to debug the issue."
}

trap error ERR

RELEASE_URL="https://github.com/newrelic/open-install-library/archive/refs/tags/${VERSION}.tar.gz"

rm -rf $SCRATCH
rm -rf recipes/
mkdir $SCRATCH
curl -sL --retry 3 "${RELEASE_URL}" | tar -xz -C $SCRATCH
mv ${SCRATCH}/open-install-library-*/recipes/ .

