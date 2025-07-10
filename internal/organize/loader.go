package organize

import (
	"encoding/json"
	"os"

	"github.com/mrizkifadil26/medix/model"
)

func LoadRawMetadata(path string) ([]model.MovieEntry, error) {
	var root model.MovieOutput
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	var entries []model.MovieEntry
	for _, group := range root.Data {
		for _, item := range group.Items {
			item.Group = group.Name
			entries = append(entries, item)
		}
	}

	return entries, nil
}

func LoadScatteredIcons(sources []SourceDir, excludeDirs []string) ([]model.IconEntry, error) {
	var all []model.IconEntry
	for _, src := range sources {
		icons, err := IndexScatteredIcons(src.Path, src.Source, excludeDirs)
		if err != nil {
			return nil, err
		}
		all = append(all, icons...)
	}
	return all, nil
}
