#!/usr/bin/env bash

# Exit script if you try to use an uninitialized variable.
set -o nounset

# Exit script if a statement returns a non-true return value.
set -o errexit

# Use the error status of the first failure,
# rather than that of the last item in a pipeline.
set -o pipefail

printf "\n**************************************************\n"

export GIT_TAG=$(git describe --tags | tr -d "v")

printf "Generating Homebrew formula for git tag: ${GIT_TAG} \n"

release_asset_url="https://github.com/newrelic/newrelic-cli/archive/v${GIT_TAG}.tar.gz"

# Download the release asset so we can get a proper sha256
wget $release_asset_url

asset_file=$(find ${PWD} -type f -name "v${GIT_TAG}.tar.gz")

printf "\nAsset gzip: ${asset_file}"

SHA256="$(openssl sha256 < $asset_file | sed 's/(stdin)= //')"

printf "\nNew SHA256: ${SHA256} \n"
printf "\n**************************************************\n"

# Change to homebrew-core repo
cd homebrew-core

# Set git config to our GitHub "machine user" nr-developer-toolkit
# https://developer.github.com/v3/guides/managing-deploy-keys/#machine-users
git config --global user.email "$GH_USER_EMAIL"
git config --global user.name "$GH_USER_NAME"

# Add the original Homebrew/homebrew-core as the upstream repo
git remote add upstream git@github.com:Homebrew/homebrew-core.git

# List our remotes for CI clarity
git remote -v

# Need to fetch so we have the upstream/master branch locally
git fetch upstream

# Ensure our local master branch is up to date with the
# latest code from Homebrew/homebrew-core.
# Abort the rebase if encounter merge conflicts.
git rebase upstream/master || git rebase --abort

homebrew_formula_file='Formula/newrelic-cli.rb'
tmp_formula_file='Formula/newrelic-cli.rb.tmp'

# Set variables for lines to replace in the formula (lines 4 and 5)
formula_url='  url "https:\/\/github.com\/newrelic\/newrelic-cli\/archive\/v'${GIT_TAG}'.tar.gz"'
formula_sha256='  sha256 "'${SHA256}'"'

# Make temporary copy of existing formula file
cp $homebrew_formula_file $tmp_formula_file

# Replace lines 4 and 5 in the formula file with updated URL and sha256
sed -e '4s/.*/'"${formula_url}"'/' -e '5s/.*/'"${formula_sha256}"'/' $tmp_formula_file > $homebrew_formula_file

# Remove the temporary file
rm $tmp_formula_file

# Display diff (without a pager so script can continue)
git --no-pager diff

# Create new branch, commit updates, push new release branch to newrelic-forks/homebrew-core
homebrew_release_branch="release/${GIT_TAG}"
git checkout -b $homebrew_release_branch
git add Formula/newrelic-cli.rb
git status
git commit -m "newrelic-cli ${GIT_TAG}" # homebrew recommended commit message format

printf "Pushing up new release branch: ${homebrew_release_branch}"

git push --set-upstream origin $homebrew_release_branch
