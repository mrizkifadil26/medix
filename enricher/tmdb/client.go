package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func (c *Client) doRequest(
	endpoint string,
	params url.Values,
	result any,
) error {
	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	if params != nil {
		reqURL.RawQuery = params.Encode()
	}

	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	return nil
}
