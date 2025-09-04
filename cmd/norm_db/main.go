package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mrizkifadil26/medix/normdb"
)

func main() {
	// open SQLite connection
	dsn := "db/sqlite/media.db"
	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	// group movie files
	movieGroups, _, err := normdb.GroupFiles(conn, "movie")
	if err != nil {
		log.Fatalf("grouping movies failed: %v", err)
	}

	// normalize movies into DB
	var groups []normdb.MovieGroup
	for _, g := range movieGroups {
		groups = append(groups, *g)
	}

	if err := normdb.NormalizeMovies(conn, groups, "movie"); err != nil {
		log.Fatalf("normalize movies failed: %v", err)
	}

	fmt.Printf("✅ Normalized %d movies into DB\n", len(groups))

	// group TV files
	// _, shows, err := normdb.GroupFiles(conn, "tv")
	// if err != nil {
	// 	log.Fatalf("grouping TV files failed: %v", err)
	// }

	// fmt.Printf("📺 Grouped %d shows (TV not yet normalized)\n", len(shows))

	// exit cleanly
	os.Exit(0)
}
