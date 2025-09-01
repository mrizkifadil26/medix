# Makefile to build media tools, generate data, and serve static site

# Enable recursive globbing
SHELL := /bin/bash
.ONESHELL:

# Job types
TYPES := media icon

# --- Directories ---
CONFIG_DIR		:= config
OUTPUT_DIR		:= data
BIN_DIR			:= bin

# --- Executable Commands ---
SCANNER_CMD		:= ./cmd/scan
NORMALIZE_CMD	:= ./cmd/normalize
ENRICH_CMD		:= ./cmd/enrich
PROGRESS_CMD	:= ./cmd/progress
SERVER_CMD		:= ./cmd/server
WEBGEN_CMD		:= ./cmd/webgen
DEV_CMD			:= ./cmd/dev

DEPLOY_SCRIPT	:= ./scripts/deploy.sh

.PHONY: all movies tvshows index-icons sync progress webgen \
        build-webgen build-scan build-progress build-server build-index \
        build-all serve test test-slugify deploy clean help

# --- Default target ---
all: movies tvshows

TYPES := media icon

scan-all:
	@for type in $(TYPES); do \
		shopt -s globstar; \
		for config in $(CONFIG_DIR)/$$type/**/*.json; do \
			[ -f "$$config" ] || continue; \
			name=$$(basename $$config .json); \
			out="$(OUTPUT_DIR)/scanned/$$type/$$name.json"; \
			echo "Scanning $$type: $$config → $$out"; \
			mkdir -p "$$(dirname $$out)"; \
			$(GO) run $(SCANNER_CMD) --config "$$config" --output "$$out"; \
		done \
	done

# Scan specific type: make scan-media or make scan-icon
scan-%:
	@if [ -z "$(label)" ] || [ -z "$(type)" ]; then \
		echo "❌ Missing arguments. Usage: make scan-one label=LABEL type=TYPE"; \
		exit 1; \
	fi; \
	config="$(CONFIG_DIR)/scanner/$*/$(type).$(label).json";
	out="$(OUTPUT_DIR)/scanned/$*/$(type).$(label).json"; \

	echo "Scanning $*: $$config → $$out"; \
	mkdir -p "$$(dirname $$out)"; \
	$(GO) run $(SCANNER_CMD) --config "$$config" --output "$$out"; \

# Manual scan (no config): make scan ROOT=path MODE=files|dirs
scan:
	@$(GO) run $(SCANNER_CMD) $(ARGS)

scan-test:
	@$(GO) run $(SCANNER_CMD) \
		-config="$(CONFIG_DIR)/config.json" \
		$(ARGS)

# --- Icon index generation ---
icon:
	@if [ -z "$(type)" ] || [ -z "$(name)" ]; then \
		echo "Usage: make icon type=<type> name=<name>"; \
		exit 1; \
	fi
	@echo "Running for $(type)/$(name)"

	$(GO) run $(SCANNER_CMD) \
		--type icon \
		--content $(type) \
		--name $(type).$(name) \
		--config "config/scan.icon.json"

normalize:
	@$(GO) run $(NORMALIZE_CMD) \
		--config="config/normalizer/$(media)/$(type).$(label).json" \
		--output="output/normalized/$(media)/$(type).$(label).json"

normalize-all:
	@shopt -s globstar; \
	for config in config/normalizer/media/**/*.json; do \
		[ -f "$$config" ] || continue; \
		name=$$(basename $$config .json); \
		out="output/normalized/media/$$name.json"; \
		echo "Normalizing: $$config → $$out"; \
		mkdir -p "$$(dirname $$out)"; \
		$(GO) run $(NORMALIZE_CMD) --config "$$config" --output "$$out"; \
	done

enrich:
	@$(GO) run $(ENRICH_CMD) \
		--config="config/enricher/$(media)/$(type).$(label).json" \
		--output="output/enriched/$(media)/$(type).$(label).json"

enrich-refresh:
	@$(GO) run $(ENRICH_CMD) \
		--config="config/enricher/$(media)/$(type).$(label).json" \
		--output="output/enriched/$(media)/$(type).$(label).json" \
		--refresh

enrich-all:
	@for type in $(TYPES); do \
		shopt -s globstar; \
		for config in config/enricher/media/**/*.json; do \
			[ -f "$$config" ] || continue; \
			name=$$(basename $$config .json); \
			out="output/enriched/media/$$name.json"; \
			echo "Enriching $$type: $$config → $$out"; \
			mkdir -p "$$(dirname $$out)"; \
			$(GO) run $(ENRICH_CMD) --config "$$config" --output "$$out"; \
		done \
	done

# --- Progress report ---
progress:
	@$(GO) run $(PROGRESS_CMD)

# --- Static site generation ---
webgen:
	@$(GO) run $(WEBGEN_CMD)

# --- Build individual binaries ---
build-webgen:
	mkdir -p $(BIN_DIR)
	@$(GO) build -o $(BIN_DIR)/webgen $(WEBGEN_CMD)

build-server:
	@mkdir -p $(BIN_DIR)
	@$(GO) build -o $(BIN_DIR)/server $(SERVER_CMD)

# --- Build all tools ---
build-all: build-webgen build-server

# --- Local dev ---
dev:
	@$(GO) run $(DEV_CMD)

# --- Testing ---
test:
	go test -v ./...

# --- Deploy ---
deploy:
	bash $(DEPLOY_SCRIPT)

# --- Maintenance ---
format:
	@$(GO) fmt ./...

tidy:
	@$(GO) mod tidy

clean:
	@rm -rf $(OUTPUT_DIR) dist $(BIN_DIR)
	@echo "[DONE] Cleaned build artifacts"

# --- Help ---
help:
	@echo ""
	@echo "🎬 Media Build Commands:"
	@echo "   make progress        Create progress.json"
	@echo ""
	@echo "🛠️ Build Commands:"
	@echo "   make build-sync        Build sync binary"
	@echo "   make build-webgen    Build static site generator binary"
	@echo "   make build-all       Build all binaries into ./bin/"
	@echo ""
	@echo "🌐 Site & Serve:"
	@echo "   make webgen          Build static HTML site with webgen"
	@echo "   make dry-run         Run webgen in dry mode"
	@echo "   make run             Run webgen with configured input/output"
	@echo "   make serve           Run local web server"
	@echo "   make dev             Run dev tool"
	@echo ""
	@echo "🧪 Testing:"
	@echo "   make test            Run all unit tests"
	@echo "   make test-slugify    Test slugify utilities"
	@echo ""
	@echo "🚀 Deployment:"
	@echo "   make deploy          Deploy using deploy.sh"
	@echo ""
	@echo "🧹 Maintenance:"
	@echo "   make clean           Remove output/, dist/, and bin/"
	@echo "   make format          Format Go code"
	@echo "   make tidy            Run go mod tidy"
	@echo ""
