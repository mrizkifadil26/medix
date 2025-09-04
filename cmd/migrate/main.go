package main

import (
	"flag"
	"log"

	"github.com/mrizkifadil26/medix/internal/db"
)

func main() {
	dbPath := flag.String("db", "app.db", "SQLite database file")
	migrationsDir := flag.String("migrations", "migrations", "Migrations directory")
	flag.Parse()

	d := db.Open(*dbPath)
	defer d.Close()

	db.RunMigrations(d, *migrationsDir)

	log.Println("Migrations applied successfully!")
}
