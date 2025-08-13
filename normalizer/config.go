package normalizer

type FieldConfig struct {
	Name      string            `json:"name,omitempty"`      // e.g., "items.#.itemName"
	Replace   map[string]string `json:"replace,omitempty"`   // e.g., {"from": " - ", "to": ": "}
	Normalize []string          `json:"normalize,omitempty"` // e.g., ["stripBrackets", "titlecase"]
	Extract   string            `json:"extract,omitempty"`   // e.g., "year"
	Format    string            `json:"format,omitempty"`    // e.g., "{{title}} {{year}}"
	From      map[string]string `json:"from,omitempty"`      // e.g., {"title": "...", "year": "..."}
	SaveAs    string            `json:"saveAs,omitempty"`    // e.g., "items.#.metadata.title"
}

type Config struct {
	Root       string        `json:"root"`   // Path to input data file
	Fields     []FieldConfig `json:"fields"` // List of normalization rules
	Verbose    bool          `json:"verbose,omitempty"`
	OutputPath string        `json:"outputPath,omitempty"`
}
