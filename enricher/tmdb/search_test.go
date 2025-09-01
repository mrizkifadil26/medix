package tmdb_test

/*
func TestSearchMovie_WithBearerToken(t *testing.T) {
	expectedToken := "test-token"

	tests := []struct {
		name           string
		query          tmdb.SearchQuery
		expectError    bool
		mockStatusCode int
	}{
		{
			name: "Valid query with year",
			query: tmdb.SearchQuery{
				Query: "Doraemon",
				Year:  "2024",
			},
			expectError:    false,
			mockStatusCode: http.StatusOK,
		},
		{
			name: "Valid query with primary_release_year",
			query: tmdb.SearchQuery{
				Query:       "Doraemon",
				PrimaryYear: "2024",
			},
			expectError:    false,
			mockStatusCode: http.StatusOK,
		},
		{
			name: "Missing query (should fail)",
			query: tmdb.SearchQuery{
				Year: "2024",
			},
			expectError:    true,
			mockStatusCode: http.StatusOK,
		},
		{
			name: "API returns 401 (unauthorized)",
			query: tmdb.SearchQuery{
				Query: "Doraemon",
			},
			expectError:    true,
			mockStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		tc := tc // pin
		t.Run(tc.name, func(t *testing.T) {
			mockResults := tmdb.SearchResult{
				Results: []tmdb.SearchItem{
					{Title: "Doraemon: Nobita's Sky Utopia", ReleaseDate: "2024-03-01"},
				},
			}

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Assert Bearer token
				auth := r.Header.Get("Authorization")
				if auth != "Bearer "+expectedToken {
					t.Errorf("expected Authorization: Bearer %s, got %s", expectedToken, auth)
				}

				// Mock status code
				w.WriteHeader(tc.mockStatusCode)

				if tc.mockStatusCode == http.StatusOK {
					_ = json.NewEncoder(w).Encode(mockResults)
				} else {
					_ = json.NewEncoder(w).Encode(map[string]any{"status_message": "Unauthorized"})
				}
			}))
			defer ts.Close()

			client := tmdb.NewClient(expectedToken)
			client.BaseURL = ts.URL

			results, err := client.Search("movie", tc.query)

			if tc.expectError && err == nil {
				t.Fatalf("expected error, got none")
			}
			if !tc.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.expectError && len(results) != 1 {
				t.Errorf("expected 1 result, got %d", len(results))
			}
		})
	}
}

func TestSearchQuery_Validation(t *testing.T) {
	tests := []struct {
		name        string
		query       tmdb.SearchQuery
		expectedErr error
	}{
		{
			name:        "Empty query",
			query:       tmdb.SearchQuery{},
			expectedErr: errors.New("field Query is required"),
		},
		{
			name:  "Only query",
			query: tmdb.SearchQuery{Query: "Tenet"},
		},
		{
			name:  "Query with year",
			query: tmdb.SearchQuery{Query: "Tenet", Year: "2020"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.query.Validate()
			if tc.expectedErr != nil && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if tc.expectedErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.expectedErr != nil && !strings.Contains(err.Error(), tc.expectedErr.Error()) {
				t.Fatalf("expected error '%s', got '%s'", tc.expectedErr.Error(), err.Error())
			}
		})
	}
}

func TestSearchMovie_RealAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skip real API test in short mode")
	}

	const token = "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiI3N2QxNGJiN2JkODYyY2E0ZTE4MzBiZWNiODgxNGU3NyIsIm5iZiI6MTU2MzYzMjEyOC43NzUsInN1YiI6IjVkMzMyMjAwYWU2ZjA5MDAwZTdiNWJlZiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.ZpKTx-psLaWjwS-zRvpDLO7QKmoNJnF_xubbb8vn-48"

	client := tmdb.NewClient(token)

	items, err := client.SearchMovie("Due Date", 2010)
	if err != nil {
		t.Fatalf("real TMDB search failed: %v", err)
	}
	if len(items) == 0 {
		t.Errorf("no results returned from real TMDB search")
	} else {
		t.Logf("✅ Got %d results. First title: %s", len(items), items[0].Title)
	}
}
*/
