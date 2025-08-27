package local

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mrizkifadil26/medix/utils/datawrapper"
	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type MediaSourceEnricher struct{}

type MediaSource struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Extension string `json:"ext"`
	Size      int64  `json:"size"`
}

func (s *MediaSourceEnricher) Name() string {
	return "media"
}

func (s *MediaSourceEnricher) Enrich(
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
			if strings.EqualFold(extStr, ".mkv") &&
				typ.Raw() == "file" &&
				!strings.Contains(strings.ToUpper(nameStr), "IMAX") {
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

		nameStr, ok := name.Raw().(string)
		if !ok {
			panic("name is not string")
		}

		pathStr, ok := path.Raw().(string)
		if !ok {
			panic("path is not string")
		}

		extStr, ok := ext.Raw().(string)
		if !ok {
			panic("ext is not string")
		}

		sizeStr, ok := size.Raw().(float64)
		if !ok {
			panic("size is not number")
		}

		fmt.Println("Selected", nameStr, pathStr, extStr, sizeStr)

		media := MediaSource{
			Name:      strings.TrimSuffix(nameStr, filepath.Ext(nameStr)),
			Path:      pathStr,
			Extension: extStr,
			Size:      int64(sizeStr),
		}
		fmt.Println("Media", media)

		jsonpath.Set(root, fmt.Sprintf("items.%d.media", idx), media)
	}

	return root, nil
}
