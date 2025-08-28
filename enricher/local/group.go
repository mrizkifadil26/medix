package local

import (
	"path/filepath"

	"github.com/mrizkifadil26/medix/utils/datawrapper"
	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type CollectionFilter struct{}

func (f CollectionFilter) Name() string {
	return "collection"
}

func (f CollectionFilter) Apply(
	data datawrapper.Data,
	errs *[]error,
) {
	groupNode, ok := data.Get("group")
	if !ok {
		return
	}

	groupRaw, ok := groupNode.Raw().([]any)
	if !ok || len(groupRaw) == 0 {
		return
	}

	// build group objects
	var groups []Group
	for _, g := range groupRaw {
		label, _ := g.(string)
		groups = append(groups, Group{
			Label: label,
			Path:  filepath.Join(groupsPath(groups), label),
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

	// attach to child
	_ = jsonpath.Set(data, "collection", collection)
}

// helper: build path for previous groups
func groupsPath(groups []Group) string {
	var parts []string
	for _, g := range groups {
		parts = append(parts, g.Label)
	}
	return filepath.Join(parts...)
}
