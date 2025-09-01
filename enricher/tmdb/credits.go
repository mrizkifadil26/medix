package tmdb

import (
	"fmt"
)

// CreditService handles fetching movie credits from TMDb
type CreditService struct {
	client *Client
	cache  *creditsCache
}

// NewCreditService creates a new CreditService
func NewCreditService(client *Client, cache *creditsCache) *CreditService {
	return &CreditService{client: client, cache: cache}
}

// FetchCredits fetches actors, directors, and producers for a given TMDb movie ID
func (s *CreditService) FetchCredits(tmdbID int) (TMDbCredits, error) {
	if tmdbID == 0 {
		return TMDbCredits{}, fmt.Errorf("invalid TMDbID")
	}

	// --- Check cache first ---
	key := fmt.Sprint(tmdbID) // convert TMDbID to string for cache
	if cached, ok := s.cache.Get("credits", key); ok {
		return cached, nil
	}

	creditsRaw, err := s.client.GetMovieCredits(tmdbID)
	if err != nil {
		return TMDbCredits{}, fmt.Errorf("failed to fetch credits: %w", err)
	}

	credits := TMDbCredits{}

	// --- Actors ---
	for _, c := range creditsRaw.Cast {
		credits.Actors = append(credits.Actors, Person{
			ID:   c.ID,
			Name: c.Name,
			Role: c.Character, // character name for actor
		})
	}

	// --- Directors & Producers ---
	for _, crew := range creditsRaw.Crew {
		switch crew.Job {
		case "Director":
			credits.Directors = append(credits.Directors, Person{
				ID:   crew.ID,
				Name: crew.Name,
				Role: "director",
			})
		case "Producer":
			credits.Producers = append(credits.Producers, Person{
				ID:   crew.ID,
				Name: crew.Name,
				Role: "producer",
			})
		}
	}

	// --- Store in cache ---
	s.cache.Put("credits", key, credits)
	_ = s.cache.Save()

	return credits, nil
}

func (c *Client) GetMovieCredits(tmdbID int) (*CreditsRaw, error) {
	if tmdbID == 0 {
		return nil, fmt.Errorf("invalid TMDbID")
	}

	endpoint := fmt.Sprintf("%s/movie/%d/credits", c.BaseURL, tmdbID)

	var result CreditsRaw
	if err := c.doRequest(endpoint, nil, &result); err != nil {
		return nil, fmt.Errorf("failed to fetch credits: %w", err)
	}

	return &result, nil
}

// CreditsRaw represents the TMDb /credits response
type CreditsRaw struct {
	Cast []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Character string `json:"character"`
	} `json:"cast"`

	Crew []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Job  string `json:"job"`
	} `json:"crew"`
}
