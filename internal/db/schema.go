package db

import "database/sql"

func Init(db *sql.DB) error {
	query := `
	-- Blog posts
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT NOT NULL UNIQUE,
		content TEXT NOT NULL,
		published INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL
	);

	-- SCS session storage
	CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		data BLOB NOT NULL,
		expiry DATETIME NOT NULL
	);

	-- Helpful indexes
	CREATE INDEX IF NOT EXISTS idx_posts_published
		ON posts(published);

	CREATE INDEX IF NOT EXISTS idx_posts_created_at
		ON posts(created_at DESC);

	CREATE INDEX IF NOT EXISTS idx_sessions_expiry
		ON sessions(expiry);
	`
	_, err := db.Exec(query)
	return err
}
