package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

type WalkDirFunc func(dir string, entries []os.DirEntry)

type WalkFileFunc func(file string)

func walkDirs(
	root string,
	maxDepth int,
	cache *dirCache,
	fn WalkDirFunc,
) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}

		// Use cache if enabled
		// if cache != nil && cache.Has(path) {
		// 	return filepath.SkipDir
		// }

		depth := calcDepth(root, path)
		if depth > maxDepth {
			return filepath.SkipDir
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return nil // Ignore broken dirs
		}

		if isEmpty(entries) {
			return nil // ✅ Exclude completely empty folders
		}

		hasSubdir := false
		for _, entry := range entries {
			if entry.IsDir() {
				hasSubdir = true
				break
			}
		}

		if !hasSubdir && depth >= 1 {
			// This is a leaf folder — no subdirectories inside
			fn(path, entries)
			return filepath.SkipDir // Don’t go deeper
		}

		// Cache current path
		// if cache != nil {
		// 	cache.Set(path)
		// }

		return nil
	})
}

func walkFiles(
	root string,
	maxDepth int,
	exts []string,
	cache *dirCache,
	fn WalkFileFunc,
) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		if calcDepth(root, path) > maxDepth {
			return nil
		}

		if len(exts) == 0 || hasValidExt(path, exts) {
			fn(path)
		}

		return nil
	})
}

func calcDepth(root, path string) int {
	rel, _ := filepath.Rel(root, path)
	if rel == "." {
		return 0
	}
	return len(splitPath(rel))
}

func splitPath(p string) []string {
	if p == "" || p == "." {
		return nil
	}

	return strings.Split(filepath.ToSlash(p), "/")
}

// func isLeaf(entries []os.DirEntry) bool {
// 	for _, e := range entries {
// 		if e.IsDir() {
// 			return false
// 		}
// 	}

// 	return true
// }

func buildGroupLabel(root, path string) []string {
	rel, err := filepath.Rel(root, path)
	if err != nil || rel == "." || rel == "" {
		return nil
	}

	parts := splitPath(filepath.ToSlash(rel))
	return parts // include everything, even just one folder
}

func isEmpty(entries []os.DirEntry) bool {
	return len(entries) == 0
}

func CountTargetDirs(
	root, label string,
	maxDepth int,
	cache *dirCache,
	filterFn func(path string, entries []os.DirEntry) bool,
) int {
	count := 0
	_ = walkDirs(root, maxDepth, cache, func(path string, entries []os.DirEntry) {
		if filterFn == nil || filterFn(path, entries) {
			count++
		}
	})

	return count
}
