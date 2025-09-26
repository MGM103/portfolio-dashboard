.PHONY: build clean run install

BINARY_NAME := portfolio-dashboard
OUTPUT_DIR := bin
BINARY_PATH := $(OUTPUT_DIR)/$(BINARY_NAME)
MAIN_FILE := ./cmd/portfolio-dashboard/main.go
INSTALL_DIR := $(HOME)/.local/bin

run:
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(MAIN_FILE)
	@$(BINARY_PATH)

build:
	@echo "Building $(BINARY_NAME) executable in: $(OUTPUT_DIR)"
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(MAIN_FILE)

install: build
	@echo "Deploying $(BINARY_NAME) to $(INSTALL_DIR)"
	@mkdir -p $(INSTALL_DIR)
	cp $(BINARY_PATH) $(INSTALL_DIR)/$(BINARY_NAME)

clean:
	@echo "Cleaning executables..."
	@rm -rf $(OUTPUT_DIR)
