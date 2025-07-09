package indexer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/util"
)

func BuildIconIndex(cfg IconIndexerConfig) (model.IconIndex, error) {
	log.Println("🔍 Starting icon indexing...")

	dirMap := make(map[string][]model.IconEntry)
	for _, src := range cfg.Sources {
		if err := collectIcons(src.Path, src.Source, cfg.ExcludeDirs, dirMap); err != nil {
			log.Printf("⚠️ Failed indexing %s: %v", src.Path, err)
		}
	}
	// return index, nil
	index := model.IconIndex{
		Type:        "genre",
		GeneratedAt: time.Now(),
	}

	tree := groupIcons(dirMap, cfg.ExcludeDirs)

	var roots []string
	for dir := range tree {
		roots = append(roots, dir)
	}
	sort.Strings(roots)

	for _, root := range roots {
		entries := tree[root]
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name < entries[j].Name
		})
		index.Data = append(index.Data, model.IconGroup{
			Name:  root,
			Items: entries,
		})
	}

	total := 0
	for _, g := range index.Data {
		total += len(g.Items)
	}
	fmt.Printf("✅ Indexed %d icons\n", total)

	return index, nil
}

func collectIcons(baseDir, source string, excludeDirs []string, dirMap map[string][]model.IconEntry) error {
	return filepath.WalkDir(baseDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("⚠️ Walk error in %s: %v", path, err)
			return nil
		}

		if d.IsDir() && isExcluded(filepath.Base(path), excludeDirs) {
			return filepath.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(strings.ToLower(d.Name()), ".ico") {
			return nil
		}

		info, err := os.Stat(path)
		if err != nil {
			log.Printf("⚠️ Failed to stat %s: %v", path, err)
			return nil
		}

		relDir, err := filepath.Rel(baseDir, filepath.Dir(path))
		if err != nil || relDir == "." {
			return nil
		}

		dirMap[relDir] = append(dirMap[relDir], model.IconEntry{
			ID:       util.Slugify(d.Name()),
			Name:     d.Name(),
			Size:     info.Size(),
			Source:   source,
			FullPath: path,
			Type:     "icon",
		})

		return nil
	})
}

func groupIcons(dirMap map[string][]model.IconEntry, excludeDirs []string) map[string][]model.IconEntry {
	tree := make(map[string][]model.IconEntry)

	for dir, icons := range dirMap {
		parts := strings.SplitN(dir, string(os.PathSeparator), 2)
		root := parts[0]
		sub := ""
		if len(parts) > 1 {
			sub = parts[1]
		}

		if sub != "" && !isExcluded(filepath.Base(sub), excludeDirs) {
			collectionName := filepath.Base(sub)
			for i := range icons {
				icons[i].Type = "icon"
			}

			tree[root] = append(tree[root], model.IconEntry{
				Name:  collectionName,
				Type:  "collection",
				Items: icons,
			})
		} else {
			for _, icon := range icons {
				tree[root] = append(tree[root], icon)
			}
		}
	}

	return tree
}

func isExcluded(name string, excludeDirs []string) bool {
	for _, ex := range excludeDirs {
		if name == ex {
			return true
		}
	}

	return false
}
