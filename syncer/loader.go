package syncer

import (
	"fmt"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/utils"
)

// LoadMedia loads raw media entries (movies or TV shows)
func LoadMedia(path string) ([]model.MediaEntry, error) {
	var entries model.MediaOutput
	err := utils.LoadJSON(path, &entries)
	if err != nil {
		return nil, fmt.Errorf("load media: %w", err)
	}

	return entries.Items, nil
}

// LoadIcons loads the raw icon map
func LoadIcons(path string) ([]model.IconEntry, error) {
	var icons model.IconIndex
	err := utils.LoadJSON(path, &icons)
	if err != nil {
		return nil, fmt.Errorf("load icons: %w", err)
	}

	return icons.Items, nil
}
