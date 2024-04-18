#
# Makefile fragment for installing tools
#

GO               ?= go
GOFMT            ?= gofmt
GO_MOD_OUTDATED  ?= go-mod-outdated
BUILD_DIR        ?= ./bin/

# Go file to track tool deps with go modules
TOOL_DIR     ?= tools
TOOL_CONFIG  ?= $(TOOL_DIR)/tools.go

GOTOOLS ?= $(shell cd $(TOOL_DIR) && go list -e -f '{{ .Imports }}' -tags tools | tr -d '[]')

tools: check-version git-hooks
	@echo "=== $(PROJECT_NAME) === [ tools            ]: Installing tools required by the project..."
	@cd $(TOOL_DIR) && $(GO) mod download
	@cd $(TOOL_DIR) && $(GO) install $(GOTOOLS)
	@cd $(TOOL_DIR) && $(GO) mod tidy

tools-outdated: check-version
	@echo "=== $(PROJECT_NAME) === [ tools-outdated   ]: Finding outdated tool deps with $(GO_MOD_OUTDATED)..."
	@cd $(TOOL_DIR) && $(GO) list -u -m -json all | $(GO_MOD_OUTDATED) -direct -update

tools-update: check-version
	@echo "=== $(PROJECT_NAME) === [ tools-update     ]: Updating tools required by the project..."
	@cd $(TOOL_DIR) && for x in $(GOTOOLS); do \
		$(GO) get -u $$x; \
	done
	@cd $(TOOL_DIR) && $(GO) mod tidy

.PHONY: tools tools-update tools-outdated
