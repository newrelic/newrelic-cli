
#
# Makefile fragment for events actions
#
EMBEDDED_EVENTS_PATH         ?= internal/install/recipes/files/events.src

events: events-clean events-init

events-init:
	@echo "=== $(PROJECT_NAME) === [ events-init    ]: initializing..."
	@if [ -n "${SEGMENT_WRITE_KEY}" ]; then \
		echo "${SEGMENT_WRITE_KEY}" > ${EMBEDDED_EVENTS_PATH}; \
	fi

events-clean:
	@echo "=== $(PROJECT_NAME) === [ events-clean    ]: cleaning events..."
	@rm -f $(EMBEDDED_EVENTS_PATH)

.PHONY: events events-clean events-init
