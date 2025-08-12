package scanner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func Scan(
	root string,
	options ScanOptions,
	outputOptions OutputOptions,
	tags []string,
) (ScanOutput, error) {
	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inputPath := filepath.Clean(root)
	output := ScanOutput{
		Version:     "0.1.0",
		GeneratedAt: time.Now().Format(time.RFC3339),
		SourcePath:  inputPath,
		Mode:        options.Mode,
	}

	// Normalize input path
	info, err := os.Stat(inputPath)
	if err != nil {
		return output, fmt.Errorf("input path error: %w", err)
	}

	if !info.IsDir() {
		return output, fmt.Errorf("input path is not a directory")
	}

	logLevel := determineLogLevel(options)

	walkOpts := WalkOptions{
		MaxDepth:        options.Depth,
		SkipEmptyDirs:   options.SkipEmpty,
		SkipRoot:        options.SkipRoot,
		OnlyLeafDirs:    options.OnlyLeaf,
		MinIncludeDepth: options.MinIncludeDepth,
		IncludeHidden:   options.IncludeHidden,
		IncludePatterns: options.IncludePatterns,
		ExcludePatterns: options.ExcludePatterns,
		IncludeExts:     options.IncludeExts,
		ExcludeExts:     options.ExcludeExts,
		StopOnError:     options.StopOnError,
		SkipOnError:     options.SkipOnError,
		EnableProgress:  options.EnableProgress,

		IncludeErrors: outputOptions.IncludeErrors,
		IncludeStats:  outputOptions.IncludeStats,

		Debug: DebugOptions{
			Enable: true,
			Level:  logLevel,
			LogFunc: func(e DebugEvent) {
				fmt.Println(logLevel)
				if shouldLog(e.Level, logLevel) {
					fmt.Printf("[%s] %s - %s (%v)\n", e.Level, e.Path, e.Message, e.Detail)
				}
			},
		},
	}

	walker := NewWalker(ctx, walkOpts)

	var (
		items = make([]ScanEntry, 0)
		// jobs  []concurrency.TaskFunc
		mu sync.Mutex // to protect shared output
	)

	handleFile := func(path string, size int64) error {
		rel, _ := filepath.Rel(inputPath, path)
		entry := ScanEntry{
			Name:       filepath.Base(path),
			Path:       path,
			RelPath:    rel,
			Size:       &size,
			Type:       "file",
			GroupPath:  filepath.Dir(rel),
			GroupLabel: strings.Split(filepath.Dir(rel), string(filepath.Separator)),
		}

		mu.Lock()
		items = append(items, entry)
		mu.Unlock()

		return nil
	}

	handleDir := func(path string, entries []os.DirEntry) error {
		rel, _ := filepath.Rel(inputPath, path)
		entry := ScanEntry{
			Path:       path,
			RelPath:    rel,
			Name:       filepath.Base(path),
			Type:       "directory",
			GroupPath:  filepath.Dir(rel),
			GroupLabel: strings.Split(filepath.Dir(rel), string(filepath.Separator)),
		}

		if options.IncludeChildren {
			entry.Children = collectChildren(path, inputPath)
		}

		mu.Lock()
		items = append(items, entry)
		mu.Unlock()

		return nil
	}

	switch options.Mode {
	case "files":
		walker.OnVisitFile = handleFile

	case "dirs":
		walker.OnVisitDir = handleDir

	case "mixed":
		walker.OnVisitFile = handleFile
		walker.OnVisitDir = handleDir

	default:
		return output, fmt.Errorf("unsupported scan mode: %s", options.Mode)
	}

	if err := walker.Walk(inputPath); err != nil {
		return output, fmt.Errorf("scan failed: %w", err)
	}

	if tags != nil {
		output.Tags = tags
	}
	output.DurationMs = time.Since(start).Milliseconds()
	output.ItemCount = len(items)
	output.Items = items

	if outputOptions.IncludeStats && walker.Stats != nil {
		output.Stats = walker.GetStats()
	}

	return output, nil
}

func collectChildren(root, base string) []ScanEntry {
	var children []ScanEntry

	entries, err := os.ReadDir(root)
	if err != nil {
		return children
	}

	for _, entry := range entries {
		fullPath := filepath.Join(root, entry.Name())
		relPath, _ := filepath.Rel(base, fullPath)

		child := ScanEntry{
			Path:       fullPath,
			RelPath:    relPath,
			Name:       entry.Name(),
			Type:       "directory",
			GroupPath:  filepath.Dir(relPath),
			GroupLabel: strings.Split(filepath.Dir(relPath), string(filepath.Separator)),
		}

		if entry.IsDir() {
			child.Type = "directory"
			child.Children = collectChildren(fullPath, base)
		} else {
			child.Type = "file"

			// Stat the file to get ModTime and Size
			info, err := entry.Info()
			fileSize := info.Size()
			if err == nil {
				child.ModTime = info.ModTime().Format(time.RFC3339)
				child.Size = &fileSize
				child.Ext = filepath.Ext(entry.Name())
			}
		}

		children = append(children, child)
	}

	return children
}
