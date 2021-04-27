.PHONY: linux docs local
# ==================== [START] Global Variable Declaration =================== #
SHELL := /bin/bash
# 'shell' removes newlines
BASE_DIR := $(shell pwd)

BINARY := "terraform-provider-godaddy"

COMMIT := $(shell git rev-parse --short HEAD)

UNAME_S := $(shell uname -s)

VERSION := $(shell grep "version=" install.sh | cut -d= -f2)

# exports all variables
export
# ===================== [END] Global Variable Declaration ==================== #

linux:
	@echo "Pulling latest image"
	@docker-compose -f "${BASE_DIR}/docker/docker-compose.yml" pull
	@echo "Compile and build"
	@docker-compose -f "${BASE_DIR}/docker/docker-compose.yml" run --rm builder
	@echo "Cleaning up resources"
	@docker-compose -f "${BASE_DIR}/docker/docker-compose.yml" down
	
docs:
	@go generate

local:
	go build -o $(BINARY) -ldflags='-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)' .
	rm -rf ~/.terraform/plugins/terraform-godaddy
	rm -rf ~/.terraform.d/plugins/github.com/n3integration/godaddy/$(VERSION)/darwin_amd64/terraform-provider-godaddy
	mkdir -p ~/.terraform.d/plugins/github.com/n3integration/godaddy/$(VERSION)/darwin_amd64/
	mv $(BINARY) ~/.terraform.d/plugins/github.com/n3integration/godaddy/$(VERSION)/darwin_amd64/
	chmod +x ~/.terraform.d/plugins/github.com/n3integration/godaddy/$(VERSION)/darwin_amd64/$(BINARY)

