package scannerV2

import (
	"context"
	"fmt"
	"io/fs"
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

	walkOpts := WalkOptions{
		MaxDepth:        options.Depth,
		SkipEmptyDirs:   options.SkipEmpty,
		OnlyLeafDirs:    options.OnlyLeaf,
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
		},
	}

	walker := NewWalker(ctx, walkOpts)

	var (
		items = make([]ScanEntry, 0)
		// jobs  []concurrency.TaskFunc
		mu sync.Mutex // to protect shared output
	)

	switch options.Mode {
	case "files":
		walker.OnVisitFile = func(path string, size int64) error {
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

	case "dirs":
		walker.OnVisitDir = func(path string, entries []os.DirEntry) error {
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

	case "mixed":
		walker.OnVisitFile = func(path string, size int64) error {
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

		walker.OnVisitDir = func(path string, entries []fs.DirEntry) error {
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

	default:
		return output, fmt.Errorf("unsupported scan mode: %s", options.Mode)
	}

	if err := walker.Walk(inputPath); err != nil {
		return output, fmt.Errorf("scan failed: %w", err)
	}

	// taskExec, err := SelectExecutor(options.Concurrency)
	// taskExec, err := SelectExecutor(1)
	// if err != nil {
	// 	return output, fmt.Errorf("concurrency error: %w", err)
	// }

	// exec := concurrency.FromTaskExecutor(taskExec)
	// if err := exec(ctx, jobs); err != nil {
	// return output, fmt.Errorf("execution error: %w", err)
	// }

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

// func (o ScanOptions) IsParallel() bool {
// 	return o.Concurrency > 1
// }

func scanFlat(base string, opts ScanOptions) []ScanEntry {
	entries, _ := os.ReadDir(base)
	var out []ScanEntry

	// Create lowercase extension filter map (same as collectSubEntries)
	extFilter := make(map[string]bool)
	// for _, ext := range opts.SubExts {
	// 	extFilter[strings.ToLower(ext)] = true
	// }

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if len(extFilter) > 0 && !extFilter[ext] {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil
		}
		size := info.Size()

		fullPath := filepath.Join(base, entry.Name())
		out = append(out, ScanEntry{
			Path: fullPath,
			Name: entry.Name(),
			Size: &size,
		})
	}

	return out
}

func scanNested(base string, opts ScanOptions) []ScanEntry {
	return scanRecursive(base, 0, opts)
}

func scanRecursive(path string, depth int, opts ScanOptions) []ScanEntry {
	// if opts.SubDepth >= 0 && depth >= opts.SubDepth {
	// 	return nil
	// }

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	// Lowercase extension filter
	extFilter := make(map[string]bool)
	// for _, ext := range opts.SubExts {
	// 	extFilter[strings.ToLower(ext)] = true
	// }

	var out []ScanEntry

	for _, entry := range entries {
		full := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			sub := scanRecursive(full, depth+1, opts)
			if opts.SkipEmpty && len(sub) == 0 {
				continue
			}

			out = append(out, ScanEntry{
				Path:     full,
				Name:     entry.Name(),
				Children: sub,
			})
		} else {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if len(extFilter) > 0 && !extFilter[ext] {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				continue // skip if we can't get info
			}
			size := info.Size()

			out = append(out, ScanEntry{
				Path: full,
				Name: entry.Name(),
				Size: &size,
			})
		}
	}

	return out
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
