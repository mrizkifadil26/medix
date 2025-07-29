package normalizer

type NormalizeJob struct {
	Input  any      `json:"input"`            // string | []string | path to JSON
	Output string   `json:"output,omitempty"` // optional output path
	Steps  []string `json:"steps"`            // normalization steps
}
