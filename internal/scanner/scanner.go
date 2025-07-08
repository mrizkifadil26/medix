package scanner

import (
	"os"
	"path/filepath"
	"time"

	"github.com/schollz/progressbar/v3"

	"github.com/mrizkifadil26/medix/model"
)

func ScanDirectory(contentType, root string) model.RawOutput {
	result := model.RawOutput{
		Type:        contentType,
		GeneratedAt: time.Now(),
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return result
	}

	var genres []os.DirEntry
	for _, g := range entries {
		if g.IsDir() {
			genres = append(genres, g)
		}
	}

	bar := progressbar.NewOptions(len(genres),
		progressbar.OptionSetDescription("Scanning genres"),
		progressbar.OptionSetWidth(30),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "#", SaucerPadding: " ", BarStart: "[", BarEnd: "]"}),
	)

	for _, genre := range genres {
		genrePath := filepath.Join(root, genre.Name())
		items := scanGenre(genrePath)
		if len(items) > 0 {
			result.Data = append(result.Data, model.GenreBlock{
				Genre: genre.Name(),
				Items: items,
			})
		}

		bar.Add(1)
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
