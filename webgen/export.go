package webgen

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func shorten(path string) string {
	return filepath.Clean(path)
}

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func EnsureDirs(paths ...string) {
	for _, path := range paths {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func CopyFile(src, dst string) error {
	fmt.Printf("[FILE] %s → %s\n", shorten(src), shorten(dst))
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}

func CopyDir(src, dst string) error {
	fmt.Printf("[DIR ] %s → %s\n", shorten(src), shorten(dst))
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Compute relative path from src → path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		return CopyFile(path, targetPath)
	})
}
