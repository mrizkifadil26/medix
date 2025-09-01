package db

type Movie struct {
	ID     int
	Title  string
	Year   string
	TMDBID int
}

func (db *DB) InsertMovie(title, year string, tmdbID int) error {
	_, err := db.Exec(`
        INSERT INTO movie (title, year, tmdb_id)
        VALUES (?, ?, ?)
        ON CONFLICT(tmdb_id) DO NOTHING;
    `, title, year, tmdbID)

	return err
}

func (db *DB) GetMovieByTMDBID(tmdbID int) (*Movie, error) {
	row := db.QueryRow(`SELECT id, title, year, tmdb_id FROM movie WHERE tmdb_id = ?`, tmdbID)
	var m Movie
	if err := row.Scan(&m.ID, &m.Title, &m.Year, &m.TMDBID); err != nil {
		return nil, err
	}

	return &m, nil
}
