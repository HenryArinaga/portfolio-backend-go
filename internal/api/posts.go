// internal/api/posts.go
package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

var database *sql.DB

func SetDB(db *sql.DB) {
	database = db
}

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Published bool      `json:"published"`
}

func ListPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query(`
		SELECT id, title, slug, content, created_at, published
		FROM posts
		WHERE published = 1
		ORDER BY created_at DESC
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []Post{}

	for rows.Next() {
		var p Post
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.Content,
			&p.CreatedAt,
			&p.Published,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetPostBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[len("/api/posts/"):]
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	var p Post
	err := database.QueryRow(`
		SELECT id, title, slug, content, created_at, published
		FROM posts
		WHERE slug = ? AND published = 1
	`, slug).Scan(
		&p.ID,
		&p.Title,
		&p.Slug,
		&p.Content,
		&p.CreatedAt,
		&p.Published,
	)

	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
