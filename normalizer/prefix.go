package normalizer

import "strings"

var knownPrefixes = []string{
	"[1080p]", "[720p]", "[WEBRip]", "[BluRay]", "[YIFY]", "[x264]", "[H264]", "[AC3]",
}

func RemoveKnownPrefixes(s string) string {
	for _, prefix := range knownPrefixes {
		s = strings.ReplaceAll(s, prefix, "")
	}
	return s
}
