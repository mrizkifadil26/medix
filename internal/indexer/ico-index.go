package indexer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mrizkifadil26/medix/util"
)

type IconIndex struct {
	Type        string      `json:"type"` // "genre"
	GeneratedAt time.Time   `json:"generated_at"`
	Groups      []IconGroup `json:"groups"`
}

type IconGroup struct {
	ID    string      `json:"id,omitempty"` // e.g. "sci-fi"
	Name  string      `json:"name"`         // genre name like "Sci-Fi"
	Items []IconEntry `json:"items"`
}

type IconEntry struct {
	ID        string      `json:"id,omitempty"`
	Name      string      `json:"name"`
	Extension string      `json:"extension,omitempty"`
	Size      int64       `json:"size,omitempty"`
	Source    string      `json:"source,omitempty"`
	FullPath  string      `json:"full_path,omitempty"`
	Type      string      `json:"type"` // "icon" or "collection"
	Items     []IconEntry `json:"items,omitempty"`
}

const (
	PersonalDir   = "/mnt/c/Users/Rizki/OneDrive/Pictures/Icons/Personal Icon Pack/Movies/ICO"
	DownloadedDir = "/mnt/c/Users/Rizki/OneDrive/Pictures/Icons/Downloaded Icon Pack/Movie Icon Pack/downloaded"
	OutputPath    = "data/ico.index.json"
	ExcludeDir    = "Collection"
)

func BuildIconIndex() (IconIndex, error) {
	// Index personal icons first (higher priority)
	// if err := scanIcoDir(PersonalDir, "personal", index, true); err != nil {
	// 	return nil, err
	// }

	// // Index downloaded icons (lower priority, don't overwrite)
	// if err := scanIcoDir(DownloadedDir, "downloaded", index, false); err != nil {
	// 	return nil, err
	// }

	// return index, nil
	index := IconIndex{
		Type:        "genre",
		GeneratedAt: time.Now(),
	}

	fmt.Println("🔍 Indexing icons...")
	dirMap := make(map[string][]IconEntry)
	collectIcons(PersonalDir, "personal", dirMap)
	collectIcons(DownloadedDir, "downloaded", dirMap)

	// Convert map to sorted slice
	tree := make(map[string][]IconEntry)
	for dir, icons := range dirMap {
		parts := strings.SplitN(dir, string(os.PathSeparator), 2)
		root := parts[0]
		sub := ""
		if len(parts) > 1 {
			sub = parts[1]
		}

		if sub != "" && filepath.Base(sub) != ExcludeDir {
			// treat folder as collection name only
			collectionName := filepath.Base(sub)
			for i := range icons {
				icons[i].Type = "icon"
			}
			collectionEntry := IconEntry{
				Name:  collectionName,
				Type:  "collection",
				Items: icons,
			}
			tree[root] = append(tree[root], collectionEntry)
		} else {
			for _, icon := range icons {
				icon.Type = "icon"
				tree[root] = append(tree[root], icon)
			}
		}
	}

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
		index.Groups = append(index.Groups, IconGroup{
			ID:    util.Slugify(root),
			Name:  root,
			Items: entries,
		})
	}

	total := 0
	for _, g := range index.Groups {
		total += len(g.Items)
	}
	fmt.Printf("✅ Indexed %d icons\n", total)

	// index.Sources = append(index.Sources, personalIcons...)
	// index.Sources = append(index.Sources, downloadedIcons...)

	return index, nil
}

// func scanIcoDir(dir string, source string, index IconIndex, overwrite bool) error {
// 	entries, err := os.ReadDir(dir)
// 	if err != nil {
// 		return fmt.Errorf("read dir error (%s): %w", dir, err)
// 	}

// 	for _, entry := range entries {
// 		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".ico") {
// 			continue
// 		}

// 		name := strings.TrimSuffix(entry.Name(), ".ico")
// 		if name == "" {
// 			continue
// 		}

// 		if _, exists := index[name]; exists && !overwrite {
// 			continue // skip if already exists and overwrite is false
// 		}

// 		index[name] = IconEntry{
// 			Source: source,
// 			Path:   filepath.Join(dir, entry.Name()),
// 		}
// 	}

// 	return nil
// }

func collectIcons(baseDir, source string, dirMap map[string][]IconEntry) {
	err := filepath.WalkDir(baseDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("⚠️ Walk error in %s: %v", path, err)
			return nil
		}

		if d.IsDir() {
			if filepath.Base(path) == ExcludeDir {
				return filepath.SkipDir
			}
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
		if err != nil {
			relDir = "UNKNOWN"
		}

		// Ignore top-level icons (those directly under baseDir)
		if relDir == "." {
			return nil
		}

		dirMap[relDir] = append(dirMap[relDir], IconEntry{
			ID:        util.Slugify(d.Name()),
			Name:      d.Name(),
			Size:      info.Size(),
			Source:    source,
			FullPath:  path,
			Type:      "icon",
			Extension: filepath.Ext(d.Name()),
		})

		return nil
	})

	if err != nil {
		log.Printf("⚠️ Failed walking %s: %v", baseDir, err)
	}
}

func SaveIconIndex(index IconIndex) error {
	_ = os.MkdirAll(filepath.Dir(OutputPath), 0755)
	file, err := os.Create(OutputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", " ")
	return enc.Encode(index)
}
