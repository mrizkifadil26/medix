package organizer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mrizkifadil26/medix/utils"
)

// Apply takes an OrganizeResult and executes file operations (e.g., move).
func Apply(result OrganizeResult) {
	for _, change := range result.Changes {
		src := change.Source
		dst := change.Target

		switch change.Action {
		case "copy":
			fmt.Printf("🟡 Copy: %s → %s\n", src, dst)

			/*
				err := copyFile(src, dst)
				if err != nil {
					fmt.Printf("❌ Failed to copy from %s → %s: %v\n", src, dst, err)
				} else {
					fmt.Printf("✅ Copied: %s → %s\n", src, dst)
				}
			*/

		case "move":
			fmt.Printf("🟡 Move: %s → %s\n", src, dst)

			/*
				err := moveFile(src, dst)
				if err != nil {
					fmt.Printf("❌ Failed to move from %s → %s: %v\n", src, dst, err)
				} else {
					fmt.Printf("✅ Moved: %s → %s\n", src, dst)
				}
			*/

		default:
			fmt.Printf("⚠️ Unknown action: %s for %s\n", change.Action, src)
		}
	}
}

// moveFile moves a file from src to dst, creating dst directories if needed.
func moveFile(src, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}
	return os.Rename(src, dst)
}

// copyFile copies a file from src to dst using a utility function.
func copyFile(src, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		return err
	}
	return utils.CopyFile(src, dst)
}
