package local

import (
	"strings"

	"github.com/mrizkifadil26/medix/utils/datawrapper"
	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type IconFilter struct{}

func (f IconFilter) Name() string {
	return "icon"
}

func (f IconFilter) Apply(
	item datawrapper.Data,
	errs *[]error,
) {
	children, ok := item.Get("children")
	if !ok {
		return
	}

	for _, idx := range children.Keys() {
		childNode, _ := children.Get(idx)

		ext := safeString(childNode, "ext", errs)
		if ext == ".mkv" {
			_ = jsonpath.Set(item, "media", childNode.Raw())
		}

		name := safeString(childNode, "name", errs)
		if strings.Contains(strings.ToLower(name), "imax") {
			_ = jsonpath.Set(item, "extras.imax", childNode.Raw())
		}
	}
}
