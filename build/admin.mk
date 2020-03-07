
OUTPUT_DIR  ?= ./output/
LATEST_VER  ?= $(shell curl -s https://api.github.com/repos/newrelic/newrelic-cli/releases | jq -r '.[].tag_name' | head -n 1)
HOMEBREW_UPSTREAM ?= git@github.com:newrelic-forks/homebrew-core.git

admin: admin-init admin-homebrew

admin-clean:
	@echo "=== $(PROJECT_NAME) === [ admin-clean     ]: cleaning admin cache..."
	@rm -rf $(OUTPUT_DIR)/*

admin-init:
	@echo "=== $(PROJECT_NAME) === [ admin-init     ]: cloning homebrew..."
	@test -d $(OUTPUT_DIR)/homebrew || \
		git clone --depth 10 $(HOMEBREW_UPSTREAM) $(OUTPUT_DIR)/homebrew

admin-homebrew:
	@echo "=== $(PROJECT_NAME) === [ admin-homebrew ]: updating homebrew..."
	@mkdir -p $(OUTPUT_DIR)/homebrew
	@pushd $(OUTPUT_DIR)/homebrew \
		; git checkout master \
		; git fetch origin master \
		; git reset --hard origin/master  \
		; git checkout -b newrelic_$(LATEST_VER) \
		; popd
	@pushd $(OUTPUT_DIR)/homebrew; git checkout -b newrelic_$(LATEST_VER);popd
	@sed -i '' -e 's/v[0-9]*\.[0-9]*\.[0-9]*/$(LATEST_VER)/g' \
		output/homebrew/Formula/newrelic-cli.rb
	@pushd $(OUTPUT_DIR)/homebrew \
		; git add Formula/newrelic-cli.rb \
		; git commit -m 'newrelic-cli $(LATEST_VER)' \
		; git push origin +HEAD \
		; popd
	@echo "=== $(PROJECT_NAME) === [ admin-homebrew ]: Now go create a pull request"
