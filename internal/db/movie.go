package db

import (
	"database/sql"
	"time"
)

// Movie represents a movie entry
type Movie struct {
	ID         int64
	Title      string
	Year       int
	MediaID    int64
	Source     string
	Group      string
	Collection string
	Thumbnail  *int64

	CreatedAt time.Time
}

// UpsertMovie inserts or updates a movie based on the UNIQUE constraint
func UpsertMovie(db *sql.DB, m Movie) (int64, error) {
	now := time.Now()
	res, err := db.Exec(`
		INSERT INTO movie (title, year, media_id, source, "group", "collection", created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(title, year, source, "group", "collection")
		DO UPDATE SET media_id = excluded.media_id, created_at = excluded.created_at
	`, m.Title, m.Year, m.MediaID, m.Source, m.Group, m.Collection, now)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil || id == 0 {
		// fallback: fetch existing ID in case of conflict
		row := db.QueryRow(`
			SELECT id FROM movie
			WHERE title=? AND year=? AND source=? AND "group"=? AND "collection"=?
		`, m.Title, m.Year, m.Source, m.Group, m.Collection)
		err = row.Scan(&id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

// UpdateMovieThumbnail sets the thumbnail for a movie
func UpdateMovieThumbnail(db *sql.DB, movieID, fileID int64) error {
	_, err := db.Exec(`UPDATE movie SET thumbnail_id = ? WHERE id = ?`, fileID, movieID)
	return err
}

// UpsertMovieSubtitle links a subtitle to a movie
func UpsertMovieSubtitle(db *sql.DB, movieID int64, fileID int64, language string) error {
	_, err := db.Exec(`
		INSERT INTO subtitle (media_type, media_id, subtitle_id, language)
		VALUES ('movie', ?, ?, ?)
		ON CONFLICT(media_type, media_id, subtitle_id) DO NOTHING
	`, movieID, fileID, language)
	return err
}

type MovieWithExtras struct {
	Movie     Movie
	Thumbnail *MediaThumbnail
	Subtitles []MediaSubtitle
}

// BatchUpsertMovies inserts or updates movies along with their thumbnails and subtitles
func BatchUpsertMovies(db *sql.DB, batch []MovieWithExtras) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Movie upsert statement
	movieStmt, err := tx.Prepare(`
		INSERT INTO movie (title, year, media_id, source, "group", "collection", created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(title, year, source, "group", "collection")
		DO UPDATE SET media_id = excluded.media_id
		RETURNING id
	`)
	if err != nil {
		return err
	}
	defer movieStmt.Close()

	// Thumbnail upsert
	thumbStmt, err := tx.Prepare(`
		INSERT INTO media_thumbnail (icon_id, desktop_ini_id)
		VALUES (?, ?)
		RETURNING id
	`)
	if err != nil {
		return err
	}
	defer thumbStmt.Close()

	// Subtitles upsert
	subStmt, err := tx.Prepare(`
		INSERT INTO media_subtitle (media_type, media_id, subtitle_id, language)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(media_type, media_id, subtitle_id) DO NOTHING
	`)
	if err != nil {
		return err
	}
	defer subStmt.Close()

	// Update parent reference in files
	updateFileStmt, err := tx.Prepare(`
		UPDATE files
		SET parent_type = ?, parent_id = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer updateFileStmt.Close()

	now := time.Now()

	for _, b := range batch {
		var movieID int64
		if err := movieStmt.QueryRow(
			b.Movie.Title, b.Movie.Year, b.Movie.MediaID, b.Movie.Source,
			b.Movie.Group, b.Movie.Collection, now,
		).Scan(&movieID); err != nil {
			return err
		}

		// Thumbnail
		if b.Thumbnail != nil {
			var thumbID int64
			if err := thumbStmt.QueryRow(b.Thumbnail.IconID, b.Thumbnail.DesktopIniID).Scan(&thumbID); err != nil {
				return err
			}

			if _, err := tx.Exec(`UPDATE movie SET thumbnail_id = ? WHERE id = ?`, thumbID, movieID); err != nil {
				return err
			}
		}

		// Subtitles + link to parent
		for _, s := range b.Subtitles {
			if _, err := subStmt.Exec(s.MediaType, movieID, s.SubtitleID, s.Language); err != nil {
				return err
			}

			if _, err := updateFileStmt.Exec("movie", movieID, s.SubtitleID); err != nil {
				return err
			}
		}

		// Main movie file → link parent
		if _, err := updateFileStmt.Exec("movie", movieID, b.Movie.MediaID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
