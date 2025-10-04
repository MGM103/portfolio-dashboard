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
	@echo "Deploying $(BINARY_NAME) and wrapper script to $(INSTALL_DIR)"
	@mkdir -p $(INSTALL_DIR)
	@cp $(BINARY_PATH) $(INSTALL_DIR)/$(BINARY_NAME)
	@install -m 755 ./portfolio-app $(INSTALL_DIR)/portfolio-app
	@mkdir -p $(CONFIG_DIR)
	@echo "Config dir created at location: $(CONFIG_DIR)"

	@if [ ! -f $(CONFIG_DIR)/.env ]; then \
		echo "Creating skeleton .env at $(CONFIG_DIR)/.env, please enter there values."; \
		{ \
			echo "CMC_API_KEY="; \
			echo "DB_PATH="; \
		} > $(CONFIG_DIR)/.env; \
		chmod 600 $(CONFIG_DIR)/.env; \
	else \
		echo "$(CONFIG_DIR)/.env already exists, not overwriting."; \
	fi

clean:
	@echo "Cleaning executables..."
	@rm -rf $(OUTPUT_DIR)
