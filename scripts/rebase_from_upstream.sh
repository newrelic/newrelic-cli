#!/usr/bin/env bash

# Exit script if you try to use an uninitialized variable.
set -o nounset

# Exit script if a statement returns a non-true return value.
set -o errexit

# Use the error status of the first failure,
# rather than that of the last item in a pipeline.
set -o pipefail

cd homebrew-core

echo "User email: ${GH_USER_EMAIL}"
echo "User name: ${GH_USER_NAME}"

# Set git config to our GitHub "machine user" nr-developer-toolkit
# https://developer.github.com/v3/guides/managing-deploy-keys/#machine-users
git config --global user.email "$GH_USER_EMAIL"
git config --global user.name "$GH_USER_NAME"

echo "Adding remote upstream homebrew-core..."

git branch

# Add the original Homebrew/homebrew-core as the upstream repo
git remote add upstream https://github.com/Homebrew/homebrew-core.git

# List our remotes for CI clarity
git remote -v

echo "Fetching upstream homebrew-core..."

# Need to fetch so we have the upstream/master branch locally
git fetch upstream

# Ensure our local master branch is up to date with the
# latest code from Homebrew/homebrew-core.
# Abort the rebase if encounter merge conflicts.
git rebase upstream/master

exitCode=$?

if [ $exitCode -ne 0 ]; then
  echo " "
  echo "Failed to rebase on top of upstream/master likely due to a merge conflict."
  echo "Please rebase the homebrew pull request locally and fix any conflicts before merging."
  echo " "

  git rebase --abort

  exit $exitCode
fi

git push --set-upstream origin master
