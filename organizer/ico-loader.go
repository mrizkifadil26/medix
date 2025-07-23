package organizer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/utils"
)

// Regex to match (alt), (alt 2), (alt 3), ...
var altSuffixRegex = regexp.MustCompile(`(?i)\s*\(alt(?: \d+)?\)$`)

func loadIcons(sources []IconSource) map[string][]model.IconRef {
	iconMap := make(map[string][]model.IconRef)

	for _, src := range sources {
		entries, err := os.ReadDir(src.Path)
		if err != nil {
			fmt.Printf("⚠️ Failed to read icon source %s (%s): %v\n", src.Source, src.Path, err)
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".ico") {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				continue
			}

			fullName := entry.Name()
			baseName := stripAltSuffix(fullName) // Strip "(alt)", "(alt 2)", etc.
			slug := utils.Slugify(baseName)

			meta := model.IconRef{
				Name:     fullName,
				FullPath: filepath.Join(src.Path, fullName),
				Size:     info.Size(),
			}

			iconMap[slug] = append(iconMap[slug], meta)
		}
	}

	return iconMap
}

// stripAltSuffix removes any trailing (alt), (alt 2), etc. from a file name (without extension)
func stripAltSuffix(name string) string {
	name = strings.TrimSuffix(name, filepath.Ext(name))
	return strings.TrimSpace(altSuffixRegex.ReplaceAllString(name, ""))
}
