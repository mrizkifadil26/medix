# Makefile to generate sidebar JSON files

OUTPUT_DIR = output
SCANNER = ./cmd/scanner
PROGRESS = ./cmd/progress
BUILDER := ./cmd/builder
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

.PHONY: watch
watch:
	go run $(BUILDER) --watch

# Deploy using shell script
.PHONY: deploy
deploy:
	bash $(DEPLOY_SCRIPT)