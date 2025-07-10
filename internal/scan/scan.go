package scan

const maxConcurrency = 8

func ScanAll[T any](cfg ScanConfig, strategy ScanStrategy[T]) T {
	return strategy.Scan(cfg.Sources)
}

// result := model.RawOutput{
// 	Type:        "raw",
// 	GeneratedAt: time.Now(),
// }

// cache := &dirCache{}
// genreMap := make(map[string]model.RawGenre)
// for _, root := range cfg.Sources {
// 	scanSingleRoot(cfg.ContentType, root, genreMap, cache)
// }

// // Sort and reassemble
// var genreNames []string
// for name := range genreMap {
// 	genreNames = append(genreNames, name)
// }
// sort.Strings(genreNames)

// for _, name := range genreNames {
// 	result.Data = append(result.Data, genreMap[name])
// }

// fmt.Println("All sources scanned.")
// return result
// }

// func scanSingleRoot(contentType, root string, genreMap map[string]model.RawGenre, cache *dirCache) {
// 	entries, err := os.ReadDir(root)
// 	if err != nil {
// 		fmt.Printf("Failed to read root directory: %v\n", err)
// 		return
// 	}

// 	var genres []os.DirEntry
// 	for _, g := range entries {
// 		if g.IsDir() {
// 			genres = append(genres, g)
// 		}
// 	}

// 	fmt.Printf("Scanning %d genres in %s\n", len(genres), root)
// 	bar := progressbar.NewOptions(len(genres),
// 		progressbar.OptionSetDescription("Scanning genres"),
// 		progressbar.OptionSetWidth(30),
// 		progressbar.OptionShowCount(),
// 		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "#", SaucerPadding: " ", BarStart: "[", BarEnd: "]"}),
// 	)

// 	var mu sync.Mutex
// 	var wg sync.WaitGroup
// 	sem := make(chan struct{}, maxConcurrency)
// 	for _, genre := range genres {
// 		wg.Add(1)
// 		sem <- struct{}{}

// 		go func(genre os.DirEntry) {
// 			defer wg.Done()
// 			defer func() { <-sem }()
// 			defer bar.Add(1)

// 			genreName := genre.Name()
// 			genrePath := filepath.Join(root, genreName)

// 			items := scanGenre(genrePath, cache, contentType)
// 			if len(items) == 0 {
// 				return
// 			}

// 			mu.Lock()
// 			existing := genreMap[genreName]
// 			existing.Name = genreName
// 			existing.Items = append(existing.Items, items...)
// 			genreMap[genreName] = existing
// 			mu.Unlock()
// 		}(genre)
// 	}

// 	wg.Wait()
// }

// func scanGenre(genrePath string, cache *dirCache, contentType string) []model.RawEntry {
// 	entries := cache.Read(genrePath)
// 	if entries == nil {
// 		return nil
// 	}

// 	// Sort title folders
// 	sort.Slice(entries, func(i, j int) bool {
// 		return entries[i].Name() < entries[j].Name()
// 	})

// 	var items []model.RawEntry

// 	for _, entry := range entries {
// 		if !entry.IsDir() {
// 			continue
// 		}

// 		titlePath := filepath.Join(genrePath, entry.Name())
// 		subEntries := cache.Read(titlePath)
// 		if subEntries == nil {
// 			continue
// 		}

// 		status := resolveStatus(subEntries)
// 		ico := findIcon(titlePath, subEntries)

// 		list := extractChildren(titlePath, subEntries, cache, contentType)
// 		itemType := "single"
// 		if contentType == "movies" && len(list) > 0 {
// 			itemType = "collection"
// 		}

// 		rawEntry := model.RawEntry{
// 			Type:   itemType,
// 			Name:   entry.Name(),
// 			Path:   titlePath,
// 			Status: status,
// 			Icon:   ico,
// 		}

// 		switch contentType {
// 		case "movies":
// 			if len(list) > 0 {
// 				rawEntry.Items = &model.RawEntryItems{Entries: list}
// 			}
// 		case "tvshows":
// 			seasons := extractSeasonNames(subEntries)
// 			if len(seasons) > 0 {
// 				rawEntry.Items = &model.RawEntryItems{Seasons: seasons}
// 			}
// 		}

// 		items = append(items, rawEntry)
// 	}

// 	return items
// }
