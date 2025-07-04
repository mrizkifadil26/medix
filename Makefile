# Makefile to generate sidebar JSON files

OUTPUT_DIR = output
SCANNER = ./cmd/scanner
PROGRESS = ./cmd/progress

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
	rm -f $(OUTPUT_DIR)/movies_sidebar.json $(OUTPUT_DIR)/tvshows_sidebar.json

# Create progress.json report
.PHONY: progress
progress:
	go run $(PROGRESS)

.PHONY: build-static
build-static:
	go run ./cmd/builder
