package scanner

import (
	"os"
	"path/filepath"
	"time"

	"mrizkifadil26.github.io/medix/model"
)

func ScanDirectory(contentType, root string) model.RawOutput {
	result := model.RawOutput{
		Type:        contentType,
		GeneratedAt: time.Now(),
	}

	genres, err := os.ReadDir(root)
	if err != nil {
		return result
	}

	for _, genre := range genres {
		if !genre.IsDir() {
			continue
		}
		genrePath := filepath.Join(root, genre.Name())
		items := scanGenre(genrePath)
		if len(items) > 0 {
			result.Data = append(result.Data, model.GenreBlock{
				Genre: genre.Name(),
				Items: items,
			})
		}
	}

	return result
}

func scanGenre(genrePath string) []model.RawItem {
	var items []model.RawItem

	entries, err := os.ReadDir(genrePath)
	if err != nil {
		return items
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		entryPath := filepath.Join(genrePath, entry.Name())
		subDirs, _ := os.ReadDir(entryPath)

		if isCollection(subDirs) {
			var children []model.RawChild
			for _, sub := range subDirs {
				if !sub.IsDir() {
					continue
				}
				subPath := filepath.Join(entryPath, sub.Name())
				children = append(children, model.RawChild{
					Name:   sub.Name(),
					Path:   subPath,
					Status: ResolveStatus(subPath),
				})
			}

			items = append(items, model.RawItem{
				Type:     "collection",
				Name:     entry.Name(),
				Path:     entryPath,
				Status:   ResolveStatus(entryPath),
				Children: children,
			})
		} else {
			items = append(items, model.RawItem{
				Type:   "single",
				Name:   entry.Name(),
				Path:   entryPath,
				Status: ResolveStatus(entryPath),
			})
		}
	}

	return items
}
