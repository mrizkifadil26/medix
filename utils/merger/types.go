package merger

type MergeOptions[T any] struct {
	KeyFn   func(item T) string
	EqualFn func(a, b T) bool
	MergeFn func(oldItem, newItem T) T // optional
}

type MergeResult[T any] struct {
	Added   []T
	Removed []T
	Changed []T
	Merged  []T
}
