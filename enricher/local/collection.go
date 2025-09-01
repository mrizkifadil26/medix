package local

import (
	"fmt"

	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type CollectionFilter struct{}

func (f CollectionFilter) Name() string {
	return "collection"
}

func (f CollectionFilter) Apply(
	data any,
	errs *[]error,
) {
	groupNode, err := jsonpath.Get(data, "group_label")
	if err != nil {
		*errs = append(*errs, fmt.Errorf("directory has no group label"))
		return
	}

	groupRaw, ok := groupNode.([]any)
	if !ok || len(groupRaw) == 0 {
		return
	}

	// enforce length rules
	if len(groupRaw) > 2 {
		*errs = append(*errs, fmt.Errorf("invalid group_label length: %d (must be 1 or 2)", len(groupRaw)))
		return
	}

	// build group objects
	if len(groupRaw) == 2 {
		// build groups (only the first is group)
		var groups []Group
		firstLabel, _ := groupRaw[0].(string)
		groups = append(groups, Group{
			Label: firstLabel,
			Path:  firstLabel,
		})

		collectionName, _ := groupRaw[1].(string)

		collection := Collection{
			Name:  collectionName,
			Group: groups,
		}

		// attach to child
		_ = jsonpath.Set(data, "collection", collection)
	}
}
