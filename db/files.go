package db

import "time"

type File struct {
	ID           int64      `db:"id"`
	AbsolutePath string     `db:"absolute_path"`
	RelativePath string     `db:"relative_path"`
	Filename     string     `db:"filename"`
	Extension    string     `db:"extension"`
	Type         string     `db:"type"` // video, subtitle, thumbnail
	Size         string     `db:"size"`
	ModTime      *time.Time `db:"modtime"`

	ParentType string `db:"parent_type"` // "movie" | "tv"
	ParentID   int64  `db:"parent_id"`

	CreatedAt *time.Time `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// Insert file
func (db *DB) InsertFile(f *File) (int64, error) {
	res, err := db.Exec(`
		INSERT INTO files (
			absolute_path, relative_path, filename, extension, type, size, modtime,
			parent_type, parent_id, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
	`, f.AbsolutePath, f.RelativePath, f.Filename, f.Extension, f.Type, f.Size, f.ModTime, f.ParentType, f.ParentID)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// Get by ID
func (db *DB) GetFile(id int64) (*File, error) {
	row := db.QueryRow(`SELECT * FROM files WHERE id = ? AND deleted_at IS NULL`, id)
	var f File

	if err := row.Scan(
		&f.ID, &f.AbsolutePath, &f.RelativePath, &f.Filename, &f.Extension,
		&f.Type, &f.Size, &f.ModTime, &f.ParentType, &f.ParentID,
		&f.CreatedAt, &f.DeletedAt,
	); err != nil {
		return nil, err
	}

	return &f, nil
}
