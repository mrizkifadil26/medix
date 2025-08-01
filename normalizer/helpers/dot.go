package normalizer

import "strings"

func DotToSpace(s string) string {
	return strings.ReplaceAll(s, ".", " ")
}
