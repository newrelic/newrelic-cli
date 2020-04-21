#############################
# Global vars
#############################
PROJECT_NAME := $(shell basename $(shell pwd))
PROJECT_VER  ?= $(shell git describe --tags --always --dirty | sed -e '/^v/s/^v\(.*\)$$/\1/g')
# Last released version (not dirty) without leading v
PROJECT_VER_TAGGED  := $(shell git describe --tags --always --abbrev=0 | sed -e '/^v/s/^v\(.*\)$$/\1/g')

SRCDIR       ?= .
GO            = go

# The root module (from go.mod)
PROJECT_MODULE  ?= $(shell $(GO) list -m)

#############################
# Targets
#############################
all: build

# Humans running make:
build: check-version clean lint test cover-report compile

# Build command for CI tooling
build-ci: check-version clean lint test compile-only

# All clean commands
clean: cover-clean compile-clean release-clean

# Import fragments
include build/chocolatey.mk
include build/compile.mk
include build/deps.mk
include build/docker.mk
include build/document.mk
include build/lint.mk
include build/release.mk
include build/snapcraft.mk
include build/test.mk
include build/util.mk

.PHONY: all build build-ci clean
