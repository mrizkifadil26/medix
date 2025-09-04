package db

import "database/sql"

// TV represents a TV show entry
type TV struct {
	ID     int64
	Title  string
	Source string
	Group  string

	Thumbnail *int64
}

// TVSeason represents a season of a TV show
type TVSeason struct {
	ID           int64
	TVID         int64
	SeasonNumber int
	EpisodeCount int
	FirstAirDate string
	LastAirDate  string
	TMDBID       *int64
}

// TVEpisode represents a single episode
type TVEpisode struct {
	ID            int64
	SeasonID      int64
	MediaID       int64 // file reference
	EpisodeNumber int
	TMDBID        *int64
}

type TVWithExtras struct {
	TV         TV
	Seasons    []TVSeason
	Episodes   []TVEpisode
	Thumbnails []MediaThumbnail
	Subtitles  []MediaSubtitle
}

// ---------------- TV ----------------
func UpsertTV(db *sql.DB, t TV) (int64, error) {
	res, err := db.Exec(`
		INSERT INTO tv (title, source, "group", thumbnail_id, created_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(title, source, "group")
		DO UPDATE SET thumbnail_id = excluded.thumbnail_id
	`, t.Title, t.Source, t.Group, t.Thumbnail)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// ---------------- TV Season ----------------
func UpsertTVSeason(db *sql.DB, s TVSeason) (int64, error) {
	res, err := db.Exec(`
		INSERT INTO tv_season (tv_id, season_number, episode_count, first_air_date, last_air_date, tmdb_id)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(tv_id, season_number) DO UPDATE SET
			episode_count = excluded.episode_count,
			first_air_date = excluded.first_air_date,
			last_air_date = excluded.last_air_date,
			tmdb_id = excluded.tmdb_id
	`, s.TVID, s.SeasonNumber, s.EpisodeCount, s.FirstAirDate, s.LastAirDate, s.TMDBID)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// ---------------- TV Episode ----------------
func UpsertTVEpisode(db *sql.DB, e TVEpisode) (int64, error) {
	res, err := db.Exec(`
		INSERT INTO tv_episode (season_id, media_id, episode_number, tmdb_id)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(season_id, episode_number) DO UPDATE SET
			media_id = excluded.media_id,
			tmdb_id = excluded.tmdb_id
	`, e.SeasonID, e.MediaID, e.EpisodeNumber, e.TMDBID)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// BatchUpsertTVFull performs a full TV batch insert/update
func BatchUpsertTVWithExtras(db *sql.DB, batch TVWithExtras) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1️⃣ Upsert TV
	var tvID int64
	res, err := tx.Exec(`
		INSERT INTO tv (title, source, "group", thumbnail_id, created_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(title, source, "group")
		DO UPDATE SET thumbnail_id = excluded.thumbnail_id
	`, batch.TV.Title, batch.TV.Source, batch.TV.Group, batch.TV.Thumbnail)
	if err != nil {
		return err
	}
	tvID, _ = res.LastInsertId()

	// 2️⃣ Upsert Seasons
	stmtSeason, err := tx.Prepare(`
		INSERT INTO tv_season (tv_id, season_number, episode_count, first_air_date, last_air_date, tmdb_id)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(tv_id, season_number) DO UPDATE SET
			episode_count = excluded.episode_count,
			first_air_date = excluded.first_air_date,
			last_air_date = excluded.last_air_date,
			tmdb_id = excluded.tmdb_id
	`)
	if err != nil {
		return err
	}
	defer stmtSeason.Close()

	for _, s := range batch.Seasons {
		if _, err := stmtSeason.Exec(tvID, s.SeasonNumber, s.EpisodeCount, s.FirstAirDate, s.LastAirDate, s.TMDBID); err != nil {
			return err
		}
	}

	// 3️⃣ Upsert Episodes
	stmtEpisode, err := tx.Prepare(`
		INSERT INTO tv_episode (season_id, media_id, episode_number, tmdb_id)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(season_id, episode_number) DO UPDATE SET
			media_id = excluded.media_id,
			tmdb_id = excluded.tmdb_id
	`)
	if err != nil {
		return err
	}
	defer stmtEpisode.Close()

	for _, e := range batch.Episodes {
		if _, err := stmtEpisode.Exec(e.SeasonID, e.MediaID, e.EpisodeNumber, e.TMDBID); err != nil {
			return err
		}
	}

	// 4️⃣ Upsert Thumbnails
	stmtThumbInsert, err := tx.Prepare(`INSERT INTO media_thumbnail (icon_id, desktop_ini_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer stmtThumbInsert.Close()

	stmtThumbUpdate, err := tx.Prepare(`UPDATE tv SET thumbnail_id = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmtThumbUpdate.Close()

	for _, t := range batch.Thumbnails {
		res, err := stmtThumbInsert.Exec(t.IconID, t.DesktopIniID)
		if err != nil {
			return err
		}
		thumbID, _ := res.LastInsertId()
		if _, err := stmtThumbUpdate.Exec(thumbID, tvID); err != nil {
			return err
		}
	}

	// 5️⃣ Upsert Subtitles
	stmtSub, err := tx.Prepare(`
		INSERT INTO media_subtitle (media_type, media_id, subtitle_id, language)
		VALUES ('tv', ?, ?, ?)
		ON CONFLICT(media_type, media_id, subtitle_id) DO NOTHING
	`)
	if err != nil {
		return err
	}
	defer stmtSub.Close()

	for _, s := range batch.Subtitles {
		if _, err := stmtSub.Exec(tvID, s.SubtitleID, s.Language); err != nil {
			return err
		}
	}

	return tx.Commit()
}
