package db

import (
	"database/sql"
	"time"
)

type File struct {
	ID           int64     `db:"id"`
	AbsolutePath string    `db:"absolute_path"`
	RelativePath string    `db:"relative_path"`
	Filename     string    `db:"filename"`
	Extension    string    `db:"extension"`
	Type         string    `db:"type"` // video, subtitle, thumbnail
	Size         int64     `db:"size"`
	ModTime      time.Time `db:"modtime"`

	ParentType  string `db:"parent_type"` // "movie" | "tv"
	ParentID    *int64 `db:"parent_id"`
	Fingerprint string `db:"fingerprint"`

	CreatedAt time.Time  `db:"created_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// UpsertFile inserts or updates a file record based on absolute_path.
func UpsertFile(db *sql.DB, f *File) error {
	now := time.Now()

	_, err := db.Exec(`
		INSERT INTO files (
			absolute_path, relative_path, filename, extension, type, size, modtime,
			parent_type, parent_id, fingerprint, created_at, deleted_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(fingerprint) DO UPDATE SET
			absolute_path=excluded.absolute_path,
			relative_path=excluded.relative_path,
			filename=excluded.filename,
			extension=excluded.extension,
			type=excluded.type,
			size=excluded.size,
			modtime=excluded.modtime,
			parent_type=excluded.parent_type,
			parent_id=excluded.parent_id,
			deleted_at=excluded.deleted_at
	`, f.AbsolutePath, f.RelativePath, f.Filename, f.Extension, f.Type, f.Size,
		f.ModTime, f.ParentType, f.ParentID, f.Fingerprint, now, f.DeletedAt)

	return err
}

// BatchUpsertFiles inserts or updates multiple files in a single transaction.
func BatchUpsertFiles(db *sql.DB, files []File) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO files (
			absolute_path, relative_path, filename, extension, type, size, modtime,
			parent_type, parent_id, fingerprint, created_at, deleted_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(fingerprint) DO UPDATE SET
			absolute_path=excluded.absolute_path,
			relative_path=excluded.relative_path,
			filename=excluded.filename,
			extension=excluded.extension,
			type=excluded.type,
			size=excluded.size,
			modtime=excluded.modtime,
			parent_type=excluded.parent_type,
			parent_id=excluded.parent_id,
			deleted_at=excluded.deleted_at
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()
	for _, f := range files {
		_, err := stmt.Exec(
			f.AbsolutePath,
			f.RelativePath,
			f.Filename,
			f.Extension,
			f.Type,
			f.Size,
			f.ModTime,
			f.ParentType,
			f.ParentID,
			f.Fingerprint,
			now,
			f.DeletedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// ListFiles returns all files from the `files` table
func ListFiles(db *sql.DB) ([]File, error) {
	rows, err := db.Query(`
		SELECT id, absolute_path, relative_path, filename, extension, type, size, modtime, parent_type, parent_id, fingerprint
		FROM files
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var f File
		var parentID sql.NullInt64
		if err := rows.Scan(
			&f.ID,
			&f.AbsolutePath,
			&f.RelativePath,
			&f.Filename,
			&f.Extension,
			&f.Type,
			&f.Size,
			&f.ModTime,
			&f.ParentType,
			&parentID,
			&f.Fingerprint,
		); err != nil {
			return nil, err
		}

		if parentID.Valid {
			f.ParentID = &parentID.Int64
		}

		files = append(files, f)
	}

	return files, nil
}
