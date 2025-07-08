package util

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Me, Myself & Irene.ico", "me-myself-and-irene"},
		{"Alien (1979).ico", "alien-1979"},
		{"She's Out of My League (Alt).ico", "shes-out-of-my-league-alt"},
		{"Romeo + Juliet.ico", "romeo-juliet"},
		{"Alien³.ico", "alien3"},
		{"9½ Weeks.ico", "91-2-weeks"},
		{"You've Got Mail.ico", "youve-got-mail"},
		{"Zack Snyder’s Justice League.ico", "zack-snyders-justice-league"},
		{"Café Society.ico", "cafe-society"},
		{"The Devil's Advocate.ico", "the-devils-advocate"},
	}

	for _, tt := range tests {
		got := Slugify(tt.input)
		if got != tt.expected {
			t.Errorf("Slugify(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}
