# ==================== [START] Global Variable Declaration =================== #
SHELL := /bin/bash
# 'shell' removes newlines
BASE_DIR := $(shell pwd)

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
	
local:
	go build ./plugin/terraform-godaddy
	rm -rf ~/.terraform/plugins/terraform-godaddy
	rm -rf ~/.terraform.d/plugins/github.com/n3integration/godaddy/$(VERSION)/darwin_amd64/terraform-provider-godaddy
	mkdir -p ~/.terraform.d/plugins/github.com/n3integration/godaddy/$(VERSION)/darwin_amd64/
	mv terraform-godaddy ~/.terraform.d/plugins/github.com/n3integration/godaddy/$(VERSION)/darwin_amd64/terraform-provider-godaddy
	chmod +x ~/.terraform.d/plugins/github.com/n3integration/godaddy/$(VERSION)/darwin_amd64/terraform-provider-godaddy

