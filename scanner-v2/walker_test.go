package scannerV2_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	scannerV2 "github.com/mrizkifadil26/medix/scanner-v2"
)

/* func setupTestData(t *testing.T) string {
	root := t.TempDir()

	// Create: root/Action/Inception (2010)/Inception.2010.mkv
	movieDir := filepath.Join(root, "Action", "Inception (2010)")
	require.NoError(t, os.MkdirAll(movieDir, 0755))

	movieFile := filepath.Join(movieDir, "Inception.2010.mkv")
	require.NoError(t, os.WriteFile(movieFile, []byte("mock data"), 0644))

	return root
} */

func setupTestData(t *testing.T) string {
	dir := t.TempDir()

	// Create structure:
	// dir/
	//   file1.txt
	//   sub/
	//     file2.txt
	if err := os.WriteFile(filepath.Join(dir, "file1.txt"), []byte("abc"), 0644); err != nil {
		t.Fatal(err)
	}

	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(sub, "file2.txt"), []byte("12345"), 0644); err != nil {
		t.Fatal(err)
	}

	return dir
}

/*
func TestWalkDirs(t *testing.T) {
	root := setupTestData(t)

	var visited []string
	err := scannerV2.WalkDirs(root, scannerV2.WalkOptions{
		MaxDepth: 2,
	}, func(path string, entries []os.DirEntry) {
		rel, _ := filepath.Rel(root, path)
		visited = append(visited, rel)
	})
	require.NoError(t, err)

	require.Contains(t, visited, ".") // root
	require.Contains(t, visited, "Action")
	require.Contains(t, visited, filepath.Join("Action", "Inception (2010)"))
}

func TestWalkFiles(t *testing.T) {
	root := setupTestData(t)

	var files []string
	var totalSize int64
	err := scannerV2.WalkFiles(root, scannerV2.WalkOptions{
		MaxDepth: 2,
		Exts:     []string{".mkv"},
	}, func(path string, size int64) {
		rel, _ := filepath.Rel(root, path)
		files = append(files, rel)
		totalSize += size
	})
	require.NoError(t, err)

	require.Len(t, files, 1)
	require.Equal(t, filepath.Join("Action", "Inception (2010)", "Inception.2010.mkv"), files[0])
	require.Equal(t, int64(len("mock data")), totalSize)
}

func TestWalkFiles_IgnoresOtherExtensions(t *testing.T) {
	root := setupTestData(t)

	// Add a non-matching file
	txtFile := filepath.Join(root, "Action", "Inception (2010)", "readme.txt")
	require.NoError(t, os.WriteFile(txtFile, []byte("ignore me"), 0644))

	var files []string
	err := scannerV2.WalkFiles(root, scannerV2.WalkOptions{
		MaxDepth: 2,
		Exts:     []string{".mkv"},
	}, func(path string, size int64) {
		files = append(files, path)
	})
	require.NoError(t, err)

	require.Len(t, files, 1)
	require.True(t, strings.HasSuffix(files[0], "Inception.2010.mkv"))
}

func TestWalkDirs_MaxDepthLimit(t *testing.T) {
	root := setupTestData(t)

	// maxDepth = 1 should skip "Action/Inception (2010)"
	var visited []string
	err := scannerV2.WalkDirs(root, scannerV2.WalkOptions{
		MaxDepth: 1,
	}, func(path string, entries []os.DirEntry) {
		rel, _ := filepath.Rel(root, path)
		visited = append(visited, rel)
	})
	require.NoError(t, err)

	require.Contains(t, visited, ".")
	require.Contains(t, visited, "Action")
	require.NotContains(t, visited, filepath.Join("Action", "Inception (2010)"))
}

func TestWalkFiles_MaxDepthLimit(t *testing.T) {
	root := setupTestData(t)

	var files []string
	err := scannerV2.WalkFiles(root, scannerV2.WalkOptions{
		MaxDepth: 1,
		Exts:     []string{".mkv"},
	}, func(path string, size int64) {
		files = append(files, path)
	})
	require.NoError(t, err)

	require.Len(t, files, 0)
}

func TestWalkDirs_EmptyFolder(t *testing.T) {
	root := t.TempDir() // empty

	var visited []string
	err := scannerV2.WalkDirs(root, scannerV2.WalkOptions{
		MaxDepth: 0,
	}, func(path string, entries []os.DirEntry) {
		rel, _ := filepath.Rel(root, path)
		visited = append(visited, rel)
	})
	require.NoError(t, err)

	require.Equal(t, []string{"."}, visited) // Only root is visited
}

func TestWalkFiles_EmptyFolder(t *testing.T) {
	root := t.TempDir() // empty

	var files []string
	err := scannerV2.WalkFiles(root, scannerV2.WalkOptions{
		MaxDepth: 1,
		Exts:     []string{".mkv"},
	}, func(path string, size int64) {
		files = append(files, path)
	})
	require.NoError(t, err)

	require.Empty(t, files)
}

func TestWalkFiles_NoExtFilter(t *testing.T) {
	root := setupTestData(t)

	// Add additional file with a different extension
	altFile := filepath.Join(root, "Action", "Inception (2010)", "poster.jpg")
	require.NoError(t, os.WriteFile(altFile, []byte("image"), 0644))

	var files []string
	err := scannerV2.WalkFiles(root, scannerV2.WalkOptions{
		MaxDepth: 2,
		Exts:     []string{},
	}, func(path string, size int64) {
		rel, _ := filepath.Rel(root, path)
		files = append(files, rel)
	})
	require.NoError(t, err)

	require.Len(t, files, 2)
	require.Contains(t, files, filepath.Join("Action", "Inception (2010)", "Inception.2010.mkv"))
	require.Contains(t, files, filepath.Join("Action", "Inception (2010)", "poster.jpg"))
}

func TestWalkDirs_SkipEmpty(t *testing.T) {
	root := t.TempDir()

	// Create a non-empty directory with one file
	nonEmptyDir := filepath.Join(root, "NonEmpty")
	require.NoError(t, os.MkdirAll(nonEmptyDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(nonEmptyDir, "file.mkv"), []byte("data"), 0644))

	// Create an empty directory
	emptyDir := filepath.Join(root, "Empty")
	require.NoError(t, os.MkdirAll(emptyDir, 0755))

	var visited []string
	err := scannerV2.WalkDirs(root, scannerV2.WalkOptions{
		MaxDepth:  1,
		SkipEmpty: true,
	}, func(path string, entries []os.DirEntry) {
		rel, _ := filepath.Rel(root, path)
		visited = append(visited, rel)
	})
	require.NoError(t, err)

	require.Contains(t, visited, ".")        // root
	require.Contains(t, visited, "NonEmpty") // has file
	require.NotContains(t, visited, "Empty") // is truly empty
}

func TestWalkDirs_OnlyLeaf(t *testing.T) {
	root := t.TempDir()

	// Parent -> Child -> Grandchild
	require.NoError(t, os.MkdirAll(filepath.Join(root, "Parent", "Child", "Grandchild"), 0755))

	var visited []string
	err := scannerV2.WalkDirs(root, scannerV2.WalkOptions{
		MaxDepth: -1,
		OnlyLeaf: true,
	}, func(path string, entries []os.DirEntry) {
		rel, _ := filepath.Rel(root, path)
		visited = append(visited, rel)
	})
	require.NoError(t, err)

	// Should only visit leaf dirs: Grandchild
	require.Contains(t, visited, filepath.Join("Parent", "Child", "Grandchild"))
	require.NotContains(t, visited, "Parent")
	require.NotContains(t, visited, filepath.Join("Parent", "Child"))
}
*/

func TestIncludeErrors_CollectsErrors(t *testing.T) {
	root := setupTestData(t)

	opts := scannerV2.WalkOptions{
		IncludeStats:  true,
		IncludeErrors: true,
		StopOnError:   false,
		SkipOnError:   true,
		MaxDepth:      -1,
	}
	w := scannerV2.NewWalker(context.Background(), opts)

	// fmt.Println(opts.PrettyPrint())

	// Simulate error for first file
	w.OnVisitFile = func(path string, size int64) error {
		if filepath.Base(path) == "file1.txt" {
			return errors.New("test file error")
		}

		return nil
	}

	err := w.Walk(root)
	if err != nil {
		t.Fatalf("unexpected walk error: %v", err)
	}

	stats := w.GetStats()
	if stats.ErrorsCount != 1 {
		t.Errorf("expected ErrorsCount=1, got %d", stats.ErrorsCount)
	}
	if len(stats.Errors) != 1 {
		t.Errorf("expected 1 stored error, got %d", len(stats.Errors))
	}
	if stats.Errors[0].Error() != "test file error" {
		t.Errorf("unexpected error content: %v", stats.Errors[0])
	}
}

func TestIncludeErrors_False_NoStorage(t *testing.T) {
	root := setupTestData(t)

	opts := scannerV2.WalkOptions{
		IncludeStats:  true,
		IncludeErrors: false,
		StopOnError:   false,
		SkipOnError:   true,
		MaxDepth:      -1,
	}
	w := scannerV2.NewWalker(context.Background(), opts)

	w.OnVisitFile = func(path string, size int64) error {
		return errors.New("err without storage")
	}

	err := w.Walk(root)
	if err != nil {
		t.Fatalf("unexpected walk error: %v", err)
	}

	stats := w.GetStats()
	if stats.ErrorsCount == 0 {
		t.Errorf("expected ErrorsCount > 0")
	}
	if len(stats.Errors) != 0 {
		t.Errorf("expected 0 stored errors when IncludeErrors=false, got %d", len(stats.Errors))
	}
}

func TestStopOnError_StopsImmediately(t *testing.T) {
	root := setupTestData(t)

	opts := scannerV2.WalkOptions{
		IncludeStats:  true,
		IncludeErrors: true,
		StopOnError:   true,
		MaxDepth:      -1,
	}
	w := scannerV2.NewWalker(context.Background(), opts)

	called := 0
	w.OnVisitFile = func(path string, size int64) error {
		called++
		return errors.New("stop now")
	}

	err := w.Walk(root)
	if err == nil || err.Error() != "stop now" {
		t.Fatalf("expected stop error, got %v", err)
	}
	if called != 1 {
		t.Errorf("expected to stop after first file, got %d calls", called)
	}
}

func TestSkipOnError_Skips(t *testing.T) {
	root := setupTestData(t)

	opts := scannerV2.WalkOptions{
		IncludeStats:  true,
		IncludeErrors: true,
		StopOnError:   false,
		SkipOnError:   true,
		MaxDepth:      -1,
	}
	w := scannerV2.NewWalker(context.Background(), opts)

	called := 0
	w.OnVisitFile = func(path string, size int64) error {
		called++
		return errors.New("skip this")
	}

	err := w.Walk(root)
	if err != nil {
		t.Fatalf("unexpected walk error: %v", err)
	}
	if called == 0 {
		t.Errorf("expected files to be visited")
	}
	if w.Stats.ErrorsCount == 0 {
		t.Errorf("expected error count increment")
	}
}
