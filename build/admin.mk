
OUTPUT_DIR  ?= ./output/
HOMEBREW_UPSTREAM ?= git@github.com:newrelic-forks/homebrew-core.git

admin: admin-homebrew

admin-homebrew:
	@echo "=== $(PROJECT_NAME) === [ admin-homebrew ]: updating homebrew..."
	@HUB_REMOTE=$(HOMEBREW_UPSTREAM) brew bump-formula-pr --url https://github.com/newrelic/$(PROJECT_NAME)/archive/$(PROJECT_VER_TAGGED).tar.gz $(PROJECT_NAME)
