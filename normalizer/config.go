package normalizer

// type FieldConfig struct {
// 	Name      string            `json:"name,omitempty"`      // e.g., "items.#.itemName"
// 	Replace   map[string]string `json:"replace,omitempty"`   // e.g., {"from": " - ", "to": ": "}
// 	Normalize []string          `json:"normalize,omitempty"` // e.g., ["stripBrackets", "titlecase"]
// 	Extract   string            `json:"extract,omitempty"`   // e.g., "year"
// 	Format    string            `json:"format,omitempty"`    // e.g., "{{title}} {{year}}"
// 	From      map[string]string `json:"from,omitempty"`      // e.g., {"title": "...", "year": "..."}
// 	SaveAs    string            `json:"saveAs,omitempty"`    // e.g., "items.#.metadata.title"
// }

type Config struct {
	Root       string  `json:"root"`   // Path to input data file
	Fields     []Field `json:"fields"` // List of normalization rules
	Verbose    bool    `json:"verbose,omitempty"`
	OutputPath string  `json:"outputPath,omitempty"`
}

// Each field and its actions
type Field struct {
	Name    string   `json:"name"`    // JSONPath selector
	Actions []Action `json:"actions"` // List of actions to apply
}

// Action types
type Action struct {
	Type string `json:"type"` // "replace", "transform", "save", "extract", "template", etc.

	// For replace
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`

	// For transform
	Methods []string `json:"methods,omitempty"`

	// For save / extract / template
	Target   string `json:"target,omitempty"`
	Pattern  string `json:"pattern,omitempty"`
	Template string `json:"template,omitempty"`
}
