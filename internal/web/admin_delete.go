package web

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	adminapi "github.com/henryarin/portfolio-backend-go/internal/api/admin"
)

type AdminDeletePostData struct {
	Post adminapi.Post
}

func AdminDeletePost(t *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 4 {
			http.NotFound(w, r)
			return
		}

		id, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		post, err := adminapi.GetPostByID(db, id)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := t.ExecuteTemplate(w, "layout", AdminDeletePostData{Post: post}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
