package local

import (
	"github.com/mrizkifadil26/medix/utils/datawrapper"
	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type SubtitlesFilter struct{}

func (f SubtitlesFilter) Name() string {
	return "subtitle"
}

func (f SubtitlesFilter) Apply(
	item datawrapper.Data,
	errs *[]error,
) {
	children, ok := item.Get("children")
	if !ok {
		return
	}

	subs := map[string]any{}
	for _, idx := range children.Keys() {
		childNode, _ := children.Get(idx)

		ext := safeString(childNode, "ext", errs)
		if ext == ".srt" {
			lang := safeString(childNode, "lang", errs)
			subs[lang] = childNode.Raw()
		}
	}

	if len(subs) > 0 {
		_ = jsonpath.Set(item, "subtitles", subs)
	}
}
