#!/usr/bin/env bash

##
# https://circleci.com/docs/2.0/using-shell-scripts/#set-error-flags
##

# Exit script if you try to use an uninitialized variable.
set -o nounset

# Exit script if a statement returns a non-true return value.
set -o errexit

# Use the error status of the first failure,
# rather than that of the last item in a pipeline.
set -o pipefail

printf "\n**************************************************\n"
printf "Generating Homebrew formula for git tag: ${GIT_TAG}"

asset_file=$(find ${PWD}/dist -type f -name "newrelic-cli_${GIT_TAG}_Darwin_x86_64*")

printf "\nAsset: ${asset_file}"

SHA256="$(openssl sha256 < $asset_file | sed 's/(stdin)= //')"

printf "\nNew SHA256: ${SHA256}"
printf "\n**************************************************\n\n"

homebrew_repo_name="newrelic-forks/homebrew-core"
upstream_homebrew="git@github.com:${homebrew_repo_name}.git"

printf "\nPreparing pull request to ${homebrew_repo_name}... \n"
printf "Cloning ${homebrew_repo_name}...\n"

# Set git config to our GitHub "machine user" nr-developer-toolkit
# https://developer.github.com/v3/guides/managing-deploy-keys/#machine-users
git config --global user.email "developer-toolkit-team@newrelic.com"
git config --global user.name "New Relic Developer Toolkit Bot"

# Clone homebrew-core fork
git clone $upstream_homebrew

# Change to local homebrew-core and output updates
cd homebrew-core

homebrew_formula_file='Formula/newrelic-cli.rb'
tmp_formula_file='Formula/newrelic-cli.rb.tmp'

# Set variables for lines to replace in the formula (lines 4 and 5)
formula_url='  url "https:\/\/github.com\/newrelic\/newrelic-cli\/archive\/v'${GIT_TAG}'.tar.gz"'
formula_sha256='  sha256 "'${SHA256}'"'

# Make temporary copy of existing formula file
cp $homebrew_formula_file $tmp_formula_file

# Replace lines 4 and 5 in the formula file, using the .tmp file as a template
sed -e '4s/.*/'"${formula_url}"'/' -e '5s/.*/'"${formula_sha256}"'/' $tmp_formula_file > $homebrew_formula_file

# Remove the temporary file
rm $tmp_formula_file

# Display diff (without a pager so script can continue)
git --no-pager diff

homebrew_release_branch="release/${GIT_TAG}"

# Create new branch, commit updates, push new release branch to newrelic-forks/homebrew-core
git checkout -b $homebrew_release_branch
git add Formula/newrelic-cli.rb
git status
git commit -m "newrelic-cli ${GIT_TAG}" # homebrew recommended commit message format
git push origin $homebrew_release_branch
