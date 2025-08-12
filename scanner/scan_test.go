package scanner_test

/*
func TestScan_Basic(t *testing.T) {
	root := t.TempDir()

	// Create test files
	files := []string{
		"Action/Inception (2010)/Inception.2010.mkv",
		"Comedy/DumbAndDumber.1994.mp4",
	}

	for _, relPath := range files {
		fullPath := filepath.Join(root, relPath)
		require.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0755))
		require.NoError(t, os.WriteFile(fullPath, []byte("dummy"), 0644))
	}

	opts := scannerV2.ScanOptions{
		Mode:  "files",
		Exts:  []string{".mkv"},
		Depth: -1, // Unlimited depth
	}

	out, err := scannerV2.Scan(root, opts)

	require.NoError(t, err)
	require.Len(t, out.Items, 1)

	entry := out.Items[0]
	require.Contains(t, entry.ItemPath, "Inception.2010.mkv")
	require.Equal(t, "Inception.2010.mkv", entry.ItemName)
	require.NotNil(t, entry.ItemSize)
}
	/*

func TestScan_FilterMP4(t *testing.T) {
	root := t.TempDir()

	files := []string{
		"Comedy/Movie1.mp4",
		"Comedy/Movie2.mkv", // Should be ignored
	}

	for _, f := range files {
		full := filepath.Join(root, f)
		_ = os.MkdirAll(filepath.Dir(full), 0755)
		_ = os.WriteFile(full, []byte("x"), 0644)
	}

	opts := scannerV2.ScanOptions{
		Mode:  "files",
		Exts:  []string{".mp4"},
		Depth: -1, // Unlimited depth
	}

	out, err := scannerV2.Scan(root, opts)
	require.NoError(t, err)
	require.Len(t, out.Items, 1)
	require.Equal(t, "Movie1.mp4", out.Items[0].ItemName)
}

func TestScan_DepthLimit(t *testing.T) {
	root := t.TempDir()

	paths := []string{
		"Level1/File1.mp4",        // Should be included (depth 1)
		"Level1/Level2/File2.mp4", // Should be excluded (depth 2)
	}

	for _, p := range paths {
		full := filepath.Join(root, p)
		_ = os.MkdirAll(filepath.Dir(full), 0755)
		_ = os.WriteFile(full, []byte("data"), 0644)
	}

	opts := scannerV2.ScanOptions{
		Mode:  "files",
		Exts:  []string{".mp4"},
		Depth: 1,
	}

	out, err := scannerV2.Scan(root, opts)
	require.NoError(t, err)
	require.Len(t, out.Items, 1)
	require.Equal(t, "File1.mp4", out.Items[0].ItemName)
}

func TestScan_EmptyFolder(t *testing.T) {
	root := t.TempDir()

	_ = os.MkdirAll(filepath.Join(root, "EmptyDir"), 0755)

	opts := scannerV2.ScanOptions{
		Mode:  "files",
		Exts:  []string{".mp4"},
		Depth: -1, // Unlimited depth
	}

	out, err := scannerV2.Scan(root, opts)
	require.NoError(t, err)
	require.Len(t, out.Items, 0)
}

func TestScan_ExcludePath(t *testing.T) {
	root := t.TempDir()

	files := []string{
		"Include/Movie.mp4",
		"Exclude/Secret.mp4",
	}

	for _, f := range files {
		full := filepath.Join(root, f)
		_ = os.MkdirAll(filepath.Dir(full), 0755)
		_ = os.WriteFile(full, []byte("data"), 0644)
	}

	opts := scannerV2.ScanOptions{
		Mode:    "files",
		Exts:    []string{".mp4"},
		Depth:   -1,
		Exclude: []string{"Exclude"},
	}

	out, err := scannerV2.Scan(root, opts)
	require.NoError(t, err)
	require.Len(t, out.Items, 1)
	require.Contains(t, out.Items[0].ItemPath, "Include/Movie.mp4")
}

func TestScan_ModeDirs(t *testing.T) {
	root := t.TempDir()

	_ = os.MkdirAll(filepath.Join(root, "Genre/Film1 (2020)"), 0755)
	_ = os.MkdirAll(filepath.Join(root, "Genre/Film2 (2021)"), 0755)

	opts := scannerV2.ScanOptions{
		Mode:  "dirs",
		Depth: 2,
	}

	out, err := scannerV2.Scan(root, opts)
	require.NoError(t, err)
	require.Len(t, out.Items, 4) // root + Genre + both Film dirs
}

func TestScan_MultiExtDeep(t *testing.T) {
	root := t.TempDir()

	files := []string{
		"Series/Season1/Episode1.mkv",
		"Series/Season1/Episode2.mp4",
		"Series/Season2/Episode3.avi", // not in filter
	}

	for _, p := range files {
		full := filepath.Join(root, p)
		_ = os.MkdirAll(filepath.Dir(full), 0755)
		_ = os.WriteFile(full, []byte("ep"), 0644)
	}

	opts := scannerV2.ScanOptions{
		Mode:  "files",
		Exts:  []string{".mkv", ".mp4"},
		Depth: 3,
	}

	out, err := scannerV2.Scan(root, opts)
	require.NoError(t, err)
	require.Len(t, out.Items, 2)
} */

/*
func TestScan_TVShowStyle(t *testing.T) {
	root := t.TempDir()

	// Simulate TV Show layout
	tvShow := filepath.Join(root, "Breaking Bad")
	seasons := []string{
		filepath.Join(tvShow, "Season 1"),
		filepath.Join(tvShow, "Season 2"),
	}

	// Create nested season folders
	for _, season := range seasons {
		require.NoError(t, os.MkdirAll(season, 0755))
		// Add dummy video in each season
		file := filepath.Join(season, "Episode1.mkv")
		require.NoError(t, os.WriteFile(file, []byte("dummy"), 0644))
	}

	opts := scannerV2.ScanOptions{
		Mode:      "dirs",
		Depth:     -1, // Only go one level down, scan "Breaking Bad", but not seasons
		LeafDepth: 1,  // Only take "Breaking Bad" and nest children
	}

	out, err := scannerV2.Scan(root, opts)
	require.NoError(t, err)
	require.Len(t, out.Items, 1)

	entry := out.Items[0]
	require.Equal(t, "Breaking Bad", entry.ItemName)
	require.GreaterOrEqual(t, len(entry.SubEntries), 2)

	subNames := append([]string{}, entry.SubEntries...)

	require.Contains(t, subNames, "Season 1")
	require.Contains(t, subNames, "Season 2")
}
*/

/*
func TestScan_ConcurrencyEffect(t *testing.T) {
	tmp := t.TempDir()

	for i := 0; i < 5; i++ {
		filePath := filepath.Join(tmp, fmt.Sprintf("movie%d.mkv", i))
		os.WriteFile(filePath, []byte("dummy"), 0644)
	}

	start := time.Now()
	_, err := scannerV2.Scan(tmp, scannerV2.ScanOptions{
		Mode:        "files",
		Depth:       1,
		Exts:        []string{".mkv"},
		Verbose:     false,
		Concurrency: 1, // sequential
	})
	require.NoError(t, err)
	durationSequential := time.Since(start)

	start = time.Now()
	_, err = scannerV2.Scan(tmp, scannerV2.ScanOptions{
		Mode:        "files",
		Depth:       1,
		Exts:        []string{".mkv"},
		Verbose:     false,
		Concurrency: 4, // parallel
	})

	require.NoError(t, err)
	durationParallel := time.Since(start)

	assert.Less(t, durationParallel, durationSequential)
}
*/
