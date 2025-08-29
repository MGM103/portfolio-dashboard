.PHONY: build clean run

BINARY_NAME := portfolio-dashboard
OUTPUT_DIR := bin
BINARY_PATH := $(OUTPUT_DIR)/$(BINARY_NAME)
MAIN_FILE := ./cmd/portfolio-dashboard/main.go

run:
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(MAIN_FILE)
	@$(BINARY_PATH)

build:
	@echo "Building $(BINARY_NAME) executable in: $(OUTPUT_DIR)"
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(MAIN_FILE)

clean:
	@echo "Cleaning executables..."
	@rm -rf $(OUTPUT_DIR)
