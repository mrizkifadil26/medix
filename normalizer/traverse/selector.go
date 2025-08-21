package traverse

import (
	"strings"
)

// Selector is a compiled selector pattern
type Selector struct {
	// parts []string
	Tokens    []string // e.g., ["items", "#", "name"]
	HashIndex []int    // precompute positions of `#`
}

// CompileSelector creates a new Selector from string
func CompileSelector(path string) Selector {
	tokens := strings.Split(path, ".")
	hashIdx := []int{}
	for i, t := range tokens {
		if t == "#" {
			hashIdx = append(hashIdx, i)
		}
	}

	return Selector{
		Tokens:    tokens,
		HashIndex: hashIdx,
	}
}

// Match checks if a JSON path matches this selector
func (s Selector) Match(path []string) bool {
	if len(path) != len(s.Tokens) {
		return false
	}

	for i, tok := range s.Tokens {
		if tok != "#" && tok != path[i] {
			return false
		}
	}

	return true

	/*
		for i := range s.parts {
			switch s.parts[i] {
			case "#": // wildcard for array index
				if _, err := strconv.Atoi(path[i]); err != nil {
					return false // only match if path[i] is a number
				}

				// accept any numeric string
				continue
			case "*": // wildcard for any key
				continue
			default:
				if s.parts[i] != path[i] {
					return false
				}
			}
		}

		return true
	*/
}

// SelectorSet holds multiple selectors
type SelectorSet struct {
	selectors []Selector
}

// NewSelectorSet compiles many patterns
func NewSelectorSet(patterns []string) SelectorSet {
	set := SelectorSet{}
	for _, p := range patterns {
		set.selectors = append(set.selectors, CompileSelector(p))
	}

	return set
}

// Match returns true if any selector matches
func (ss SelectorSet) Match(path []string) bool {
	for _, sel := range ss.selectors {
		if sel.Match(path) {
			return true
		}
	}

	return false
}
