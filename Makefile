# Makefile to build media tools, generate data, and serve static site
# --- Directories ---
OUTPUT_DIR        := output
BIN_DIR           := bin

# --- Command Sources ---
SCANNER_CMD        	= ./cmd/scan
PROGRESS_CMD       	= ./cmd/progress
SERVER_CMD         	= ./cmd/server
ICON_INDEXER_CMD   	= ./cmd/index
WEBGEN_CMD        	= ./cmd/webgen
DEV_CMD             := ./cmd/dev

DEPLOY_SCRIPT  		= ./scripts/deploy.sh

# Tools
GO       = go
RICHGO   = richgo

# Flags
DRY_FLAG = --dry
INPUT    = data
OUTPUT   = dist

.PHONY: all movies tvshows index-icons progress webgen \
        build-webgen build-scan build-progress build-server build-index \
        build-all serve test test-slugify deploy clean help

# --- Default target ---
all: movies tvshows

# --- Media source generation ---
movies:
	@$(GO) run $(SCANNER_CMD) -config "config/scan_config.json" -type movies

tv:
	@$(GO) run $(SCANNER_CMD) -config "config/scan_config.json" -type tv

# --- Icon indexing ---
index-icons:
	@$(GO) run $(ICON_INDEXER_CMD)

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
	@$(RICHGO) run $(DEV_CMD)

# --- Testing ---
test:
	go test -v ./...

test-slugify:
	go test -v ./util

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
	@echo "   make movies          Generate movies_sidebar.json"
	@echo "   make tvshows         Generate tvshows_sidebar.json"
	@echo "   make index-icons     Generate icon index JSON"
	@echo "   make progress        Create progress.json"
	@echo ""
	@echo "🛠️ Build Commands:"
	@echo "   make build-webgen    Build static site generator binary"
	@echo "   make build-all       Build all binaries into ./bin/"
	@echo ""
	@echo "🌐 Site & Serve:"
	@echo "   make webgen          Build static HTML site with webgen"
	@echo "   make serve           Run local web server"
	@echo "   make watch           Watch & rebuild with Air"
	@echo "   make watch-serve     Watch and serve concurrently"
	@echo ""
	@echo "🧪 Testing:"
	@echo "   make test            Run all unit tests"
	@echo "   make test-slugify    Test slugify only"
	@echo ""
	@echo "🚀 Deployment:"
	@echo "   make deploy          Deploy using deploy.sh"
	@echo ""
	@echo "🧹 Maintenance:"
	@echo "   make clean           Remove output/"
	@echo ""
