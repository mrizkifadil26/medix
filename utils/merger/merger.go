package merger

func Merge[T any](
	oldItems, newItems []T,
	opts MergeOptions[T],
) MergeResult[T] {
	oldMap := make(map[string]T)
	for _, item := range oldItems {
		oldMap[opts.KeyFn(item)] = item
	}

	newMap := make(map[string]T)
	for _, item := range newItems {
		newMap[opts.KeyFn(item)] = item
	}

	var added, removed, changed, merged []T

	for k, newItem := range newMap {
		if oldItem, exists := oldMap[k]; exists {
			if !opts.EqualFn(oldItem, newItem) {
				if opts.MergeFn != nil {
					merged = append(merged, opts.MergeFn(oldItem, newItem))
				} else {
					changed = append(changed, newItem)
				}
			}
		} else {
			added = append(added, newItem)
		}
	}

	for k, oldItem := range oldMap {
		if _, exists := newMap[k]; !exists {
			removed = append(removed, oldItem)
		}
	}

	return MergeResult[T]{
		Added:   added,
		Removed: removed,
		Changed: changed,
		Merged:  merged,
	}
}

// type Merger[T any, K comparable] struct {
// 	keySelector   func(T) K
// 	mergeStrategy func(oldItem, newItem T) T
// }

// func New[T any, K comparable]() *Merger[T, K] {
// 	return &Merger[T, K]{}
// }

// func (m *Merger[T, K]) WithKeySelector(f func(T) K) *Merger[T, K] {
// 	m.keySelector = f
// 	return m
// }

// func (m *Merger[T, K]) WithMergeStrategy(f func(oldItem, newItem T) T) *Merger[T, K] {
// 	m.mergeStrategy = f
// 	return m
// }

// func (m *Merger[T, K]) Merge(existing, incoming []T) []T {
// 	existingMap := make(map[K]T)
// 	for _, item := range existing {
// 		key := m.keySelector(item)
// 		existingMap[key] = item
// 	}

// 	for _, item := range incoming {
// 		key := m.keySelector(item)
// 		if oldItem, found := existingMap[key]; found {
// 			existingMap[key] = m.mergeStrategy(oldItem, item)
// 		} else {
// 			existingMap[key] = item
// 		}
// 	}

// 	result := make([]T, 0, len(existingMap))
// 	for _, item := range existingMap {
// 		result = append(result, item)
// 	}

// 	return result
// }
