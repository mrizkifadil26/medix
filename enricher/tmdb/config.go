package tmdb

type Config struct {
	APIKey               string `json:"api_key"`
	FetchCredits         bool   `json:"fetch_credits,omitempty"`         // optional, default false
	FetchRecommendations bool   `json:"fetch_recommendations,omitempty"` // optional, default false
	FetchSimilar         bool   `json:"fetch_similar,omitempty"`         // optional, default false
}
