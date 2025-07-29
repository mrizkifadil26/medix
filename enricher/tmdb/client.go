package tmdb

import (
	"net/http"
	"time"
)

type bearerTransport struct {
	Token     string
	Transport http.RoundTripper
}

func (bt *bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+bt.Token)
	req.Header.Set("Accept", "application/json")

	return bt.Transport.RoundTrip(req)
}

type Client struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: "https://api.themoviedb.org/3",
		APIKey:  apiKey,
		Client: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &bearerTransport{
				Token:     apiKey,
				Transport: http.DefaultTransport,
			},
		},
	}
}

// func (e *TMDBClient) Run(
// 	ctx context.Context,
// 	job enricher.Config,
// ) (enricher.EnrichResult, error) {
// 	start := time.Now()

// 	if e.apiKey == "" {
// 		return enricher.EnrichResult{}, fmt.Errorf("TMDB_API_KEY not set")
// 	}

// 	raw, err := os.ReadFile(job.InputPath)
// 	if err != nil {
// 		return enricher.EnrichResult{}, fmt.Errorf("read error: %w", err)
// 	}

// 	var entries []model.MediaEntry
// 	if err := json.Unmarshal(raw, &entries); err != nil {
// 		return enricher.EnrichResult{}, fmt.Errorf("unmarshal error: %w", err)
// 	}

// 	var enrichedCount int
// 	for i := range entries {
// 		meta, err := e.fetch(entries[i].Title, entries[i].Year)
// 		if err != nil {
// 			continue // skip if not found or failed
// 		}

// 		entries[i].Metadata.Enriched = true
// 		entries[i].Metadata.Source = "TMDB"
// 		entries[i].Metadata.Overview = meta.Overview
// 		entries[i].Metadata.PosterPath = meta.PosterPath
// 		entries[i].Metadata.ReleaseDate = meta.ReleaseDate
// 		enrichedCount++
// 	}

// 	out, _ := json.MarshalIndent(entries, "", "  ")
// 	if err := os.WriteFile(job.Output, out, 0644); err != nil {
// 		return enricher.EnrichResult{}, fmt.Errorf("write error: %w", err)
// 	}

// 	return enricher.EnrichResult{
// 		Type:     "media",
// 		JobName:  job.Name,
// 		Items:    entries,
// 		Count:    enrichedCount,
// 		Duration: time.Since(start).String(),
// 	}, nil
// }

// func (e *TMDBEnricher) fetch(
// 	title string,
// 	year int,
// ) (*TMDBSearchItem, error) {
// 	baseURL := "https://api.themoviedb.org/3/search/movie"
// 	params := url.Values{}
// 	params.Set("api_key", e.apiKey)
// 	params.Set("query", title)
// 	params.Set("year", fmt.Sprintf("%d", year))

// 	resp, err := e.client.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("TMDB error: status %d", resp.StatusCode)
// 	}

// 	var result TMDBSearchResult
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return nil, err
// 	}
// 	if len(result.Results) == 0 {
// 		return nil, fmt.Errorf("no result")
// 	}

// 	return &result.Results[0], nil
// }
