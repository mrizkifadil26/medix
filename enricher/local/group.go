package local

import (
	"path/filepath"

	"github.com/mrizkifadil26/medix/utils/datawrapper"
	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type Group struct {
	Label string `json:"label"`
	Path  string `json:"path"`
}

type Collection struct {
	Name  string  `json:"name"`
	Group []Group `json:"group"`
}

type CollectionEnricher struct{}

func (s *CollectionEnricher) Name() string {
	return "collection"
}

func (s *CollectionEnricher) Enrich(
	data datawrapper.Data,
) error {
	root := data.Raw()

	itemsNode, ok := data.Get("items")
	if !ok {
		panic("items not found")
	}

	for _, idx := range itemsNode.Keys() {
		itemNode, _ := itemsNode.Get(idx)

		// check if "group" exists
		groupNode, ok := itemNode.Get("group")
		if !ok {
			continue
		}

		// group is expected to be []string
		groupRaw, ok := groupNode.Raw().([]any)
		if !ok || len(groupRaw) == 0 {
			continue
		}

		// build group objects
		var groups []Group
		for _, g := range groupRaw {
			label, _ := g.(string)
			groups = append(groups, Group{
				Label: label,
				Path:  filepath.Join(groupsPath(groups), label), // build nested path
			})
		}

		var collectionName string
		if len(groups) > 1 {
			// last becomes collection
			collectionName = groups[len(groups)-1].Label
			groups = groups[:len(groups)-1]
		}

		collection := Collection{
			Name:  collectionName,
			Group: groups,
		}

		jsonpath.Set(root, "items."+idx.(string)+".collection", collection)
	}

	return nil
}

// helper: build path for previous groups
func groupsPath(groups []Group) string {
	var parts []string
	for _, g := range groups {
		parts = append(parts, g.Label)
	}
	return filepath.Join(parts...)
}
