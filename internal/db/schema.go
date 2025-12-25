package db

import "database/sql"

func Init(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT NOT NULL UNIQUE,
		content TEXT NOT NULL,
		published BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL
	);
	`
	_, err := db.Exec(query)
	return err
}
