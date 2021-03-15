# ==================== [START] Global Variable Declaration =================== #
SHELL := /bin/bash
# 'shell' removes newlines
BASE_DIR := $(shell pwd)

UNAME_S := $(shell uname -s)

# exports all variables
export
# ===================== [END] Global Variable Declaration ==================== #

local:
	go build ./plugin/terraform-godaddy
	rm ~/.terraform/plugins/terraform-godaddy
	mv terraform-godaddy ~/.terraform/plugins/terraform-godaddy
	chmod +x ~/.terraform/plugins/terraform-godaddy