package syncer

type SyncConfig struct {
	Name       string `json:"name"`       // "movies" or "tvshows"
	MediaInput string `json:"mediaInput"` // Input raw media JSON
	IconInput  string `json:"iconInput"`  // Input icon index JSON
	OutMedia   string `json:"outMedia"`   // Synced media output JSON
	OutIcon    string `json:"outIcon"`    // Synced icon index output JSON
}
