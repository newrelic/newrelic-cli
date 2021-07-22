#
# Makefile fragment for Generate
#

GO           ?= go

PACKAGES ?= $(shell $(GO) list ./...)

GOTOOLS += github.com/newrelic/tutone/cmd/tutone

# Generate then lint fixes
generate: tools generate-tutone lint-fix

generate-tutone:
	@echo "=== $(PROJECT_NAME) === [ generate-tutone  ]: Running tutone generate..."
	@tutone -c .tutone.yml generate -l debug

.PHONY: generate generate-tutone
