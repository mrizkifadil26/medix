package iconmap

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GenerateIndex(cfg Config) (IconIndex, error) {
	start := time.Now()
	log.Println("🔍 Starting icon indexing...")

	dirMap := make(map[string][]*IconEntry)
	for _, src := range cfg.Sources {
		if err := collectIcons(src.Path, src.Name, cfg.ExcludeDirs, dirMap); err != nil {
			log.Printf("⚠️ Failed indexing %s: %v", src.Path, err)
		}
	}

	// grouped := flattenIcons(dirMap, cfg.ExcludeDirs)

	// var entries []IconEntry
	// totalItems := 0

	// for groupName, items := range grouped {
	// 	sort.Slice(items, func(i, j int) bool {
	// 		return items[i].Name < items[j].Name
	// 	})
	// 	for _, item := range items {
	// 		item.Group = groupName
	// 		entries = append(entries, item)
	// 		totalItems++
	// 	}
	// }

	flattened := flattenIcons(dirMap, cfg.ExcludeDirs)

	index := IconIndex{
		Type:           cfg.Type,
		Version:        "0.1.0",
		GeneratedAt:    time.Now(),
		TotalItems:     len(flattened),             // count of all icons
		TotalSources:   len(cfg.Sources),           // number of source directories
		GroupCount:     len(getRootGroups(dirMap)), // number of top-level groups
		ScanDurationMs: time.Since(start).Milliseconds(),
		Items:          flattened,
	}

	log.Printf("✅ Indexed %d icons in %v", len(flattened), time.Since(start))
	return index, nil
}

func flattenIcons(dirMap map[string][]*IconEntry, excludeDirs []string) []IconEntry {
	var flat []IconEntry

	for dir, icons := range dirMap {
		groupName := getSubGroupName(dir, excludeDirs)

		for _, icon := range icons {
			icon.Type = "icon"
			icon.Group = groupName
			flat = append(flat, *icon)
		}
	}

	return flat
}

func getSubGroupName(dir string, excludeDirs []string) string {
	parts := strings.SplitN(dir, string(os.PathSeparator), 2)
	if len(parts) < 2 {
		return ""
	}

	sub := filepath.Base(parts[1])
	if isExcluded(sub, excludeDirs) {
		return ""
	}

	return sub
}

func getRootGroups(dirMap map[string][]*IconEntry) map[string]struct{} {
	rootGroups := make(map[string]struct{})
	for dir := range dirMap {
		root := strings.SplitN(dir, string(os.PathSeparator), 2)[0]
		rootGroups[root] = struct{}{}
	}

	return rootGroups
}
