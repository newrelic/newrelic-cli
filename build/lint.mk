#
# Makefile fragment for Linting
#

GO           ?= go
MISSPELL     ?= misspell
GOFMT        ?= gofmt
GOIMPORTS    ?= goimports

COMMIT_LINT_CMD   ?= go-gitlint
COMMIT_LINT_REGEX ?= "(Scoop update|(chore|docs|feat|fix|refactor|tests?)(\([^\)]+\))?:) .*"
COMMIT_LINT_START ?= "2021-02-23"
COMMIT_MSG_FILE   ?= ""

GOLINTER      = golangci-lint

EXCLUDEDIR      ?= .git
SRCDIR          ?= .
GO_PKGS         ?= $(shell ${GO} list ./... | grep -v -e "/vendor/" -e "/example")
FILES           ?= $(shell find ${SRCDIR} -type f \( -name '*.go' -o -name '*.yml' -o -name '*.json' -o -name '*.mk' \))
GO_FILES        ?= $(shell find $(SRCDIR) -type f -name "*.go" | grep -v -e ".git/" -e '/vendor/' -e '/example/')
PROJECT_MODULE  ?= $(shell $(GO) list -m)

GO_MOD_OUTDATED ?= go-mod-outdated

GOTOOLS += github.com/client9/misspell/cmd/misspell \
           github.com/llorllale/go-gitlint/cmd/go-gitlint \
           github.com/psampaz/go-mod-outdated \
           github.com/golangci/golangci-lint/cmd/golangci-lint \
           golang.org/x/tools/cmd/goimports \
		   gotest.tools/gotestsum

lint: deps spell-check gofmt lint-commit golangci goimports outdated tools-outdated
lint-fix: deps gofmt-fix goimports

#
# Check spelling on all the files, not just source code
#
spell-check: deps
	@echo "=== $(PROJECT_NAME) === [ spell-check      ]: Checking for spelling mistakes with $(MISSPELL)..."
	@$(MISSPELL) -source text $(FILES)

spell-check-fix: deps
	@echo "=== $(PROJECT_NAME) === [ spell-check-fix  ]: Fixing spelling mistakes with $(MISSPELL)..."
	@$(MISSPELL) -source text -w $(FILES)

gofmt: deps
	@echo "=== $(PROJECT_NAME) === [ gofmt            ]: Checking file format with $(GOFMT)..."
	@$(GOFMT) -e -l -s -d $(GO_FILES)

gofmt-fix: deps
	@echo "=== $(PROJECT_NAME) === [ gofmt-fix        ]: Fixing file format with $(GOFMT)..."
	@$(GOFMT) -e -l -s -w $(GO_FILES)

goimports: deps
	@echo "=== $(PROJECT_NAME) === [ goimports        ]: Checking imports with $(GOIMPORTS)..."
	@$(GOIMPORTS) -l -w -local $(PROJECT_MODULE) $(GO_FILES)

lint-commit: deps
	@echo "=== $(PROJECT_NAME) === [ lint-commit      ]: Checking that commit messages are properly formatted ($(COMMIT_LINT_CMD))..."
	@$(COMMIT_LINT_CMD) --since=$(COMMIT_LINT_START) --subject-minlen=10 --subject-maxlen=120 --subject-regex=$(COMMIT_LINT_REGEX) --msg-file=$(COMMIT_MSG_FILE)

golangci: deps
	@echo "=== $(PROJECT_NAME) === [ golangci-lint    ]: Linting using $(GOLINTER) ($(COMMIT_LINT_CMD))..."
	@$(GOLINTER) run --build-tags unit,integration

outdated: deps
	@echo "=== $(PROJECT_NAME) === [ outdated         ]: Finding outdated deps with $(GO_MOD_OUTDATED)..."
	@$(GO) list -u -m -json all | $(GO_MOD_OUTDATED) -direct -update

.PHONY: lint spell-check spell-check-fix gofmt gofmt-fix lint-fix lint-commit outdated goimports
