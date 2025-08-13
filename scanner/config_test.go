package scanner_test

import (
	"testing"

	"github.com/mrizkifadil26/medix/scanner"
	"github.com/mrizkifadil26/medix/utils"
)

func TestDefaultConfigOnly(t *testing.T) {
	defaultCfg := scanner.DefaultConfig()

	if defaultCfg.Root != "" {
		t.Errorf("Expected default Root to be empty or nil, got %v", defaultCfg.Root)
	}

	if defaultCfg.Options == nil {
		t.Fatal("Expected Options to be non-nil in default config")
	}
}

func TestOverrideWithFileConfig(t *testing.T) {
	defaultCfg := scanner.DefaultConfig()

	fileRoot := "/file/root"
	fileMode := "dirs"
	fileCfg := scanner.Config{
		Root: fileRoot,
		Options: &scanner.ScanOptions{
			Mode:  fileMode,
			Depth: 5,
		},
	}

	merged, err := utils.Merge(defaultCfg, fileCfg, utils.MergeOptions{
		Overwrite: true,
		Recursive: true,
	})
	if err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	if merged.Root != fileRoot {
		t.Errorf("Expected Root to be %q, got %v", fileRoot, merged.Root)
	}

	if merged.Options == nil || merged.Options.Mode != fileMode || merged.Options.Depth != 5 {
		t.Errorf("Expected Options.Mode to be %q and Depth 5, got %+v", fileMode, merged.Options)
	}
}

func TestOverrideWithCLIConfig(t *testing.T) {
	defaultCfg := scanner.DefaultConfig()

	cliRoot := "/cli/root"
	cliFormat := "yaml"
	cliCfg := scanner.Config{
		Root: cliRoot,
		Output: &scanner.OutputOptions{
			Format: cliFormat,
		},
	}

	merged, err := utils.Merge(defaultCfg, cliCfg, utils.MergeOptions{
		Overwrite: true,
		Recursive: true,
	})
	if err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	if merged.Root != cliRoot {
		t.Errorf("Expected Root to be %q, got %v", cliRoot, merged.Root)
	}

	if merged.Output == nil || merged.Output.Format != cliFormat {
		t.Errorf("Expected Output.Format to be %q, got %v", cliFormat, merged.Output.Format)
	}
}

func TestFileThenCLIOverride(t *testing.T) {
	defaultCfg := scanner.DefaultConfig()

	fileRoot := "/file/root"
	fileMode := "dirs"
	fileCfg := scanner.Config{
		Root: fileRoot,
		Options: &scanner.ScanOptions{
			Mode:  fileMode,
			Depth: 5,
		},
	}

	cliRoot := "/cli/root"
	cliFormat := "yaml"
	cliCfg := scanner.Config{
		Root: cliRoot,
		Output: &scanner.OutputOptions{
			Format: cliFormat,
		},
	}

	mergedFile, err := utils.Merge(defaultCfg, fileCfg, utils.MergeOptions{
		Overwrite: true,
		Recursive: true,
	})
	if err != nil {
		t.Fatalf("Merge file failed: %v", err)
	}

	finalMerged, err := utils.Merge(mergedFile, cliCfg, utils.MergeOptions{
		Overwrite: true,
		Recursive: true,
	})
	if err != nil {
		t.Fatalf("Merge CLI failed: %v", err)
	}

	if finalMerged.Root != cliRoot {
		t.Errorf("Expected Root to be %q, got %v", cliRoot, finalMerged.Root)
	}

	if finalMerged.Options == nil || finalMerged.Options.Mode != fileMode || finalMerged.Options.Depth != 5 {
		t.Errorf("Expected Options.Mode to be %q and Depth 5, got %+v", fileMode, finalMerged.Options)
	}

	if finalMerged.Output == nil || finalMerged.Output.Format != cliFormat {
		t.Errorf("Expected Output.Format to be %q, got %v", cliFormat, finalMerged.Output.Format)
	}
}

func TestMergeWithEmptyOverride(t *testing.T) {
	defaultCfg := scanner.DefaultConfig()

	emptyCfg := scanner.Config{}

	merged, err := utils.Merge(defaultCfg, emptyCfg, utils.MergeOptions{
		Overwrite: true,
		Recursive: true,
	})
	if err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	// Expect merged to equal defaultCfg (no fields overwritten)
	if merged.Options == nil {
		t.Fatal("Expected Options to remain set after merge with empty override")
	}
	// Add more assertions as needed...
}

func TestPartialNestedOverride(t *testing.T) {
	defaultCfg := scanner.DefaultConfig()

	fileCfg := scanner.Config{
		Options: &scanner.ScanOptions{
			Mode:  "files",
			Depth: 3,
		},
	}

	cliCfg := scanner.Config{
		Options: &scanner.ScanOptions{
			Depth: 7, // override only Depth
		},
	}

	mergedFile, err := utils.Merge(defaultCfg, fileCfg, utils.MergeOptions{
		Overwrite: true,
		Recursive: true,
	})
	if err != nil {
		t.Fatalf("File merge failed: %v", err)
	}

	finalMerged, err := utils.Merge(mergedFile, cliCfg, utils.MergeOptions{
		Overwrite: true,
		Recursive: true,
	})
	if err != nil {
		t.Fatalf("CLI merge failed: %v", err)
	}

	if finalMerged.Options.Mode != "files" {
		t.Errorf("Expected Mode 'files', got %q", finalMerged.Options.Mode)
	}
	if finalMerged.Options.Depth != 7 {
		t.Errorf("Expected Depth 7, got %d", finalMerged.Options.Depth)
	}
}

func TestNilNestedStructOverride(t *testing.T) {
	defaultCfg := scanner.DefaultConfig()

	defaultCfg.Options = &scanner.ScanOptions{
		Mode:  "default-mode",
		Depth: 2,
	}

	cliCfg := scanner.Config{
		Options: nil, // should NOT clear default Options
	}

	merged, err := utils.Merge(defaultCfg, cliCfg, utils.MergeOptions{
		Overwrite: true,
		Recursive: true,
	})
	if err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	if merged.Options == nil {
		t.Fatal("Expected Options to remain after merge with nil override")
	}
	if merged.Options.Mode != "default-mode" {
		t.Errorf("Expected Mode 'default-mode', got %q", merged.Options.Mode)
	}
}

func TestSliceMergeBehavior(t *testing.T) {
	defaultCfg := scanner.DefaultConfig()

	defaultTags := []string{"default", "base"}
	fileTags := []string{"file", "override"}
	cliTags := []string{"cli", "override"}

	defaultCfg.Tags = defaultTags
	fileCfg := scanner.Config{Tags: fileTags}
	cliCfg := scanner.Config{Tags: cliTags}

	mergedFile, err := utils.Merge(defaultCfg, fileCfg, utils.MergeOptions{
		Overwrite: true,
		Recursive: true,
	})
	if err != nil {
		t.Fatalf("File merge failed: %v", err)
	}

	finalMerged, err := utils.Merge(mergedFile, cliCfg, utils.MergeOptions{
		Overwrite: true,
		Recursive: true,
	})
	if err != nil {
		t.Fatalf("CLI merge failed: %v", err)
	}

	// Adjust this depending on your slice merge strategy
	expected := cliTags // usually override replaces slice

	if finalMerged.Tags == nil {
		t.Fatal("Expected Tags to be set")
	}

	if len(finalMerged.Tags) != len(expected) {
		t.Errorf("Expected Tags length %d, got %d", len(expected), len(finalMerged.Tags))
	}
}

/*
func TestValidationFailsWithoutRoot(t *testing.T) {
	cfg := scanner.DefaultConfig()
	cfg.Root = nil // explicitly unset root

	if cfg.Root != nil {
		t.Fatal("Root should be nil for test")
	}

	// Simulate validation logic from your main
	if cfg.Root == nil || *cfg.Root == "" {
		// Pass test if detected
		return
	}

	t.Fatal("Validation failed to detect missing Root")
}
*/

/*
func TestApplyDefaults(t *testing.T) {
	cfg := scanner.Config{}
	cfg.ApplyDefaults()

	if cfg.Options.Mode != "files" {
		t.Errorf("Expected Mode to default to 'files', got '%s'", cfg.Options.Mode)
	}

	if cfg.Options.Depth != 1 {
		t.Errorf("Expected Depth to default to 1, got %d", cfg.Options.Depth)
	}
}
*/
