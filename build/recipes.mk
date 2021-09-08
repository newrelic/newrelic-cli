
#
# Makefile fragment for recipe actions
#
EMBEDDED_RECIPES_PATH         ?= internal/install/recipes/files
RECIPE_ARCHIVE_VERSION        ?= $(shell curl -Ls https://s3.us-east-1.amazonaws.com/nr-downloads-main/install/open-install-library/currentVersion.txt)
RECIPE_ARCHIVE_URL            ?= https://download.newrelic.com/install/open-install-library/${RECIPE_ARCHIVE_VERSION}/recipes.zip

recipes: recipes-clean recipes-fetch

recipes-fetch:
	@echo "=== $(PROJECT_NAME) === [ recipes-fetch       ]: fetching recipes..."
	curl -sL -o ${EMBEDDED_RECIPES_PATH}/recipes.zip ${RECIPE_ARCHIVE_URL}

	@echo "=== $(PROJECT_NAME) === [ recipes-fetch       ]: extracting recipes..."
	unzip ${EMBEDDED_RECIPES_PATH}/recipes.zip -d ${EMBEDDED_RECIPES_PATH}

recipes-clean: 
	@echo "=== $(PROJECT_NAME) === [ recipes-clean       ]: cleaning recipe files..."
	find ${EMBEDDED_RECIPES_PATH} -mindepth 1 ! -regex '^${EMBEDDED_RECIPES_PATH}/.keep' -delete

.PHONY: recipes recipes-clean recipes-fetch
