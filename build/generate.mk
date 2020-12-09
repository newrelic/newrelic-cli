#
# Makefile fragment for Generate
#

GO           ?= go

PACKAGES ?= $(shell $(GO) list ./...)

GOTOOLS += github.com/newrelic/tutone/cmd/tutone

# Generate then lint fixes
generate: generate-run generate-tutone lint-fix

generate-tutone:
	@echo "=== $(PROJECT_NAME) === [ generate-tutone  ]: Running tutone generate..."
	@tutone -c .tutone.yml generate -l debug

generate-run: tools generate-tutone
	@echo "=== $(PROJECT_NAME) === [ generate         ]: Running generate..."
	@for p in $(PACKAGES); do \
		echo "=== $(PROJECT_NAME) === [ generate         ]:     $$p"; \
			$(GO) generate -x $$p ; \
	done

.PHONY: generate generate-run generate-tutone
