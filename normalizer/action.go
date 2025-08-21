package normalizer

type Action struct {
	Type      string // replace, transform, extract, format
	Target    string
	HashIndex []int
	Params    map[string]any // action-specific params
}
