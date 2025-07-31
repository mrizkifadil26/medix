package classifier

type Config struct {
	Inputs  []string `json:"inputs" yaml:"inputs" validate:"required"`
	Outputs []Rule   `json:"outputs" yaml:"outputs" validate:"required,dive"`
}

type Rule struct {
	Name   string                 `json:"name" yaml:"name" validate:"required"`
	Output string                 `json:"output" yaml:"output" validate:"required"`
	Match  map[string]interface{} `json:"match" yaml:"match"`                 // entry must match this
	Set    map[string]interface{} `json:"set,omitempty" yaml:"set,omitempty"` // inject to entry + top-level
}

type Entry map[string]interface{}

type Output struct {
	SourceCount int                      `json:"source_count"`
	Items       []map[string]interface{} `json:"items"`
}
