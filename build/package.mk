#
# Makefile fragment for CloudFormation Registry
#

FIND_CMD    ?= find
ZIP         ?= zip
BUILD_DIR   ?= ./bin
OUT_DIR     ?= package/
COMPILE_OS  ?= darwin linux windows

package-clean:
	@echo "=== $(PROJECT_NAME) === [ package-clean    ]: removing stale packages..."
	@rm -rfv $(BUILD_DIR)/$(OUT_DIR)/*

package: clean package-zip

package-zip: package-clean compile-all
	@echo "=== $(PROJECT_NAME) === [ package          ]: creating packages..."
	@mkdir -p $(BUILD_DIR)/$(OUT_DIR)
	@for b in $(BINS); do \
		for os in $(COMPILE_OS); do \
			PACKAGE_FILES=`$(FIND_CMD) $(BUILD_DIR)/$$os/ -type f` ; \
			ZIP_FILE=$(BUILD_DIR)/$(OUT_DIR)/$(PROJECT_NAME)-$$os-$(PROJECT_VER) ; \
			echo "=== $(PROJECT_NAME) === [ package          ]:     $${ZIP_FILE}.zip"; \
			$(ZIP) -qj $${ZIP_FILE}.zip $$PACKAGE_FILES ; \
		done \
	done


.PHONY: package package-clean package-only package-darwin package-linux package-windows
