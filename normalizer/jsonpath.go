package normalizer

import (
	"fmt"
	"strings"
)

type ResolveOptions struct {
	InjectNilOnMissing bool
}

func ResolvePath(json any, path string) ([]any, error) {
	vals, _, err := ResolvePathWithOptions(json, path, ResolveOptions{})
	return vals, err
}

func ResolvePathWithOptions(
	json any,
	path string,
	opts ResolveOptions,
) ([]any, []error, error) {
	tokens := tokenize(path)
	return walk(json, tokens, []string{}, opts)
}

func walk(
	data any,
	path []string,
	trail []string,
	opts ResolveOptions,
) ([]any, []error, error) {
	if len(path) == 0 {
		// Reached the leaf node
		return []any{data}, nil, nil
	}

	token := path[0]
	rest := path[1:]
	trail = append(trail, token)

	switch d := data.(type) {
	case map[string]any:
		val, ok := d[token]
		if !ok {
			msg := fmt.Errorf("field %q not found in object at path %q", token, strings.Join(trail, "."))
			if opts.InjectNilOnMissing {
				return []any{nil}, []error{msg}, nil
			}
			return nil, nil, msg
		}
		return walk(val, rest, trail, opts)

	case []any:
		if token == "#" {
			var results []any
			var allErrors []error
			for i, item := range d {
				itemTrail := append(trail[:len(trail)-1], fmt.Sprintf("[%d]", i))
				vals, errs, err := walk(item, rest, itemTrail, opts)
				if err != nil {
					if opts.InjectNilOnMissing {
						allErrors = append(allErrors, fmt.Errorf("error at path %q: %w", strings.Join(itemTrail, "."), err))
						results = append(results, nil)
						continue
					}
					return nil, nil, fmt.Errorf("error at path %q: %w", strings.Join(itemTrail, "."), err)
				}
				if len(errs) > 0 {
					allErrors = append(allErrors, errs...)
				}

				results = append(results, vals...)
			}

			return results, allErrors, nil
		}

		return nil, nil, fmt.Errorf("unexpected token %q for array at path %q (expected '#')", token, strings.Join(trail, "."))

	default:
		msg := fmt.Errorf("unexpected structure at path %q; cannot continue", strings.Join(trail, "."))
		if opts.InjectNilOnMissing {
			return []any{nil}, []error{msg}, nil
		}

		return nil, nil, msg
	}
}

func tokenize(path string) []string {
	return strings.Split(path, ".")
}
