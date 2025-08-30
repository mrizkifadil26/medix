package tmdb

import "github.com/mrizkifadil26/medix/utils/cache"

// internal cache managers, not exported
type dataCache = cache.Manager[*EnrichedItem]
type genreCache = cache.Manager[map[int]string]
type langCache = cache.Manager[map[string]string]

// create all caches internally
func newCaches() (*dataCache, *genreCache, *langCache) {
	return cache.NewManager[*EnrichedItem]("tmdb.data.cache.json"),
		cache.NewManager[map[int]string]("tmdb.genres.cache.json"),
		cache.NewManager[map[string]string]("tmdb.languages.cache.json")
}
