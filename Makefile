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

# Run unit tests (only in util/)
.PHONY: test
test:
	go test -v ./...

.PHONY: build serve watch watch-serve

build:
	go run ./cmd/builder

serve:
	@echo "📡 Serving dist/ at http://localhost:8080"
	@go run ./cmd/server/main.go

watch:
	@echo "🔁 Watching files and building with Air..."
	@air

watch-serve:
	@echo "🌐 Watching & serving dist/ (in background)..."
	@make -j2 watch serve


# Deploy using shell script
.PHONY: deploy
deploy:
	bash $(DEPLOY_SCRIPT)