package local

import (
	"fmt"
	"strings"

	"github.com/mrizkifadil26/medix/utils/datawrapper"
	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type Subtitle struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Ext  string `json:"ext"`
}

// The field "subtitles" is a map of language code -> Subtitle
type Subtitles map[string]Subtitle

type SubtitleEnricher struct{}

func (s *SubtitleEnricher) Name() string {
	return "subtitle"
}

func (s *SubtitleEnricher) Enrich(
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

		subs := make(map[string]Subtitle)

		for _, cidx := range childrenNode.Keys() {
			childNode, _ := childrenNode.Get(cidx)

			name, _ := childNode.Get("name")
			ext, _ := childNode.Get("ext")
			typ, _ := childNode.Get("type")
			path, _ := childNode.Get("path")

			extStr, _ := ext.Raw().(string)
			nameStr, _ := name.Raw().(string)
			pathStr, _ := path.Raw().(string)

			if typ.Raw() == "file" && strings.EqualFold(extStr, ".srt") {
				// try detect lang by suffix
				lang := "id"                                // fallback
				base := strings.TrimSuffix(nameStr, extStr) // remove ".srt"
				if i := strings.LastIndex(base, "."); i != -1 {
					candidate := base[i+1:]
					if len(candidate) > 0 {
						lang = strings.ToLower(candidate)
					}
				}

				subs[lang] = Subtitle{
					Name: nameStr,
					Path: pathStr,
					Ext:  extStr,
				}
			}
		}

		if len(subs) > 0 {
			jsonpath.Set(root, fmt.Sprintf("items.%d.subtitles", idx), subs)
		}
	}

	return root, nil
}
