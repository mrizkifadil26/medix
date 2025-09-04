package scanner_test

/* func setupTestData(t *testing.T) string {
	root := t.TempDir()

	// Create: root/Action/Inception (2010)/Inception.2010.mkv
	movieDir := filepath.Join(root, "Action", "Inception (2010)")
	require.NoError(t, os.MkdirAll(movieDir, 0755))

	movieFile := filepath.Join(movieDir, "Inception.2010.mkv")
	require.NoError(t, os.WriteFile(movieFile, []byte("mock data"), 0644))

	return root
} */

/*
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
*/

/*
func TestWalkDirs(t *testing.T) {
	root := setupTestData(t)

	var visited []string
	err := scanner.WalkDirs(root, scanner.WalkOptions{
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
	err := scanner.WalkFiles(root, scanner.WalkOptions{
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
	err := scanner.WalkFiles(root, scanner.WalkOptions{
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
	err := scanner.WalkDirs(root, scanner.WalkOptions{
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
	err := scanner.WalkFiles(root, scanner.WalkOptions{
		MaxDepth: 1,
		Exts:     []string{".mkv"},
	}, func(path string, size int64) {
		files = append(files, path)
	})
	require.NoError(t, err)

	require.Len(t, files, 0)
}
*/
/*
func TestWalkDirs_EmptyFolder(t *testing.T) {
	root := t.TempDir() // empty

	var visitedDirs []string
	var mu sync.Mutex

	walker := &scanner.Walker{
		Context: context.Background(),
		Opts: scanner.WalkOptions{
			MaxDepth: -1,
			OnlyDirs: true,
		},
		OnVisitDir: func(path string, entries []os.DirEntry) error {
			mu.Lock()
			rel, _ := filepath.Rel(root, path)
			visitedDirs = append(visitedDirs, rel)
			mu.Unlock()

			return nil
		},
	}

	err := walker.Walk(root)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}
	require.NoError(t, err)

	require.Equal(t, []string{"."}, visitedDirs) // Only root is visited
}

func TestWalkFiles_EmptyFolder(t *testing.T) {
	root := t.TempDir() // empty

	var files []string
	var mu sync.Mutex

	walker := &scanner.Walker{
		Context: context.Background(),
		Opts: scanner.WalkOptions{
			MaxDepth:  -1,
			OnlyFiles: true,
		},
		OnVisitFile: func(path string, size int64) error {
			mu.Lock()
			files = append(files, path)
			mu.Unlock()

			return nil
		},
	}

	err := walker.Walk(root)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	require.NoError(t, err)

	require.Empty(t, files)
}

func TestWalkFiles_NoExtFilter(t *testing.T) {
	rootDir := t.TempDir()

	// Create this structure:
	// dir/
	//   Action/
	//     Inception (2010)/
	//       Inception.2010.mkv
	//       (poster.jpg will be added by test)
	actionDir := filepath.Join(rootDir, "Action", "Inception (2010)")
	require.NoError(t, os.MkdirAll(actionDir, 0755))

	// Write the main video file
	mkvFile := filepath.Join(actionDir, "Inception.2010.mkv")
	require.NoError(t, os.WriteFile(mkvFile, []byte("movie content"), 0644))

	// Add additional file with a different extension
	altFile := filepath.Join(rootDir, "Action", "Inception (2010)", "poster.jpg")
	require.NoError(t, os.WriteFile(altFile, []byte("image"), 0644))

	var files []string
	var mu sync.Mutex

	walker := &scanner.Walker{
		Context: context.Background(),
		Opts: scanner.WalkOptions{
			MaxDepth:    -1,
			IncludeExts: []string{},
		},
		OnVisitFile: func(path string, size int64) error {
			mu.Lock()
			rel, _ := filepath.Rel(rootDir, path)
			files = append(files, rel)
			mu.Unlock()

			return nil
		},
	}

	err := walker.Walk(rootDir)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	require.NoError(t, err)

	require.Len(t, files, 2)
	require.Contains(t, files, filepath.Join("Action", "Inception (2010)", "Inception.2010.mkv"))
	require.Contains(t, files, filepath.Join("Action", "Inception (2010)", "poster.jpg"))
}

func TestWalkDirs_SkipEmptyDirs(t *testing.T) {
	root := t.TempDir()

	// Create a non-empty directory with one file
	nonEmptyDir := filepath.Join(root, "NonEmpty")
	require.NoError(t, os.MkdirAll(nonEmptyDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(nonEmptyDir, "file.mkv"), []byte("data"), 0644))

	// Create an empty directory
	emptyDir := filepath.Join(root, "Empty")
	require.NoError(t, os.MkdirAll(emptyDir, 0755))

	visitedDirs := make([]string, 0)
	var mu sync.Mutex

	walker := &scanner.Walker{
		Context: context.Background(),
		Opts: scanner.WalkOptions{
			OnlyDirs:       true,
			IncludeHidden:  false,
			SkipEmptyDirs:  true,
			OnlyLeafDirs:   false,
			MaxDepth:       -1,
			IncludeStats:   false,
			EnableProgress: false,
		},
		OnVisitDir: func(path string, entries []os.DirEntry) error {
			mu.Lock()
			rel, _ := filepath.Rel(root, path)
			visitedDirs = append(
				visitedDirs,
				rel,
			)
			mu.Unlock()
			return nil
		},
		// Ignore files by not setting OnVisitFile or setting it to no-op
		OnVisitFile: func(path string, size int64) error {
			// Should never be called in this test if you want only directory walking
			t.Errorf("OnVisitFile called on %s, but should be ignored", path)
			return nil
		},
	}

	err := walker.Walk(root)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}
	require.NoError(t, err)

	require.Contains(t, visitedDirs, ".")        // root
	require.Contains(t, visitedDirs, "NonEmpty") // has file
	require.NotContains(t, visitedDirs, "Empty") // is truly empty
}

func TestWalkDirs_OnlyLeaf(t *testing.T) {
	root := t.TempDir()

	// Parent -> Child -> Grandchild
	require.NoError(t, os.MkdirAll(filepath.Join(root, "Parent", "Child", "Grandchild"), 0755))

	visitedDirs := make([]string, 0)
	var mu sync.Mutex

	walker := &scanner.Walker{
		Context: context.Background(),
		Opts: scanner.WalkOptions{
			IncludeHidden:  false,
			SkipEmptyDirs:  false,
			OnlyLeafDirs:   true,
			MaxDepth:       -1,
			IncludeStats:   false,
			EnableProgress: false,
		},
		OnVisitDir: func(path string, entries []os.DirEntry) error {
			mu.Lock()
			rel, _ := filepath.Rel(root, path)
			visitedDirs = append(
				visitedDirs,
				rel,
			)
			mu.Unlock()
			return nil
		},
		// Ignore files by not setting OnVisitFile or setting it to no-op
		OnVisitFile: func(path string, size int64) error {
			// Should never be called in this test if you want only directory walking
			t.Errorf("OnVisitFile called on %s, but should be ignored", path)
			return nil
		},
	}

	err := walker.Walk(root)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	// Should only visit leaf dirs: Grandchild
	require.Contains(t, visitedDirs, filepath.Join("Parent", "Child", "Grandchild"))
	require.NotContains(t, visitedDirs, "Parent")
	require.NotContains(t, visitedDirs, filepath.Join("Parent", "Child"))
}

func TestIncludeErrors_CollectsErrors(t *testing.T) {
	root := setupTestData(t)

	opts := scanner.WalkOptions{
		IncludeStats:  true,
		IncludeErrors: true,
		StopOnError:   false,
		SkipOnError:   true,
		MaxDepth:      -1,
	}
	w := scanner.NewWalker(context.Background(), opts)

	// Simulate error for first file
	w.OnVisitFile = func(path string, size int64) error {
		t.Logf("Visiting file: %s (size=%d)", filepath.Base(path), size)

		if filepath.Base(path) == "file1.txt" {
			t.Logf("Injecting test error for %s", path)
			return errors.New("test file error")
		}

		return nil
	}

	err := w.Walk(root)
	if err != nil {
		t.Fatalf("unexpected walk error: %v", err)
	}

	stats := w.GetStats()
	t.Logf("Collected stats: %+v", stats)

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

	opts := scanner.WalkOptions{
		IncludeStats:  true,
		IncludeErrors: false,
		StopOnError:   false,
		SkipOnError:   true,
		MaxDepth:      -1,
	}
	w := scanner.NewWalker(context.Background(), opts)

	w.OnVisitFile = func(path string, size int64) error {
		t.Logf("Visiting file: %s", filepath.Base(path))
		return errors.New("err without storage")
	}

	err := w.Walk(root)
	if err != nil {
		t.Fatalf("unexpected walk error: %v", err)
	}

	stats := w.GetStats()
	t.Logf("Collected stats: %+v", stats)

	if stats.ErrorsCount == 0 {
		t.Errorf("expected ErrorsCount > 0")
	}
	if len(stats.Errors) != 0 {
		t.Errorf("expected 0 stored errors when IncludeErrors=false, got %d", len(stats.Errors))
	}
}

func TestStopOnError_StopsImmediately(t *testing.T) {
	root := setupTestData(t)

	opts := scanner.WalkOptions{
		IncludeStats:  true,
		IncludeErrors: true,
		StopOnError:   true,
		MaxDepth:      -1,
	}
	w := scanner.NewWalker(context.Background(), opts)

	called := 0
	w.OnVisitFile = func(path string, size int64) error {
		called++

		t.Logf("Visiting #%d: %s", called, filepath.Base(path))
		return errors.New("stop now")
	}

	err := w.Walk(root)
	if err == nil || err.Error() != "stop now" {
		t.Fatalf("expected stop error, got %v", err)
	}

	t.Logf("Walk stopped after %d calls", called)

	if called != 1 {
		t.Errorf("expected to stop after first file, got %d calls", called)
	}
}

func TestSkipOnError_Skips(t *testing.T) {
	root := setupTestData(t)

	opts := scanner.WalkOptions{
		IncludeStats:  true,
		IncludeErrors: true,
		StopOnError:   false,
		SkipOnError:   true,
		MaxDepth:      -1,
	}
	w := scanner.NewWalker(context.Background(), opts)

	called := 0
	w.OnVisitFile = func(path string, size int64) error {
		called++

		t.Logf("Visiting #%d: %s", called, filepath.Base(path))
		return errors.New("skip this")
	}

	err := w.Walk(root)
	if err != nil {
		t.Fatalf("unexpected walk error: %v", err)
	}

	stats := w.GetStats()
	t.Logf("Collected stats: %+v", stats)

	if called == 0 {
		t.Errorf("expected files to be visited")
	}

	if stats.ErrorsCount == 0 {
		t.Errorf("expected error count increment")
	}
}
*/
