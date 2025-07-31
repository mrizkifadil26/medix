package scannerV2

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type WalkOptions struct {
	MaxDepth  int
	Exts      []string
	OnlyLeaf  bool // only include leaf directories
	LeafDepth int  // 0 = default, 1 = leaf-1, 2 = leaf-2, etc.
	SkipEmpty bool // skip empty directories
	Verbose   bool // log visited/skipped folders
}

func WalkDirs(root string, opts WalkOptions, fn func(path string, entries []os.DirEntry)) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if opts.Verbose {
				fmt.Printf("⚠️  Error accessing %s: %v\n", path, err)
			}

			return err
		}
		if !d.IsDir() {
			return nil
		}

		depth := pathDepth(root, path)
		if opts.MaxDepth >= 0 && depth > opts.MaxDepth {
			if opts.Verbose {
				fmt.Printf("⏭️  Skipping %s (too deep, depth %d)\n", path, depth)
			}

			return fs.SkipDir
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			if opts.Verbose {
				fmt.Printf("⚠️  Failed to read dir %s: %v\n", path, err)
			}

			return err
		}

		// Skip empty if requested
		if opts.SkipEmpty && len(entries) == 0 {
			if opts.Verbose {
				fmt.Printf("⏭️  Skipping %s (empty directory)\n", path)
			}

			return nil
		}

		// Handle leaf-depth scan mode
		if opts.LeafDepth > 0 {
			leafLevel, err := getLeafDepth(path)
			if err != nil {
				if opts.Verbose {
					fmt.Printf("⚠️  Failed to evaluate leaf depth of %s: %v\n", path, err)
				}
				return err
			}
			if leafLevel != opts.LeafDepth {
				if opts.Verbose {
					fmt.Printf("⏭️  Skipping %s (leafDepth=%d != expected=%d)\n", path, leafLevel, opts.LeafDepth)
				}
				return nil
			}

			// Special case: skip group-level node at exact depth-1
			if opts.LeafDepth == 1 && path == root {
				if opts.Verbose {
					fmt.Printf("⏭️  Skipping root %s for leafDepth=1\n", path)
				}

				return nil
			}
		} else if opts.OnlyLeaf {
			if containsDir(entries) {
				if opts.Verbose {
					fmt.Printf("⏭️  Skipping %s (not a leaf folder)\n", path)
				}
				return nil
			}

			// If LeafDepth == 1 + OnlyLeaf: only visit depth 1 and skip traversal
			if opts.MaxDepth == 1 && depth == 1 {
				fn(path, entries)
				return fs.SkipDir
			}
		}

		if opts.Verbose {
			fmt.Printf("📂 Visiting %s (%d entries)\n", path, len(entries))
		}

		fn(path, entries)
		return nil
	})
}

func WalkFiles(root string, opts WalkOptions, fn func(path string, size int64)) error {
	extMap := make(map[string]bool)
	for _, ext := range opts.Exts {
		extMap[strings.ToLower(ext)] = true
	}

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if opts.Verbose {
				fmt.Printf("⚠️  Error accessing %s: %v\n", path, err)
			}
			return err
		}

		if d.IsDir() {
			// rel, _ := filepath.Rel(root, path)
			// depth := strings.Count(rel, string(filepath.Separator))
			depth := pathDepth(root, path)
			if opts.MaxDepth >= 0 && depth > opts.MaxDepth {
				if opts.Verbose {
					fmt.Printf("⏭️  Skipping dir %s (too deep, depth %d)\n", path, depth)
				}

				return fs.SkipDir
			}

			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if len(extMap) > 0 && !extMap[ext] {
			if opts.Verbose {
				fmt.Printf("⏭️  Skipping file %s (unsupported ext: %s)\n", path, ext)
			}

			return nil
		}

		info, err := d.Info()
		if err != nil {
			if opts.Verbose {
				fmt.Printf("⚠️  Failed to get info for file %s: %v\n", path, err)
			}

			return nil
		}

		if opts.Verbose {
			fmt.Printf("📄 Found file %s (%d bytes)\n", path, info.Size())
		}

		fn(path, info.Size())
		return nil
	})
}
