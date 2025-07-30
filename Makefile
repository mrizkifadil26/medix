# Makefile to build media tools, generate data, and serve static site

# --- Directories ---
DATA_DIR        := data
RAW_DIR         := $(DATA_DIR)/raw
SYNCED_DIR      := $(DATA_DIR)/synced
ICONMAP_DIR     := $(DATA_DIR)/iconmap
ORGANIZE_DIR    := $(DATA_DIR)/organize
REPORTS_DIR     := $(DATA_DIR)/reports

OUTPUT_DIR		:= output
BIN_DIR			:= bin

# --- Executable Commands ---
SCANNER_CMD		:= ./cmd/scan
SCANNER_V2_CMD	:= ./cmd/scan-v2
PROGRESS_CMD	:= ./cmd/progress
SERVER_CMD		:= ./cmd/server
ICONMAP_CMD		:= ./cmd/iconmap
ORGANIZE_CMD	:= ./cmd/organize
SYNC_CMD    	:= ./cmd/sync
WEBGEN_CMD		:= ./cmd/webgen
DEV_CMD			:= ./cmd/dev

DEPLOY_SCRIPT	:= ./scripts/deploy.sh

# Tools
GO			:= go

# Flags
DRY_FLAG 	:= --dry
INPUT		:= data
OUTPUT   	:= dist

.PHONY: all movies tvshows index-icons sync progress webgen \
        build-webgen build-scan build-progress build-server build-index \
        build-all serve test test-slugify deploy clean help

# --- Default target ---
all: movies tvshows

# --- Media source generation ---
# media:
# 	@if [ -z "$(type)" ] || [ -z "$(name)" ]; then \
# 		echo "Usage: make media type=<type> name=<name>"; \
# 		exit 1; \
# 	fi
# 	@echo "Running for $(type)/$(name)"

# 	$(GO) run $(SCANNER_CMD) \
# 		--type media \
# 		--content $(type) \
# 		--name $(type).$(name) \
# 		--config "config/scan.media.json"

# scan-movies:
# 	$(MAKE) media type=movies name=raw
# 	$(MAKE) media type=movies name=staged
# 	$(MAKE) media type=movies name=final

# scan-tv:
# 	$(MAKE) media type=tv name=final

# scan-media: scan-movies scan-tv

# Default ARGS from Make
scan-%:
	@$(GO) run $(SCANNER_V2_CMD) \
		--config="config/scanner/$*/$(type).$(label).json" \
		--output=data/scanner/$*/$(type).$(label).json \

# Enable recursive globbing
SHELL := /bin/bash
.ONESHELL:

scan-all:
	@shopt -s globstar nullglob; \
	for f in config/scanner/**/*.json; do \
		rel_path=$${f#config/scanner/}; \
		out_path="data/scanner/$$rel_path"; \
		mkdir -p "$$(dirname "$$out_path")"; \
		echo "🔍 Scanning $$f → $$out_path"; \
		go run cmd/scan-v2/main.go --config="$$f" --output="$$out_path"; \
	done

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

# scan-icons:
# 	$(MAKE) icon type=movies name=raw
# 	$(MAKE) icon type=movies name=final
# 	$(MAKE) icon type=tv name=final

# --- Sync media and icons logically ---
sync:
	@$(GO) run $(SYNC_CMD) \
		--config="config/sync-movies.json"

# --- Organize preview/apply ---
organize-preview-%:
	@$(GO) run $(ORGANIZE_CMD) \
		--config="config/organize-$*.json" \
		--mode=preview

organize-apply-%:
	@$(GO) run $(ORGANIZE_CMD) \
		--config="config/organize-$*.json" \
		--mode=apply

organize: organize-preview organize-apply

# --- Progress report ---
progress:
	@$(GO) run $(PROGRESS_CMD)

# --- Static site generation ---
webgen:
	@$(GO) run $(WEBGEN_CMD)

dry-run:
	@$(GO) run $(WEBGEN_CMD) --input=$(INPUT) --output=$(OUTPUT) $(DRY_FLAG)

run:
	@$(GO) run $(WEBGEN_CMD) --input=$(INPUT) --output=$(OUTPUT)

# --- Build individual binaries ---
build-sync:
	@mkdir -p $(BIN_DIR)
	@$(GO) build -o $(BIN_DIR)/sync $(SYNC_CMD)

build-organize:
	@mkdir -p $(BIN_DIR)
	@$(GO) build -o $(BIN_DIR)/organize $(ORGANIZE_CMD)

build-webgen:
	mkdir -p $(BIN_DIR)
	@$(GO) build -o $(BIN_DIR)/webgen $(WEBGEN_CMD)

build-server:
	@mkdir -p $(BIN_DIR)
	@$(GO) build -o $(BIN_DIR)/server $(SERVER_CMD)

# --- Build all tools ---
build-all: build-webgen build-server build-organize

# --- Local dev ---
dev:
	@$(GO) run $(DEV_CMD)

# --- Testing ---
test:
	go test -v ./...

test-slugify:
	go test -v ./utils

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
	@echo "   make movies          Generate movies.raw.json"
	@echo "   make tv              Generate tv.raw.json"
	@echo "   make icons-movies    Generate icon index for movies"
	@echo "   make icons-tv        Generate icon index for TV shows"
	@echo "   make progress        Create progress.json"
	@echo ""
	@echo "🔄 Sync Media & Icons:"
	@echo "   make sync              Run logical sync pipeline"
	@echo ""
	@echo "🧹 Organize Media:"
	@echo "   make organize-preview-<name>  Preview organize changes"
	@echo "   make organize-apply-<name>    Apply organize changes"
	@echo ""
	@echo "🛠️ Build Commands:"
	@echo "   make build-sync        Build sync binary"
	@echo "   make build-webgen    Build static site generator binary"
	@echo "   make build-organize  Build organize binary"
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
