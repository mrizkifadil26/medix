package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations applies all SQL files in migrations dir (sorted by filename).
func RunMigrations(db *sql.DB, dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("failed to read migrations dir: %v", err)
	}

	var sqlFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles) // apply in order: 000001, 000002...

	for _, fname := range sqlFiles {
		path := filepath.Join(dir, fname)
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read migration %s: %v", fname, err)
		}

		log.Printf("applying migration: %s", fname)
		if _, err := db.Exec(string(sqlBytes)); err != nil {
			log.Fatalf("failed to execute migration %s: %v", fname, err)
		}
	}
}
