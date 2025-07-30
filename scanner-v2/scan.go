package scannerV2

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Scan(root string, options ScanOptions) (ScanOutput, error) {
	start := time.Now()
	output := ScanOutput{
		GeneratedAt: time.Now().Format(time.RFC3339),
		SourcePath:  root,
		Mode:        options.Mode,
	}

	// Normalize input path
	inputPath := filepath.Clean(root)

	// Check if input exists
	info, err := os.Stat(inputPath)
	if err != nil {
		return output, fmt.Errorf("input path error: %w", err)
	}

	if !info.IsDir() {
		return output, fmt.Errorf("input path is not a directory")
	}

	// Normalize exclude list
	excludeMap := make(map[string]bool)
	for _, ex := range options.Exclude {
		excludeMap[filepath.ToSlash(filepath.Clean(ex))] = true
	}

	// Helper: check if path is excluded
	isExcluded := func(absPath string) bool {
		rel, err := filepath.Rel(root, absPath)
		if err != nil {
			return false // fallback: don’t exclude
		}

		rel = filepath.ToSlash(rel) // make sure slashes match

		for prefix := range excludeMap {
			if strings.HasPrefix(rel, prefix) {
				return true
			}
		}
		return false
	}

	switch options.Mode {
	case "files":
		err = WalkFiles(inputPath, WalkOptions{
			MaxDepth: options.Depth,
			Exts:     options.Exts,
			Verbose:  options.Verbose,
		}, func(path string, size int64) {
			excludedCount := 0
			if isExcluded(path) {
				excludedCount++
				return
			}

			rel, _ := filepath.Rel(inputPath, path)
			entry := ScanEntry{
				// Source:     inputPath,
				GroupPath:  filepath.Dir(rel),
				ItemPath:   path,
				ItemName:   filepath.Base(path),
				ItemSize:   &size,
				GroupLabel: strings.Split(filepath.Dir(rel), string(filepath.Separator)),
			}

			output.ItemCount = len(output.Items)
			output.ExcludedCount = excludedCount
			output.Duration = time.Since(start).String()
			output.Items = append(output.Items, entry)
		})

	case "dirs":
		err = WalkDirs(
			inputPath,
			WalkOptions{
				MaxDepth:  options.Depth,
				OnlyLeaf:  options.OnlyLeaf,
				LeafDepth: options.LeafDepth,
				SkipEmpty: options.SkipEmpty,
				Verbose:   options.Verbose,
			},
			func(path string, entries []os.DirEntry) {
				excludedCount := 0
				if isExcluded(path) {
					excludedCount++
					return
				}

				rel, _ := filepath.Rel(inputPath, path)
				entry := ScanEntry{
					// Source:     inputPath,
					GroupPath: filepath.Dir(rel),
					ItemPath:  path,
					ItemName:  filepath.Base(path),
					// ItemSize:   nil,
					GroupLabel: strings.Split(filepath.Dir(rel), string(filepath.Separator)),
					SubEntries: func() []string {
						var subs []string
						for _, e := range entries {
							if e.IsDir() {
								subs = append(subs, e.Name())
							}
						}

						if len(subs) > 0 {
							return subs
						}

						return nil
					}(),
				}

				output.ItemCount = len(output.Items)
				output.ExcludedCount = excludedCount
				output.Duration = time.Since(start).String()
				output.Items = append(output.Items, entry)
			})

	default:
		return output, fmt.Errorf("unsupported scan mode: %s", options.Mode)
	}

	if err != nil {
		return output, fmt.Errorf("scan failed: %w", err)
	}

	return output, nil
}
