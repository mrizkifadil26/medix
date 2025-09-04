package db

import "database/sql"

// MediaSubtitle represents the many-to-many subtitle relation
type MediaSubtitle struct {
	ID         int64
	MediaType  string
	MediaID    int64
	SubtitleID int64
	Language   string
}

// ---------------- MediaSubtitle ----------------
func UpsertMediaSubtitle(db *sql.DB, s MediaSubtitle) (int64, error) {
	res, err := db.Exec(`
		INSERT INTO media_subtitle (media_type, media_id, subtitle_id, language)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(media_type, media_id, subtitle_id) DO NOTHING
	`, s.MediaType, s.MediaID, s.SubtitleID, s.Language)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}
