RELEASE_SCRIPT ?= ./scripts/release.sh

GOTOOLS += github.com/goreleaser/goreleaser

REL_CMD ?= goreleaser
DIST_DIR ?= ./dist

HOMEBREW_CMD ?= brew
HOMEBREW_UPSTREAM ?= git@github.com:newrelic-forks/homebrew-core.git
ARCHIVE_URL       ?= https://github.com/newrelic/$(strip $(PROJECT_NAME))/archive/v$(strip $(PROJECT_VER_TAGGED)).tar.gz

# Example usage: make release version=0.11.0
release:
	@echo "=== $(PROJECT_NAME) === [ release          ]: Generating release."
	$(RELEASE_SCRIPT) $(version)

release-clean:
	@echo "=== $(PROJECT_NAME) === [ release-clean    ]: distribution files..."
	@rm -rfv $(DIST_DIR)/*

release-publish: clean tools docker-login
	@echo "=== $(PROJECT_NAME) === [ release-publish  ]: Publishing release via $(REL_CMD)"
	$(REL_CMD)

release-homebrew:
ifeq ($(HOMEBREW_GITHUB_API_TOKEN), "")
	@echo "=== $(PROJECT_NAME) === [ admin-homebrew ]: HOMEBREW_GITHUB_API_TOKEN must be set"
	exit 1
endif
ifeq ($(shell which $(HOMEBREW_CMD)), "")
	@echo "=== $(PROJECT_NAME) === [ admin-homebrew ]: Hombrew command '$(HOMEBREW_CMD)' not found."
	exit 1
endif
	@echo "=== $(PROJECT_NAME) === [ admin-homebrew ]: updating homebrew..."
	@HUB_REMOTE=$(HOMEBREW_UPSTREAM) $(HOMEBREW_CMD) bump-formula-pr --url $(ARCHIVE_URL) $(PROJECT_NAME)
.PHONY: release release-clean release-homebrew release-publish
