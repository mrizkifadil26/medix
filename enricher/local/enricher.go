package local

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type LocalEnricher struct{}

func (s *LocalEnricher) Name() string {
	return "local"
}

var filters = []Filter{
	MediaFilter{},
	SubtitlesFilter{},
	IconFilter{},
	CollectionFilter{},
}

func (e *LocalEnricher) Enrich(
	data any,
	options map[string]string,
) (any, error) {
	var errs []error

	// Get optional filter argument from options
	filterArg := ""
	if options != nil {
		filterArg = options["filters"]
	}

	// Build allowed map
	allowed := map[string]bool{}
	if filterArg == "" {
		for _, f := range filters {
			allowed[f.Name()] = true
		}
	} else {
		for _, f := range filters {
			allowed[f.Name()] = false
		}
		for _, name := range strings.Split(filterArg, ",") {
			allowed[strings.TrimSpace(name)] = true
		}
	}

	// Get items using your Get helper
	nodes, err := jsonpath.Get(data, "items.#")
	if err != nil {
		// items key not found, nothing to enrich
		return data, nil
	}

	items, ok := nodes.([]any)
	if !ok {
		return data, fmt.Errorf("items is not an array, got %T", nodes)
	}

	for _, item := range items {
		for _, f := range filters {
			if allowed[f.Name()] {
				f.Apply(item, &errs)
			}
		}
	}

	if len(errs) > 0 {
		return data, errors.Join(errs...)
	}

	return data, nil
}

// func (e *LocalEnricher) Enrich(
// 	data datawrapper.Data,
// 	params map[string]string,
// ) (any, error) {
// 	root := data.Raw()

// 	itemsNode, ok := data.Get("items")
// 	if !ok {
// 		panic("items not found")
// 	}

// 	// --- parse filter param ---
// 	allowed := map[string]bool{
// 		"media":      true,
// 		"subtitles":  true,
// 		"icon":       true,
// 		"collection": true,
// 	}

// 	if f, ok := params["filter"]; ok && f != "" {
// 		// reset and allow only selected
// 		for k := range allowed {
// 			allowed[k] = false
// 		}

// 		for _, part := range strings.Split(f, ",") {
// 			allowed[strings.TrimSpace(part)] = true
// 		}
// 	}

// 	for _, idx := range itemsNode.Keys() {
// 		itemNode, _ := itemsNode.Get(idx)

// 		childrenNode, ok := itemNode.Get("children")
// 		if !ok {
// 			continue
// 		}

// 		// ---------- Media & IMAX ----------
// 		if allowed["media"] {
// 			var mediaSelected datawrapper.Data
// 			var imaxList []MediaSource

// 			for _, cidx := range childrenNode.Keys() {
// 				childNode, _ := childrenNode.Get(cidx)

// 				typ, _ := childNode.Get("type")
// 				ext, _ := childNode.Get("ext")
// 				name, _ := childNode.Get("name")
// 				path, _ := childNode.Get("path")
// 				size, _ := childNode.Get("size")

// 				extStr, _ := ext.Raw().(string)
// 				nameStr, _ := name.Raw().(string)
// 				pathStr, _ := path.Raw().(string)
// 				sizeFloat := size.Raw().(float64)
// 				if strings.EqualFold(extStr, ".mkv") && typ.Raw() == "file" {
// 					ms := MediaSource{
// 						Name:      strings.TrimSuffix(nameStr, filepath.Ext(nameStr)),
// 						Path:      pathStr,
// 						Extension: extStr,
// 						Size:      int64(sizeFloat),
// 					}

// 					if strings.Contains(strings.ToUpper(nameStr), "IMAX") {
// 						imaxList = append(imaxList, ms)
// 					} else if mediaSelected == nil {
// 						mediaSelected = childNode
// 						// store as main media
// 						jsonpath.Set(root, fmt.Sprintf("items.%s.media", idx), ms)
// 					}
// 				}
// 			}

// 			if len(imaxList) > 0 {
// 				jsonpath.Set(root, fmt.Sprintf("items.%s.imax", idx), imaxList)
// 			}
// 		}
// 		// ---------- Subtitles ----------
// 		if allowed["subtitles"] {
// 			subs := make(map[string]Subtitle)
// 			for _, cidx := range childrenNode.Keys() {
// 				childNode, _ := childrenNode.Get(cidx)

// 				typ, _ := childNode.Get("type")
// 				ext, _ := childNode.Get("ext")
// 				name, _ := childNode.Get("name")
// 				path, _ := childNode.Get("path")

// 				extStr, _ := ext.Raw().(string)
// 				nameStr, _ := name.Raw().(string)
// 				pathStr, _ := path.Raw().(string)

// 				if typ.Raw() == "file" && strings.EqualFold(extStr, ".srt") {
// 					lang := "id" // fallback
// 					base := strings.TrimSuffix(nameStr, extStr)
// 					if i := strings.LastIndex(base, "."); i != -1 {
// 						candidate := base[i+1:]
// 						if len(candidate) > 0 {
// 							lang = strings.ToLower(candidate)
// 						}
// 					}
// 					subs[lang] = Subtitle{Name: nameStr, Path: pathStr, Ext: extStr}
// 				}
// 			}
// 			if len(subs) > 0 {
// 				jsonpath.Set(root, fmt.Sprintf("items.%s.subtitles", idx), subs)
// 			}
// 		}

// 		// ---------- Icon ----------
// 		if allowed["icon"] {
// 			var iconSelected datawrapper.Data
// 			for _, cidx := range childrenNode.Keys() {
// 				childNode, _ := childrenNode.Get(cidx)

// 				typ, _ := childNode.Get("type")
// 				ext, _ := childNode.Get("ext")
// 				name, _ := childNode.Get("name")

// 				extStr, _ := ext.Raw().(string)
// 				nameStr, _ := name.Raw().(string)
// 				if strings.EqualFold(extStr, ".ico") &&
// 					typ.Raw() == "file" &&
// 					nameStr != "" {
// 					iconSelected = childNode
// 					break
// 				}
// 			}
// 			if iconSelected != nil {
// 				name, _ := iconSelected.Get("name")
// 				path, _ := iconSelected.Get("path")
// 				ext, _ := iconSelected.Get("ext")
// 				size, _ := iconSelected.Get("size")

// 				icon := IconSource{
// 					Name:      strings.TrimSuffix(name.Raw().(string), filepath.Ext(name.Raw().(string))),
// 					Path:      path.Raw().(string),
// 					Extension: ext.Raw().(string),
// 					Size:      int64(size.Raw().(float64)),
// 				}
// 				jsonpath.Set(root, fmt.Sprintf("items.%s.icon", idx), icon)
// 			}
// 		}

// 		// ---------- Collection ----------
// 		if allowed["collection"] {
// 			groupNode, ok := itemNode.Get("group")
// 			if ok {
// 				groupRaw, ok := groupNode.Raw().([]any)
// 				if ok && len(groupRaw) > 0 {
// 					var groups []Group
// 					for _, g := range groupRaw {
// 						label, _ := g.(string)
// 						groups = append(groups, Group{
// 							Label: label,
// 							Path:  filepath.Join(groupsPath(groups), label),
// 						})
// 					}

// 					var collectionName string
// 					if len(groups) > 1 {
// 						collectionName = groups[len(groups)-1].Label
// 						groups = groups[:len(groups)-1]
// 					}

// 					collection := Collection{
// 						Name:  collectionName,
// 						Group: groups,
// 					}
// 					jsonpath.Set(root, fmt.Sprintf("items.%s.collection", idx), collection)
// 				}
// 			}
// 		}
// 	}

// 	return root, nil
// }
