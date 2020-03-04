RELEASE_SCRIPT ?= ./tools/release.sh

GOTOOLS += github.com/goreleaser/goreleaser

REL_CMD ?= goreleaser
DIST_DIR ?= ./dist

release-clean:
	@echo "=== $(PROJECT_NAME) === [ release-clean    ]: distribution files..."
	@rm -rfv $(DIST_DIR)/*

# Example usage: make release version=0.11.0
release:
	@echo "=== $(PROJECT_NAME) === [ release          ]: Generating release."
	$(RELEASE_SCRIPT) $(version)
