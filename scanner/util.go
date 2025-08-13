package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

func containsDir(entries []os.DirEntry) bool {
	for _, e := range entries {
		if e.IsDir() {
			return true
		}
	}
	return false
}

func getDepth(base, path string) int {
	rel, _ := filepath.Rel(base, path)
	if rel == "." {
		return 0
	}

	return len(strings.Split(rel, string(filepath.Separator)))
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}

	return false
}

func getLeafDepth(start string) (int, error) {
	level := 0
	current := start

	for {
		entries, err := os.ReadDir(current)
		if err != nil {
			return 0, err
		}
		found := false
		for _, e := range entries {
			if e.IsDir() {
				current = filepath.Join(current, e.Name())
				found = true
				break
			}
		}
		if !found {
			break
		}
		level++
	}

	return level, nil
}

// isHidden checks if a file/dir is hidden (starts with .)
func isHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

func isLeafDir(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	for {
		infos, err := f.Readdir(1) // read 1 entry at a time
		if err != nil {
			break // EOF or error
		}

		if infos[0].IsDir() {
			return false // has a subdirectory
		}
	}

	return true // no subdirectories found
}

func determineLogLevel(opts ScanOptions) string {
	if opts.Trace {
		return "TRACE"
	}

	if opts.Debug { // if you have a Debug bool in opts
		return "DEBUG"
	}

	if opts.Verbose {
		return "INFO" // or VERBOSE as you want
	}

	return "ERROR"
}
