package scannerV2_test

import (
	"os"
	"reflect"
	"testing"

	scannerV2 "github.com/mrizkifadil26/medix/scanner-v2"
)

func TestParseCLI(t *testing.T) {
	// Save original args and restore after test
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{
		"scan-v2", // fake binary name
		"--root", "/test/path",
		"--mode", "dirs",
		"--exts", ".mp4,.avi",
		"--exclude", "temp,backup",
		"--depth", "3",
		"--only-leaf",
		"--leaf-depth", "2",
		"--skip-empty",
		"--concurrency", "4",
		"--verbose",
	}

	args := scannerV2.ParseCLI()

	if args.Config.Root != "/test/path" {
		t.Errorf("Expected root to be '/test/path', got '%s'", args.Config.Root)
	}

	if args.Config.Options.Mode != "dirs" {
		t.Errorf("Expected mode to be 'dirs', got '%s'", args.Config.Options.Mode)
	}

	expectedExts := []string{".mp4", ".avi"}
	if !reflect.DeepEqual(args.Config.Options.Exts, expectedExts) {
		t.Errorf("Expected exts to be %v, got %v", expectedExts, args.Config.Options.Exts)
	}

	expectedExclude := []string{"temp", "backup"}
	if !reflect.DeepEqual(args.Config.Options.Exclude, expectedExclude) {
		t.Errorf("Expected exclude to be %v, got %v", expectedExclude, args.Config.Options.Exclude)
	}

	if args.Config.Options.Depth != 3 {
		t.Errorf("Expected depth to be 3, got %d", args.Config.Options.Depth)
	}

	if !args.Config.Options.OnlyLeaf {
		t.Error("Expected OnlyLeaf to be true")
	}

	if args.Config.Options.LeafDepth != 2 {
		t.Errorf("Expected LeafDepth to be 2, got %d", args.Config.Options.LeafDepth)
	}

	if !args.Config.Options.SkipEmpty {
		t.Error("Expected SkipEmpty to be true")
	}

	if args.Config.Options.Concurrency != 4 {
		t.Errorf("Expected Concurrency to be 4, got %d", args.Config.Options.Concurrency)
	}

	if !args.Config.Options.Verbose {
		t.Error("Expected Verbose to be true")
	}
}
