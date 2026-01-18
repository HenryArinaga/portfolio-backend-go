// internal/api/admin/posts.go
package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Published bool      `json:"published"`
}

type UpdatePostRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Published bool   `json:"published"`
}

type createPostRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Published bool   `json:"published"`
}

func CreatePost(db *sql.DB, adminToken string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// only allow POST
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// admin auth check
		auth := r.Header.Get("Authorization")
		if adminToken == "" || auth != "Bearer "+adminToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// parse request body
		var req createPostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Content) == "" {
			http.Error(w, "title and content required", http.StatusBadRequest)
			return
		}

		// generate slug
		slug := strings.ToLower(req.Title)
		slug = strings.TrimSpace(slug)
		slug = strings.ReplaceAll(slug, " ", "-")

		// insert post
		_, err := db.Exec(`
			INSERT INTO posts (title, slug, content, published, created_at)
			VALUES (?, ?, ?, ?, ?)
		`,
			req.Title,
			slug,
			req.Content,
			boolToInt(req.Published),
			time.Now(),
		)

		if err != nil {
			http.Error(w, "failed to insert post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func UpdatePost(db *sql.DB, adminToken string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		auth := r.Header.Get("Authorization")
		if adminToken == "" || auth != "Bearer "+adminToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		idStr := r.URL.Path[len("/api/admin/posts/"):]
		if idStr == "" {
			http.NotFound(w, r)
			return
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid post id", http.StatusBadRequest)
			return
		}

		var req UpdatePostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		slug := strings.ToLower(strings.TrimSpace(req.Title))
		slug = strings.ReplaceAll(slug, " ", "-")

		_, err = db.Exec(`
			UPDATE posts
			SET title = ?, slug = ?, content = ?, published = ?
			WHERE id = ?
		`, req.Title, slug, req.Content, boolToInt(req.Published), id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func ListPosts(db *sql.DB, adminToken string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		auth := r.Header.Get("Authorization")
		if adminToken == "" || auth != "Bearer "+adminToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
			SELECT id, title, slug, content, created_at, published
			FROM posts
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
}
