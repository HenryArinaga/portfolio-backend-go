package web

import (
	"html/template"
	"net/http"

	"database/sql"

	adminapi "github.com/henryarin/portfolio-backend-go/internal/api/admin"
)

type AdminDashboardData struct {
	Posts []adminapi.Post
}

func AdminDashboard(t *template.Template, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := adminapi.ListPostsData(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := t.ExecuteTemplate(w, "layout", AdminDashboardData{Posts: posts}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
