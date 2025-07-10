package sync

import "regexp"

var altIDRegex = regexp.MustCompile(`-alt(?:-\d+)?$`)

func normalizeID(id string) string {
	return altIDRegex.ReplaceAllString(id, "")
}

func isAltVariant(id string) bool {
	return altIDRegex.MatchString(id)
}
