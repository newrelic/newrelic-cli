#!/usr/bin/env bash

set -e

# Environment checks
if [[ -z "${NEW_RELIC_API_KEY}" ]]; then
  echo "Please set the NEW_RELIC_API_KEY environment variable"
  exit 1
fi
if [[ -z "${NEW_RELIC_ACCOUNT_ID}" ]]; then
  echo "Please set the NEW_RELIC_ACCOUNT_ID environment variable"
  exit 1
fi
if [[ -z "${NEW_RELIC_REGION}" ]]; then
  echo "Please set the NEW_RELIC_REGION environment variable"
  exit 1
fi

# Check if we have the open install repo locally
if [ ! -d "tmp/" ]
then
    mkdir tmp/ || true
    pushd tmp
    git clone git@github.com:newrelic/open-install-library.git
    popd
fi

# Copy open install files
rsync -avzh ./tmp/open-install-library/recipes/* internal/install/recipes/files/

# Check if we have local changes to the CLI
# If so compile a new version async
CHANGES=$(git status --porcelain | wc -l)
if [ "$CHANGES" -gt 0 ]; then
    make compile-linux &
    pids[0]=$!
    make compile-darwin &
    pids[1]=$!
fi

# Build container async
docker build -t samuel:test . &
pids[2]=$!

# Remove old containers
docker rm --force cli-test &
pids[3]=$!

# Wait for everything to finish
for pid in ${pids[*]}; do
    wait $pid
done

# Start the container with CLI mounted
docker run -d \
    --name "cli-test" \
    --mount type=bind,source="$(pwd)"/bin,target=/app,readonly \
    samuel:test

# Docker exec into it
docker exec -it -e NEW_RELIC_API_KEY=$NEW_RELIC_API_KEY -e NEW_RELIC_ACCOUNT_ID=$NEW_RELIC_ACCOUNT_ID -e NEW_RELIC_REGION=$NEW_RELIC_REGION cli-test bash -c "./app/linux/newrelic install && cat /root/.newrelic/newrelic-cli.log"

# Close it again
docker rm --force cli-test
