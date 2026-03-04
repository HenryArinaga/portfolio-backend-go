package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "blog.db"
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Add new columns if they don't exist
	migrations := []string{
		`ALTER TABLE posts ADD COLUMN excerpt TEXT`,
		`ALTER TABLE posts ADD COLUMN image_url TEXT`,
	}

	for _, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			// Ignore errors if column already exists
			log.Printf("Migration note: %v", err)
		}
	}

	log.Println("Database migration completed successfully!")
}
