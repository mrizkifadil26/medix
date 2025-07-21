package iconmap

import (
	"log"
	"time"

	"github.com/mrizkifadil26/medix/model"
)

func GenerateIndex(cfg Config) (*model.IconIndex, error) {
	start := time.Now()
	log.Println("🔍 Starting icon indexing...")

	flatEntries := []model.IconEntry{}
	groupCount := 0

	for _, src := range cfg.Sources {
		grouped, err := ScanIconDirectory(src.Path, src.Name, cfg.ExcludeDirs)
		if err != nil {
			log.Printf("⚠️ Failed indexing %s: %v", src.Path, err)
			return nil, err
		}

		flat := FlattenGroupedIcons(grouped, cfg.ExcludeDirs)
		flatEntries = append(flatEntries, flat...)

		groupCount += len(grouped)
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

	// flattened := flattenIcons(dirMap, cfg.ExcludeDirs)

	index := &model.IconIndex{
		Type:           cfg.Type,
		Version:        "0.1.0",
		GeneratedAt:    time.Now(),
		TotalItems:     len(flatEntries), // count of all icons
		TotalSources:   len(cfg.Sources), // number of source directories
		GroupCount:     groupCount,       // number of top-level groups
		ScanDurationMs: time.Since(start).Milliseconds(),
		Items:          flatEntries,
	}

	log.Printf("✅ Indexed %d icons in %v", len(flatEntries), time.Since(start))
	return index, nil
}

// func getRootGroups(dirMap map[string][]*model.IconEntry) map[string]struct{} {
// 	rootGroups := make(map[string]struct{})
// 	for dir := range dirMap {
// 		root := strings.SplitN(dir, string(os.PathSeparator), 2)[0]
// 		rootGroups[root] = struct{}{}
// 	}

// 	return rootGroups
// }
