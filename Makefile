.PHONY: build clean run install

BINARY_NAME := portfolio-dashboard
OUTPUT_DIR := bin
BINARY_PATH := $(OUTPUT_DIR)/$(BINARY_NAME)
MAIN_FILE := ./cmd/portfolio-dashboard/main.go
INSTALL_DIR := $(HOME)/.local/bin
CONFIG_DIR := $(HOME)/.config/portfolio-dashboard

run: build
	@$(BINARY_PATH)

build:
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(MAIN_FILE)

install: build
	@echo "Deploying $(BINARY_NAME) and accompanying script to $(INSTALL_DIR)"
	@mkdir -p $(INSTALL_DIR)
	@cp $(BINARY_PATH) $(INSTALL_DIR)/$(BINARY_NAME)
	@install -m 755 ./portfolio-app $(INSTALL_DIR)/portfolio-app
	@echo "Updating dependencies"
	@mkdir -p $(CONFIG_DIR)
	@test -f $(CONFIG_DIR)/asset.db || cp ./data/asset.db $(CONFIG_DIR)/asset.db
	@test -f $(CONFIG_DIR)/.env || cp .env $(CONFIG_DIR)/.env

clean:
	@echo "Cleaning executables..."
	@rm -rf $(OUTPUT_DIR)
