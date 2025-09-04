package db

import "database/sql"

// MediaThumbnail represents a thumbnail
type MediaThumbnail struct {
	ID           int64
	IconID       *int64
	DesktopIniID *int64
}

// ---------------- MediaThumbnail ----------------
func UpsertMediaThumbnail(db *sql.DB, t MediaThumbnail) (int64, error) {
	res, err := db.Exec(`
		INSERT INTO media_thumbnail (icon_id, desktop_ini_id)
		VALUES (?, ?)
	`, t.IconID, t.DesktopIniID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
