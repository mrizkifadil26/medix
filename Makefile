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

DEPLOY_SCRIPT  		= ./scripts/deploy.sh

.PHONY: all movies tvshows index-icons progress webgen \
        build-webgen build-scan build-progress build-server build-index \
        build-all serve watch watch-serve test test-slugify deploy clean help

# --- Default target ---
all: movies tvshows

# --- Media source generation ---
movies:
	go run $(SCANNER_CMD) -config "config/scan_config.json" -type movies

tvshows:
	go run $(SCANNER_CMD) -config "config/scan_config.json" -type tvshows

# --- Icon indexing ---
index-icons:
	go run $(ICON_INDEXER_CMD)

# --- Progress report ---
progress:
	go run $(PROGRESS_CMD)

# --- Static site generation ---
webgen:
	go run $(WEBGEN_CMD)

# --- Build individual binaries ---
build-webgen:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/webgen $(WEBGEN_CMD)

# --- Build all tools ---
build-all: build-webgen

# --- Local server ---
dev:
	go run ./cmd/dev

# --- File watching ---
watch:
	@echo "🔁 Watching files and building with Air..."
	@air

watch-serve:
	@echo "🌐 Watching & serving dist/ (in background)..."
	@make -j2 watch serve

# --- Testing ---
test:
	go test -v ./...

test-slugify:
	go test -v ./util

# --- Deploy ---
deploy:
	bash $(DEPLOY_SCRIPT)

# --- Clean ---
clean:
	rm -rf $(OUTPUT_DIR) dist $(BIN_DIR)

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
