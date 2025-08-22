package traverse

import (
	"strconv"
	"strings"
)

// Selector is a compiled selector pattern
type Selector struct {
	Tokens    []string // e.g., ["items", "#", "name"]
	HashIndex []int    // precompute positions of `#`
}

// CompileSelector creates a new Selector from string
func CompileSelector(path string) Selector {
	tokens := strings.Split(path, ".")
	hashIdx := []int{}
	for i, token := range tokens {
		if token == "#" {
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

	for i, token := range s.Tokens {
		if !matchToken(token, path[i]) {
			return false
		}
	}

	return true
}

// SelectorSet holds multiple selectors
type SelectorSet struct {
	selectors []Selector
}

// NewSelectorSet compiles many patterns
func NewSelectorSet(patterns ...string) SelectorSet {
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

func matchToken(token, value string) bool {
	switch token {
	case "#": // numeric index
		_, err := strconv.Atoi(value)
		return err == nil
	case "*": // wildcard
		return true
	default: // exact match
		return token == value
	}
}
