# Makefile to generate sidebar JSON files and run project tools

OUTPUT_DIR     = output
SCANNER        = ./cmd/scanner
PROGRESS       = ./cmd/progress
BUILDER        = ./cmd/builder
SERVER         = ./cmd/server
ICON_INDEXER   = ./cmd/indexer
DEPLOY_SCRIPT  = ./scripts/deploy.sh

.PHONY: all movies tvshows index-icons progress build serve watch watch-serve test test-slugify deploy clean help

# Default target
all: movies tvshows

# --- Media source generation ---
movies:
	go run $(SCANNER) movies

tvshows:
	go run $(SCANNER) tvshows

# --- Icon indexing ---
index-icons:
	go run $(ICON_INDEXER)

# --- Progress report ---
progress:
	go run $(PROGRESS)

# --- Build and serve ---
build:
	go run $(BUILDER)

serve:
	@echo "📡 Serving dist/ at http://localhost:8080"
	@go run ./cmd/server/main.go

watch:
	@echo "🔁 Watching files and building with Air..."
	@air

watch-serve:
	@echo "🌐 Watching & serving dist/ (in background)..."
	@make -j2 watch serve

# --- Testing ---
test:
	go test -v ./util

test-slugify:
	go test -v ./util

# --- Deploy ---
deploy:
	bash $(DEPLOY_SCRIPT)

# --- Clean ---
clean:
	rm -rf $(OUTPUT_DIR)

# --- Help ---
help:
	@echo ""
	@echo "🎬 Media Build Commands:"
	@echo "  make movies          Generate movies_sidebar.json"
	@echo "  make tvshows         Generate tvshows_sidebar.json"
	@echo "  make index-icons     Generate icon index JSON"
	@echo "  make progress        Create progress.json"
	@echo "  make build           Run builder"
	@echo "  make serve           Run local web server"
	@echo "  make watch           Watch & rebuild with Air"
	@echo "  make watch-serve     Watch and serve concurrently"
	@echo "  make test            Run all unit tests"
	@echo "  make test-slugify    Test slugify only"
	@echo "  make deploy          Deploy using deploy.sh"
	@echo "  make clean           Remove output/"
	@echo ""
