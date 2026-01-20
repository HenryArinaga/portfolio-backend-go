package web

import (
	"html/template"
	"net/http"
)

type BlogIndexData struct {
	Posts []string
}

func BlogIndex(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := BlogIndexData{
			Posts: []string{"First SSR post", "Second SSR post"},
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := t.ExecuteTemplate(w, "layout", data); err != nil {
			http.Error(w, "template render failed", http.StatusInternalServerError)
			return
		}
	}
}
