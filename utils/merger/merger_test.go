package merger

import (
	"reflect"
	"testing"
)

type Entity struct {
	ID    string
	Name  string
	Value int
}

func TestBasicMerge(t *testing.T) {
	old := []Entity{
		{ID: "1", Name: "Old One", Value: 10},
		{ID: "2", Name: "Old Two", Value: 20},
	}

	new := []Entity{
		{ID: "2", Name: "New Two", Value: 30},   // updated
		{ID: "3", Name: "New Three", Value: 40}, // added
	}

	merger := NewWithResult[Entity, string]().
		WithKeySelector(func(e Entity) string { return e.ID }).
		WithMergeStrategy(func(old, new Entity) (Entity, bool) {
			if old.Name != new.Name || old.Value != new.Value {
				return new, true
			}

			return old, false
		})

	result := merger.Merge(old, new)

	assertEqual(t, []Entity{{ID: "3", Name: "New Three", Value: 40}}, result.Added, "added")
	assertEqual(t, []Entity{{ID: "2", Name: "New Two", Value: 30}}, result.Changed, "changed")
	assertEqual(t, []Entity{{ID: "1", Name: "Old One", Value: 10}}, result.Removed, "removed")
}

func TestNoChanges(t *testing.T) {
	data := []Entity{
		{ID: "1", Name: "Same One", Value: 100},
	}

	merger := NewWithResult[Entity, string]().
		WithKeySelector(func(e Entity) string { return e.ID }).
		WithMergeStrategy(func(old, new Entity) (Entity, bool) {
			return old, false
		})

	result := merger.Merge(data, data)

	assertEqual(t, []Entity{}, result.Added, "added")
	assertEqual(t, []Entity{}, result.Changed, "changed")
	assertEqual(t, []Entity{}, result.Removed, "removed")
	assertEqual(t, data, result.Merged, "merged")
}

func TestOnlyAdd(t *testing.T) {
	old := []Entity{}
	new := []Entity{{ID: "1", Name: "First", Value: 10}}

	merger := NewWithResult[Entity, string]().
		WithKeySelector(func(e Entity) string { return e.ID }).
		WithMergeStrategy(func(old, new Entity) (Entity, bool) {
			return new, true
		})

	result := merger.Merge(old, new)

	assertEqual(t, new, result.Added, "added")
	assertEqual(t, []Entity{}, result.Changed, "changed")
	assertEqual(t, []Entity{}, result.Removed, "removed")
	assertEqual(t, new, result.Merged, "merged")
}

func TestOnlyRemove(t *testing.T) {
	old := []Entity{{ID: "1", Name: "Gone", Value: 99}}
	new := []Entity{}

	merger := NewWithResult[Entity, string]().
		WithKeySelector(func(e Entity) string { return e.ID }).
		WithMergeStrategy(func(old, new Entity) (Entity, bool) {
			return new, true
		})

	result := merger.Merge(old, new)

	assertEqual(t, []Entity{}, result.Added, "added")
	assertEqual(t, []Entity{}, result.Changed, "changed")
	assertEqual(t, old, result.Removed, "removed")
	assertEqual(t, []Entity{}, result.Merged, "merged")
}

func assertEqual[T any](t *testing.T, expected, actual []T, name string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s mismatch.\nExpected: %#v\nActual:   %#v", name, expected, actual)
	}
}
