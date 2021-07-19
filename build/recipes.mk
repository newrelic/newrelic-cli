
#
# Makefile fragment for recipe actions
#
EMBEDDED_RECIPES_PATH         ?= internal/install/recipes/files
RECIPES_ARCHIVE_URL           ?= https://github.com/newrelic/open-install-library/releases/latest/download/recipes.zip


recipes: recipes-clean recipes-fetch

recipes-fetch:
	@echo "=== $(PROJECT_NAME) === [ recipes-fetch       ]: fetching recipes..."
	curl -sL -o ${EMBEDDED_RECIPES_PATH}/recipes.zip ${RECIPES_ARCHIVE_URL}
	@echo "=== $(PROJECT_NAME) === [ recipes-fetch       ]: extracting recipes..."
	unzip ${EMBEDDED_RECIPES_PATH}/recipes.zip -d ${EMBEDDED_RECIPES_PATH}

recipes-clean: 
	@echo "=== $(PROJECT_NAME) === [ recipes-clean       ]: cleaning recipe files..."
	find ${EMBEDDED_RECIPES_PATH} -mindepth 1 ! -regex '^${EMBEDDED_RECIPES_PATH}/.keep' -delete

.PHONY: recipes recipes-clean recipes-fetch
