#
# Makefile fragment for Docker actions
#
DOCKER_ORG        ?= newrelic
DOCKER_IMAGE_NAME ?= cli
DOCKER            ?= docker

docker-run: docker-image
	@echo "=== $(PROJECT_NAME) === [ docker-run       ]: running container:"
	$(DOCKER) run -it $(DOCKER_ORG)/$(DOCKER_IMAGE_NAME):$(PROJECT_VER)

docker-image: compile-linux
	@echo "=== $(PROJECT_NAME) === [ docker-image     ]: building docker image:"
	$(DOCKER) build -t $(DOCKER_ORG)/$(DOCKER_IMAGE_NAME):$(PROJECT_VER) .

docker-clean: docker-rm docker-rmi

docker-rm:
	@echo "=== $(PROJECT_NAME) === [ docker-clean     ]: removing docker containers:"
	@for i in "$$($(DOCKER) ps -a | grep "$(DOCKER_ORG)/$(DOCKER_IMAGE_NAME)" | cut -d' ' -f 1)"; do \
		if [ ! -z "$$i" ]; then \
			echo -n "=== $(PROJECT_NAME) === [ docker-clean     ]:     "; \
			$(DOCKER) rm -f $$i ; \
		fi ; \
	done

# Must cleanup containers first
docker-rmi: docker-rm
	@echo "=== $(PROJECT_NAME) === [ docker-clean     ]: removing docker images:"
	@for i in "$$($(DOCKER) images | grep "$(DOCKER_ORG)/$(DOCKER_IMAGE_NAME)" | tr -s ' ' | cut -d' ' -f 3)"; do \
		if [ ! -z "$$i" ]; then \
			echo "=== $(PROJECT_NAME) === [ docker-clean     ]:     $$i"; \
			$(DOCKER) rmi -f $$i ; \
		fi ; \
	done

.PHONY: docker-image docker-clean docker-rm docker-rmi docker-run
