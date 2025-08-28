package local

import (
	"strings"

	"github.com/mrizkifadil26/medix/utils/datawrapper"
	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type MediaFilter struct{}

func (f MediaFilter) Name() string { return "media" }

func (f MediaFilter) Apply(
	item datawrapper.Data,
	errs *[]error,
) {
	children, ok := item.Get("children")
	if !ok {
		return
	}

	var mainMedia any
	var imaxList []any

	for _, idx := range children.Keys() {
		childNode, _ := children.Get(idx)

		ext := safeString(childNode, "ext", errs)
		name := safeString(childNode, "name", errs)

		if strings.EqualFold(ext, ".mkv") {
			if strings.Contains(strings.ToUpper(name), "IMAX") {
				imaxList = append(imaxList, childNode.Raw())
			} else if mainMedia == nil {
				mainMedia = childNode.Raw()
			}
		}

		if mainMedia != nil {
			_ = jsonpath.Set(item, "media", mainMedia)
		}

		if len(imaxList) > 0 {
			_ = jsonpath.Set(item, "extras.imax", imaxList)
		}
	}
}
