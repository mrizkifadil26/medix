# Makefile to generate sidebar JSON files

OUTPUT_DIR = output
SCANNER = ./cmd/scanner
PROGRESS = ./cmd/progress
BUILDER := ./cmd/builder
SERVER := ./cmd/server
DEPLOY_SCRIPT := ./scripts/deploy.sh

# Default target
.PHONY: all
all: movies tvshows

# Target: generate movies_sidebar.json
.PHONY: movies
movies:
	go run $(SCANNER) movies

# Target: generate tvshows_sidebar.json
.PHONY: tvshows
tvshows:
	go run $(SCANNER) tvshows

# Clean output directory
.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)

# Create progress.json report
.PHONY: progress
progress:
	go run $(PROGRESS)

.PHONY: build
build:
	go run $(BUILDER)

# Serve static site
.PHONY: serve
serve:
	go run $(SERVER)

# Watch for changes and rebuild (requires reflex)
.PHONY: watch
watch:
	@which reflex > /dev/null || (echo "reflex not installed. Install with: go install github.com/cespare/reflex@latest" && exit 1)
	reflex -r '\.go$$|\.tmpl$$|\.json$$|\.css$$|\.js$$' -R '^dist/' -- make build

# Optional: watch and serve (requires entr or fswatch + Make tweaks)
.PHONY: watch-serve
watch-serve:
	@echo "👀 Watching for changes and rebuilding..."
	reflex -r '\.go$$|\.tmpl$$|\.json$$|\.css$$|\.js$$' -R '^dist/' -- sh -c 'make build && make serve'

# Deploy using shell script
.PHONY: deploy
deploy:
	bash $(DEPLOY_SCRIPT)