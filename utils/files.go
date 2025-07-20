package utils

import (
	"fmt"
	"io"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func EnsureDir(path string) error {
	if path == "" {
		return fmt.Errorf("directory path is empty")
	}
	return os.MkdirAll(path, os.ModePerm)
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Sync()
}
