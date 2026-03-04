package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/henryarin/portfolio-backend-go/internal/storage"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	Excerpt   string    `json:"excerpt"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	Published bool      `json:"published"`
}

type UpdatePostRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Excerpt   string `json:"excerpt"`
	ImageURL  string `json:"image_url"`
	Published bool   `json:"published"`
}

type createPostRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Excerpt   string `json:"excerpt"`
	ImageURL  string `json:"image_url"`
	Published bool   `json:"published"`
}

func CreatePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var title, content, excerpt, imageURL string
		var published bool

		ct := r.Header.Get("Content-Type")

		// Handle multipart form data (for image uploads)
		if strings.HasPrefix(ct, "multipart/form-data") {
			if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
				http.Error(w, "invalid form", http.StatusBadRequest)
				return
			}

			title = r.FormValue("title")
			content = r.FormValue("content")
			excerpt = r.FormValue("excerpt")
			published = r.FormValue("published") == "on" || r.FormValue("published") == "true"

			// Handle image upload
			file, header, err := r.FormFile("image")
			if err == nil {
				defer file.Close()
				imageURL, err = storage.SaveImage(file, header)
				if err != nil {
					http.Error(w, "failed to save image: "+err.Error(), http.StatusBadRequest)
					return
				}
			}
		} else if strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "invalid form", http.StatusBadRequest)
				return
			}

			title = r.FormValue("title")
			content = r.FormValue("content")
			excerpt = r.FormValue("excerpt")
			published = r.FormValue("published") == "on"
		} else {
			var req struct {
				Title     string `json:"title"`
				Content   string `json:"content"`
				Excerpt   string `json:"excerpt"`
				ImageURL  string `json:"image_url"`
				Published bool   `json:"published"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			title = req.Title
			content = req.Content
			excerpt = req.Excerpt
			imageURL = req.ImageURL
			published = req.Published
		}

		if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
			http.Error(w, "title and content required", http.StatusBadRequest)
			return
		}

		slug := strings.ToLower(strings.TrimSpace(title))
		slug = strings.ReplaceAll(slug, " ", "-")

		result, err := db.Exec(`
			INSERT INTO posts (title, slug, content, excerpt, image_url, published, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`,
			title,
			slug,
			content,
			nullString(excerpt),
			nullString(imageURL),
			boolToInt(published),
			time.Now(),
		)
		if err != nil {
			http.Error(w, "failed to insert post", http.StatusInternalServerError)
			return
		}

		if strings.HasPrefix(ct, "application/x-www-form-urlencoded") || strings.HasPrefix(ct, "multipart/form-data") {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}

		postID, _ := result.LastInsertId()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int64{"id": postID})
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

		// Get existing post to check for old image
		var oldImageURL sql.NullString
		err = db.QueryRow("SELECT image_url FROM posts WHERE id = ?", id).Scan(&oldImageURL)
		if err != nil {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		var title, content, excerpt, imageURL string
		var published bool
		ct := r.Header.Get("Content-Type")

		// Handle multipart form data (for image uploads)
		if strings.HasPrefix(ct, "multipart/form-data") {
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				http.Error(w, "invalid form", http.StatusBadRequest)
				return
			}

			title = r.FormValue("title")
			content = r.FormValue("content")
			excerpt = r.FormValue("excerpt")
			published = r.FormValue("published") == "on" || r.FormValue("published") == "true"

			// Handle image upload
			file, header, err := r.FormFile("image")
			if err == nil {
				defer file.Close()
				imageURL, err = storage.SaveImage(file, header)
				if err != nil {
					http.Error(w, "failed to save image: "+err.Error(), http.StatusBadRequest)
					return
				}
				// Delete old image if new one uploaded
				if oldImageURL.Valid && oldImageURL.String != "" {
					storage.DeleteImage(oldImageURL.String)
				}
			} else {
				// Keep existing image if no new upload
				imageURL = oldImageURL.String
			}
		} else if strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "invalid form", http.StatusBadRequest)
				return
			}
			title = r.FormValue("title")
			content = r.FormValue("content")
			excerpt = r.FormValue("excerpt")
			published = r.FormValue("published") == "on"
			imageURL = oldImageURL.String
		} else {
			var req UpdatePostRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			title = req.Title
			content = req.Content
			excerpt = req.Excerpt
			imageURL = req.ImageURL
			published = req.Published
		}

		slug := strings.ToLower(strings.TrimSpace(title))
		slug = strings.ReplaceAll(slug, " ", "-")

		_, err = db.Exec(`
			UPDATE posts
			SET title = ?, slug = ?, content = ?, excerpt = ?, image_url = ?, published = ?
			WHERE id = ?
		`, title, slug, content, nullString(excerpt), nullString(imageURL), boolToInt(published), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if strings.HasPrefix(ct, "application/x-www-form-urlencoded") || strings.HasPrefix(ct, "multipart/form-data") {
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
			SELECT id, title, slug, content, excerpt, image_url, created_at, published
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
			var excerpt, imageURL sql.NullString
			if err := rows.Scan(
				&p.ID,
				&p.Title,
				&p.Slug,
				&p.Content,
				&excerpt,
				&imageURL,
				&p.CreatedAt,
				&p.Published,
			); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			p.Excerpt = excerpt.String
			p.ImageURL = imageURL.String
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
		SELECT id, title, slug, content, excerpt, image_url, created_at, published
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
		var excerpt, imageURL sql.NullString
		if err := rows.Scan(
			&p.ID, &p.Title, &p.Slug,
			&p.Content, &excerpt, &imageURL, &p.CreatedAt, &p.Published,
		); err != nil {
			return nil, err
		}
		p.Excerpt = excerpt.String
		p.ImageURL = imageURL.String
		posts = append(posts, p)
	}
	return posts, nil
}

func GetPostByID(db *sql.DB, id int64) (Post, error) {
	var p Post
	var excerpt, imageURL sql.NullString
	err := db.QueryRow(`
		SELECT id, title, slug, content, excerpt, image_url, created_at, published
		FROM posts
		WHERE id = ?
	`, id).Scan(
		&p.ID, &p.Title, &p.Slug,
		&p.Content, &excerpt, &imageURL, &p.CreatedAt, &p.Published,
	)
	p.Excerpt = excerpt.String
	p.ImageURL = imageURL.String
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

		// Get image URL before deleting
		var imageURL sql.NullString
		db.QueryRow("SELECT image_url FROM posts WHERE id = ?", id).Scan(&imageURL)

		if _, err := db.Exec(`DELETE FROM posts WHERE id = ?`, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Delete associated image
		if imageURL.Valid && imageURL.String != "" {
			storage.DeleteImage(imageURL.String)
		}

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
