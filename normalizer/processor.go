package normalizer

import (
	"fmt"
	"strconv"
	"strings"
)

type ErrorHandlingOptions struct {
	ContinueOnError bool
	CollectErrors   bool
}

func Process(
	input any,
	fields []FieldConfig,
	reg *OperatorRegistry,
	opts ErrorHandlingOptions,
) (any, error) {
	root, ok := input.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("input must be a JSON object")
	}

	var allErrors []error

	// 1. Apply all modifiers first
	for _, field := range fields {
		if field.Name == "" {
			continue
		}

		if err := applyModifier(root, field, reg); err != nil {
			wrappedErr := fmt.Errorf("modifier failed for field %s: %w", field.Name, err)

			if opts.CollectErrors {
				allErrors = append(allErrors, wrappedErr)
			}

			if !opts.ContinueOnError {
				if opts.CollectErrors {
					return nil, combineErrors(allErrors)
				}

				return nil, wrappedErr
			}
		}
	}

	// 2. Apply all constructors next
	for _, field := range fields {
		if field.Format == "" {
			continue
		}

		if err := applyConstructor(root, field, reg); err != nil {
			wrappedErr := fmt.Errorf("constructor failed for format %s: %w", field.Format, err)

			if opts.CollectErrors {
				allErrors = append(allErrors, wrappedErr)
			}

			if !opts.ContinueOnError {
				if opts.CollectErrors {
					return nil, combineErrors(allErrors)
				}

				return nil, wrappedErr
			}
		}
	}

	if opts.CollectErrors && len(allErrors) > 0 {
		if rootMap, ok := input.(map[string]any); ok {
			var lines []string
			for _, err := range allErrors {
				lines = append(lines, err.Error())
			}
			rootMap["_errors"] = lines
		}

		return input, combineErrors(allErrors)
	}

	return input, nil
}

func combineErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var sb strings.Builder
	for _, err := range errs {
		sb.WriteString("- ")
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}

	return fmt.Errorf("multiple errors:\n%s", sb.String())
}

type ExpandedField struct {
	ResolvedPath string
	IndexMap     map[string]int
}

func applyModifier(
	root map[string]any,
	field FieldConfig,
	reg *OperatorRegistry,
) error {
	values, err := ResolvePath(root, field.Name)
	if err != nil {
		return fmt.Errorf("traverse %s failed: %w", field.Name, err)
	}

	// Build ops map
	fieldOps := map[string]any{}
	if field.Replace != nil {
		fieldOps["replace"] = field.Replace
	}

	if len(field.Normalize) > 0 {
		fieldOps["normalize"] = toAnySlice(field.Normalize)
	}

	if field.Extract != "" {
		fieldOps["extract"] = field.Extract
	}

	for i, val := range values {
		strVal, ok := val.(string)
		if !ok {
			continue
		}

		modified, err := reg.ApplyOperators(strVal, fieldOps)
		if err != nil {
			return fmt.Errorf("applyOperators failed on value %q (field: %s): %w", strVal, field.Name, err)
		}

		if field.SaveAs != "" {
			if err := SetPath(root, field.SaveAs, modified, i); err != nil {
				return fmt.Errorf("failed to save result for value %q (field: %s): %w", strVal, field.Name, err)
			}
		}
	}

	return nil
}

func applyConstructor(
	root map[string]any,
	field FieldConfig,
	reg *OperatorRegistry,
) error {
	if field.Format == "" || len(field.From) == 0 {
		return fmt.Errorf("constructor needs both 'format' and 'from'")
	}

	// Resolve all field values
	data := map[string]string{}
	for key, path := range field.From {
		vals, err := ResolvePath(root, path)
		if err != nil {
			return fmt.Errorf("resolve path %q for constructor key %q failed: %w", path, key, err)
		}

		if len(vals) == 0 {
			// continue
			return fmt.Errorf("no values found at path %q for constructor key %q", path, key)
		}

		if strVal, ok := vals[0].(string); ok {
			data[key] = strVal
		} else {
			return fmt.Errorf("value at path %q for constructor key %q is not a string", path, key)
		}
	}

	// Format the result
	result, err := reg.FormatFunc(field.Format, data)
	if err != nil {
		return fmt.Errorf("format failed for template %q with data %v: %w", field.Format, data, err)
	}

	if field.SaveAs != "" {
		if err := SetPath(root, field.SaveAs, result, 0); err != nil {
			return fmt.Errorf("failed to save result for value %q (field: %s): %w", result, field.Name, err)
		}
	}

	return nil
}

func SetPath(
	root map[string]any,
	path string,
	value any,
	index int, // index from Traverse match
) error {
	segments := strings.Split(path, ".")

	// Resolve `#` to index in path
	for i, s := range segments {
		if s == "#" {
			segments[i] = strconv.Itoa(index)
		}
	}

	// Walk through the path and build if needed
	curr := root
	for i := 0; i < len(segments)-1; i++ {
		key := segments[i]
		nextKey := segments[i+1]

		// Handle array access
		if idx, err := strconv.Atoi(nextKey); err == nil {
			// If current key not found or not a slice, create
			if _, ok := curr[key]; !ok {
				curr[key] = make([]any, idx+1)
			}
			slice, ok := curr[key].([]any)
			if !ok {
				return fmt.Errorf("expected slice at %s", key)
			}
			// Extend slice if needed
			if len(slice) <= idx {
				newSlice := make([]any, idx+1)
				copy(newSlice, slice)
				slice = newSlice
				curr[key] = slice
			}
			// Create map if not already
			if slice[idx] == nil {
				slice[idx] = map[string]any{}
			}
			// Descend into next map
			child, ok := slice[idx].(map[string]any)
			if !ok {
				return fmt.Errorf("expected map in slice at %s[%d]", key, idx)
			}
			curr = child
			i++ // skip next key (already handled)
		} else {
			// Handle map access
			if _, ok := curr[key]; !ok {
				curr[key] = map[string]any{}
			}
			next, ok := curr[key].(map[string]any)
			if !ok {
				return fmt.Errorf("expected map at %s", key)
			}
			curr = next
		}
	}

	// Final segment = key to set
	lastKey := segments[len(segments)-1]
	curr[lastKey] = value
	return nil
}

func toAnySlice(strs []string) []any {
	out := make([]any, len(strs))
	for i, s := range strs {
		out[i] = s
	}

	return out
}
