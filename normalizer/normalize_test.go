package normalizer_test

import (
	"encoding/json"
	"testing"

	"github.com/mrizkifadil26/medix/model"
	normalizer "github.com/mrizkifadil26/medix/normalizer"
	helpers "github.com/mrizkifadil26/medix/normalizer/helpers"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeBasicSteps(t *testing.T) {
	n := normalizer.New()

	tests := []struct {
		input string
		steps []string
		want  string
	}{
		{
			input: "Hello World.mkv",
			steps: []string{"stripExtension", "spaceToDash", "toLower"},
			want:  "hello-world",
		},
		{
			input: "[YTS] Inception.2010.1080p.BluRay.x264.mkv",
			steps: []string{"stripExtension", "stripBrackets", "dotToSpace", "collapseDashes"},
			want:  "Inception 2010 1080p BluRay x264",
		},
		{
			input: "My-File--Name",
			steps: []string{"collapseDashes"},
			want:  "My-File-Name",
		},
	}

	for _, tt := range tests {
		got := n.Normalize(tt.input, tt.steps)
		assert.Equal(t, tt.want, got, "input: %s", tt.input)
	}
}

func TestRun_StringInput(t *testing.T) {
	n := normalizer.New()

	input := "Movie.Name.2020.HDRip"
	steps := []string{"dotToSpace", "toLower"}

	result, err := n.Run(input, steps)
	assert.NoError(t, err)
	assert.Equal(t, "movie name 2020 hdrip", result)
}

func TestRun_ArrayInput(t *testing.T) {
	n := normalizer.New()

	input := []any{"A.Title", "Another.One"}
	steps := []string{"dotToSpace", "toLower"}

	result, err := n.Run(input, steps)
	assert.NoError(t, err)

	expected := []string{"a title", "another one"}
	assert.Equal(t, expected, result)
}

func TestRun_UnsupportedType(t *testing.T) {
	n := normalizer.New()

	input := 42
	_, err := n.Run(input, []string{})
	assert.Error(t, err)
}

func TestRun_WithScanJSON(t *testing.T) {
	// Load scan JSON from file
	data := `{
		"items": [
			{ "name": "Some.Movie.2020.1080p.mkv" },
			{ "name": "[YTS] Another.Film.2018.mkv" }
		]
	}`

	var parsed model.MediaOutput
	_ = json.Unmarshal([]byte(data), &parsed)

	var names []any
	for _, item := range parsed.Items {
		names = append(names, item.Name)
	}

	n := normalizer.New()
	steps := []string{"stripExtension", "dotToSpace", "stripBrackets", "toLower"}

	result, err := n.Run(names, steps)
	assert.NoError(t, err)

	expected := []string{"some movie 2020 1080p", "another film 2018"}
	assert.Equal(t, expected, result)
}

func TestExtractTitleYear(t *testing.T) {
	tests := []struct {
		name          string
		expectedTitle string
		expectedYear  int
	}{
		{"Interstellar (2014)", "Interstellar", 2014},
		{"The Godfather (1972)", "The Godfather", 1972},
		{"No Year Present", "No Year Present", 0},
		{"Edge Case (20A4)", "Edge Case (20A4)", 0},
	}

	for _, tt := range tests {
		title, year := helpers.ExtractTitleYear(tt.name)
		if title != tt.expectedTitle || year != tt.expectedYear {
			t.Errorf("❌ For %q → got (%q, %d), expected (%q, %d)", tt.name, title, year, tt.expectedTitle, tt.expectedYear)
		}
	}
}

func TestNormalizeUnicode(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Café", "Cafe"},
		{"naïve", "naive"},
		{"crème brûlée", "creme brulee"},
		{"½-life", "1-2-life"},
		{"“Quote”", `"Quote"`}, // if smart quotes are removed
		{"‘Start’ and ’End’", "'Start' and 'End'"},
		{"rock³", "rock3"},
	}

	for _, tt := range tests {
		got := helpers.NormalizeUnicode(tt.input)
		if got != tt.want {
			t.Errorf("NormalizeUnicode failed:\ninput: %q\ngot:   %q\nwant:  %q", tt.input, got, tt.want)
		}
	}
}

func TestStripExtension(t *testing.T) {
	input := "movie.name.2023.mkv"
	expected := "movie.name.2023"
	result := helpers.StripExtension(input)
	if result != expected {
		t.Errorf("StripExtension failed: got %q, want %q", result, expected)
	}
}

func TestStripBrackets(t *testing.T) {
	input := "Movie Title (2023) [BluRay]"
	expected := "Movie Title  "
	result := helpers.StripBrackets(input)
	if result != expected {
		t.Errorf("StripBrackets failed: got %q, want %q", result, expected)
	}
}

func TestReplaceSpecialChars(t *testing.T) {
	input := "Rock & Roll / Funk + Soul?"
	expected := "Rock and Roll - Funk  Soul"
	result := helpers.ReplaceSpecialChars(input)
	if result != expected {
		t.Errorf("ReplaceSpecialChars failed: got %q, want %q", result, expected)
	}
}

func TestCollapseDashes(t *testing.T) {
	input := "Hello---World___Test"
	expected := "Hello-World-Test"
	result := helpers.CollapseDashes(input)
	if result != expected {
		t.Errorf("CollapseDashes failed: got %q, want %q", result, expected)
	}
}

func TestToLower(t *testing.T) {
	input := "HeLLo WorlD"
	expected := "hello world"
	result := helpers.ToLower(input)
	if result != expected {
		t.Errorf("ToLower failed: got %q, want %q", result, expected)
	}
}

func TestSpaceToDash(t *testing.T) {
	input := "  The Matrix Reloaded  "
	expected := "The-Matrix-Reloaded"
	result := helpers.SpaceToDash(input)
	if result != expected {
		t.Errorf("SpaceToDash failed: got %q, want %q", result, expected)
	}
}

func TestRemoveKnownPrefixes(t *testing.T) {
	input := "[1080p] The Movie [BluRay]"
	expected := " The Movie "
	result := helpers.RemoveKnownPrefixes(input)
	if result != expected {
		t.Errorf("RemoveKnownPrefixes failed: got %q, want %q", result, expected)
	}
}

func TestDotToSpace(t *testing.T) {
	input := "The.Matrix.1999"
	expected := "The Matrix 1999"
	result := helpers.DotToSpace(input)
	if result != expected {
		t.Errorf("DotToSpace failed: got %q, want %q", result, expected)
	}
}
