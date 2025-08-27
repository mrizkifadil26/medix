package local

import (
	"path/filepath"
	"strings"

	"github.com/mrizkifadil26/medix/utils/datawrapper"
	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type IconSourceEnricher struct{}

type IconSource struct {
	Name      string
	Path      string
	Extension string
	Size      int64
}

func (s *IconSourceEnricher) Name() string {
	return "icon"
}

func (s *IconSourceEnricher) Enrich(
	data datawrapper.Data,
	params map[string]string,
) (any, error) {
	root := data.Raw()

	itemsNode, ok := data.Get("items")
	if !ok {
		panic("items not found")
	}

	for _, idx := range itemsNode.Keys() {
		itemNode, _ := itemsNode.Get(idx)

		childrenNode, ok := itemNode.Get("children")
		if !ok {
			panic("children not found")
		}

		var selected datawrapper.Data
		for _, cidx := range childrenNode.Keys() {
			childNode, _ := childrenNode.Get(cidx)

			typ, _ := childNode.Get("type")
			ext, _ := childNode.Get("ext")
			name, _ := childNode.Get("name")

			extStr, _ := ext.Raw().(string)
			nameStr, _ := name.Raw().(string)
			if strings.EqualFold(extStr, ".ico") &&
				typ.Raw() == "file" &&
				nameStr != "" {
				selected = childNode
				break
			}
		}

		if selected == nil {
			// no valid .mkv found → skip
			continue
		}

		name, _ := selected.Get("name")
		path, _ := selected.Get("path")
		ext, _ := selected.Get("ext")
		size, _ := selected.Get("size")

		nameStr := name.Raw().(string)
		pathStr := path.Raw().(string)
		extStr := ext.Raw().(string)
		sizeStr := size.Raw().(int64)

		icon := IconSource{
			Name:      strings.TrimSuffix(nameStr, filepath.Ext(nameStr)),
			Path:      pathStr,
			Extension: extStr,
			Size:      sizeStr,
		}

		jsonpath.Set(root, "items."+idx.(string)+".icon", icon)
	}

	return root, nil
}
