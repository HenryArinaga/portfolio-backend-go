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

func CreatePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var title, content string
		var published bool

		ct := r.Header.Get("Content-Type")

		if strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "invalid form", http.StatusBadRequest)
				return
			}

			title = r.FormValue("title")
			content = r.FormValue("content")
			published = r.FormValue("published") == "on"
		} else {
			var req struct {
				Title     string `json:"title"`
				Content   string `json:"content"`
				Published bool   `json:"published"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			title = req.Title
			content = req.Content
			published = req.Published
		}

		if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
			http.Error(w, "title and content required", http.StatusBadRequest)
			return
		}

		slug := strings.ToLower(strings.TrimSpace(title))
		slug = strings.ReplaceAll(slug, " ", "-")

		_, err := db.Exec(`
			INSERT INTO posts (title, slug, content, published, created_at)
			VALUES (?, ?, ?, ?, ?)
		`,
			title,
			slug,
			content,
			boolToInt(published),
			time.Now(),
		)
		if err != nil {
			http.Error(w, "failed to insert post", http.StatusInternalServerError)
			return
		}

		if strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func UpdatePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Path[len("/api/admin/posts/"):]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid post id", http.StatusBadRequest)
			return
		}

		var title, content string
		var published bool
		ct := r.Header.Get("Content-Type")

		if strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "invalid form", http.StatusBadRequest)
				return
			}
			title = r.FormValue("title")
			content = r.FormValue("content")
			published = r.FormValue("published") == "on"
		} else {
			var req UpdatePostRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			title = req.Title
			content = req.Content
			published = req.Published
		}

		slug := strings.ToLower(strings.TrimSpace(title))
		slug = strings.ReplaceAll(slug, " ", "-")

		_, err = db.Exec(`
			UPDATE posts
			SET title = ?, slug = ?, content = ?, published = ?
			WHERE id = ?
		`, title, slug, content, boolToInt(published), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func ListPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
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

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func ListPostsData(db *sql.DB) ([]Post, error) {
	rows, err := db.Query(`
		SELECT id, title, slug, content, created_at, published
		FROM posts
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(
			&p.ID, &p.Title, &p.Slug,
			&p.Content, &p.CreatedAt, &p.Published,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetPostByID(db *sql.DB, id int64) (Post, error) {
	var p Post
	err := db.QueryRow(`
		SELECT id, title, slug, content, created_at, published
		FROM posts
		WHERE id = ?
	`, id).Scan(
		&p.ID, &p.Title, &p.Slug,
		&p.Content, &p.CreatedAt, &p.Published,
	)
	return p, err
}

func DeletePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := strings.TrimSuffix(
			strings.TrimPrefix(r.URL.Path, "/api/admin/posts/"),
			"/delete",
		)

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid post id", http.StatusBadRequest)
			return
		}

		if _, err := db.Exec(`DELETE FROM posts WHERE id = ?`, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}
