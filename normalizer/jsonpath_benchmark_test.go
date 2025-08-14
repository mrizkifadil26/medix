package normalizer_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/mrizkifadil26/medix/normalizer"
)

// Load and cache JSON once for benchmark
var testData map[string]any

func init() {
	testSource := "../data/scanner/media/movies.final.json"
	data, err := os.ReadFile("../data/scanner/media/movies.final.json")
	if err != nil {
		panic("failed to read " + testSource + ": " + err.Error())
	}

	if err := json.Unmarshal(data, &testData); err != nil {
		panic("failed to parse test JSON: " + err.Error())
	}
}

func BenchmarkResolvePath_ItemsName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := normalizer.ResolvePath(testData, "items.#.itemName")
		if err != nil {
			b.Fatal(err)
		}
	}
}
