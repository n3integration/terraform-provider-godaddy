# ==================== [START] Global Variable Declaration =================== #
SHELL := /bin/bash
# 'shell' removes newlines
BASE_DIR := $(shell pwd)

UNAME_S := $(shell uname -s)

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
	rm ~/.terraform/plugins/terraform-godaddy
	mv terraform-godaddy ~/.terraform/plugins/terraform-godaddy
	chmod +x ~/.terraform/plugins/terraform-godaddy
