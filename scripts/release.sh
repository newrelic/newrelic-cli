#!/bin/bash

COLOR_NONE='\033[0m'
COLOR_RED='\033[0;31m'
COLOR_GREEN='\033[0;32m'
COLOR_LIGHT_GREEN='\033[1;32m'

DEFAULT_BRANCH='main'
CURRENT_GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

if [ $CURRENT_GIT_BRANCH != $DEFAULT_BRANCH ]; then
  printf "\n"
  printf "${COLOR_RED} Error: Must be on main branch to create a new release. ${COLOR_NONE}"
  printf "\n"

  exit 1
fi

# Set GOBIN env variable for Go dependencies
GOBIN=$(go env GOPATH)/bin

# Install release dependencies
go install github.com/caarlos0/svu@latest
go install github.com/x-motemen/gobump/cmd/gobump@latest
go install github.com/x-motemen/gobump/cmd/gobump@latest
go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest
go install github.com/client9/misspell/cmd/misspell@latest

VER_PACKAGE="internal/version"
VER_CMD=${GOBIN}/svu
VER_BUMP=${GOBIN}/gobump
CHANGELOG_CMD=${GOBIN}/git-chglog
CHANGELOG_FILE=CHANGELOG.md
SPELL_CMD=${GOBIN}/misspell

# Compare versions
VER_CURR=$(${VER_CMD} current)
VER_NEXT=$(${VER_CMD} next)

echo ""
echo "Comparing tag versions..."
echo "Current version: ${VER_CURR}"
echo "Next version:    ${VER_NEXT}"
echo ""

if [ "${VER_CURR}" = "${VER_NEXT}" ]; then
    VER_NEXT=$(${VER_CMD} patch)

    printf "Bumping current version ${COLOR_GREEN}${VER_CURR}${COLOR_NONE} to version ${COLOR_LIGHT_GREEN}${VER_NEXT}${COLOR_NONE} for release."
fi

GIT_USER=$(git config user.name)
GIT_EMAIL=$(git config user.email)

if [ -z "${GIT_USER}" ]; then
  echo "git user.name not set"
  exit 1
fi

if [ -z "${GIT_EMAIL}" ]; then
  echo "git user.email not set"
  exit 1
fi

echo "Generating release for ${VER_NEXT} with git user ${GIT_USER}"

# Auto-generate CLI documentation
NATIVE_OS=$(go version | awk -F '[ /]' '{print $4}')
if [ -x "bin/${NATIVE_OS}/newrelic" ]; then
   rm -rf docs/cli/*
   mkdir -p docs/cli
   bin/${NATIVE_OS}/newrelic documentation --outputDir docs/cli/ --format markdown
   git add docs/cli/*

   # Commit generated docs
   git commit --no-verify -m "chore(docs): regenerate CLI docs for ${VER_NEXT}"
fi

# Auto-generate CHANGELOG updates
${CHANGELOG_CMD} --next-tag ${VER_NEXT} -o ${CHANGELOG_FILE} --sort semver

# Fix any spelling issues in the CHANGELOG
${SPELL_CMD} -source text -w ${CHANGELOG_FILE}

# Commit CHANGELOG updates
git add ${CHANGELOG_FILE}
git commit --no-verify -m "chore(changelog): update CHANGELOG for ${VER_NEXT}"
git push --no-verify origin HEAD:${DEFAULT_BRANCH}

if [ $? -ne 0 ]; then
  echo "Failed to push branch updates, exiting"
  exit 1
fi

# Create and push new tag
git tag ${VER_NEXT}
git push --no-verify origin HEAD:${DEFAULT_BRANCH} --tags

if [ $? -ne 0 ]; then
  echo "Failed to push tag, exiting"
  exit $?
fi
