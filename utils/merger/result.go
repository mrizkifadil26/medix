package merger

type MergerWithResult[T any, K comparable] struct {
	keySelector   func(T) K
	mergeStrategy func(oldItem, newItem T) (T, bool)
}

func NewWithResult[T any, K comparable]() *MergerWithResult[T, K] {
	return &MergerWithResult[T, K]{}
}

func (m *MergerWithResult[T, K]) WithKeySelector(f func(T) K) *MergerWithResult[T, K] {
	m.keySelector = f
	return m
}

func (m *MergerWithResult[T, K]) WithMergeStrategy(f func(oldItem, newItem T) (T, bool)) *MergerWithResult[T, K] {
	m.mergeStrategy = f
	return m
}

func (m *MergerWithResult[T, K]) Merge(existing, incoming []T) MergeResult[T] {
	existingMap := make(map[K]T)
	incomingMap := make(map[K]T)

	for _, item := range existing {
		existingMap[m.keySelector(item)] = item
	}

	for _, item := range incoming {
		incomingMap[m.keySelector(item)] = item
	}

	result := MergeResult[T]{
		Added:   []T{},
		Changed: []T{},
		Merged:  []T{},
		Removed: []T{},
	}

	// Process new and updated
	for key, newItem := range incomingMap {
		if oldItem, exists := existingMap[key]; exists {
			if mergedItem, changed := m.mergeStrategy(oldItem, newItem); changed {
				result.Changed = append(result.Changed, mergedItem)
				result.Merged = append(result.Merged, mergedItem)
			} else {
				result.Merged = append(result.Merged, oldItem)
			}

			delete(existingMap, key)
		} else {
			result.Added = append(result.Added, newItem)
			result.Merged = append(result.Merged, newItem)
		}
	}

	// Remaining in existingMap are removed
	for _, removed := range existingMap {
		result.Removed = append(result.Removed, removed)
	}

	return result
}
