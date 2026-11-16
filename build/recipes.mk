
#
# Makefile fragment for recipe actions
#
EMBEDDED_RECIPES_PATH         ?= internal/install/recipes/files
RECIPE_ARCHIVE_VERSION        ?= $(shell curl -Ls https://s3.us-east-1.amazonaws.com/nr-downloads-main/install/open-install-library/currentVersion.txt)
RECIPE_ARCHIVE_URL            ?= https://s3.us-east-1.amazonaws.com/nr-downloads-main/install/open-install-library/${RECIPE_ARCHIVE_VERSION}/recipes.zip
RECIPE_ARCHIVE_SHA256_URL     ?= https://s3.us-east-1.amazonaws.com/nr-downloads-main/install/open-install-library/${RECIPE_ARCHIVE_VERSION}/recipes.zip.sha256

recipes: recipes-clean recipes-fetch

recipes-fetch:
	@echo "=== $(PROJECT_NAME) === [ recipes-fetch    ]: fetching recipes $(RECIPE_ARCHIVE_VERSION)..."
	@curl -sL -o ${EMBEDDED_RECIPES_PATH}/recipes.zip ${RECIPE_ARCHIVE_URL}
	@echo "=== $(PROJECT_NAME) === [ recipes-fetch    ]: verifying recipes integrity..."
	@if curl -sLf -o ${EMBEDDED_RECIPES_PATH}/recipes.zip.sha256 ${RECIPE_ARCHIVE_SHA256_URL} 2>/dev/null; then \
		EXPECTED=$$(cat ${EMBEDDED_RECIPES_PATH}/recipes.zip.sha256 | awk '{print $$1}'); \
		ACTUAL=$$(sha256sum ${EMBEDDED_RECIPES_PATH}/recipes.zip | awk '{print $$1}'); \
		if [ "$$EXPECTED" != "$$ACTUAL" ]; then \
			echo "=== $(PROJECT_NAME) === [ recipes-fetch    ]: ERROR: SHA-256 mismatch for recipes.zip"; \
			echo "  expected: $$EXPECTED"; \
			echo "  actual  : $$ACTUAL"; \
			exit 1; \
		fi; \
		echo "=== $(PROJECT_NAME) === [ recipes-fetch    ]: SHA-256 verified: $$ACTUAL"; \
		rm -f ${EMBEDDED_RECIPES_PATH}/recipes.zip.sha256; \
	else \
		echo "=== $(PROJECT_NAME) === [ recipes-fetch    ]: WARNING: recipes.zip.sha256 not available; skipping integrity check"; \
	fi
	@echo "=== $(PROJECT_NAME) === [ recipes-fetch    ]: extracting recipes..."
	@unzip ${EMBEDDED_RECIPES_PATH}/recipes.zip -d ${EMBEDDED_RECIPES_PATH}

recipes-clean:
	@echo "=== $(PROJECT_NAME) === [ recipes-clean    ]: cleaning recipe files..."
	@find ${EMBEDDED_RECIPES_PATH} -mindepth 1 ! -regex '^${EMBEDDED_RECIPES_PATH}/.keep' -delete

.PHONY: recipes recipes-clean recipes-fetch
