// internal/web/blog.go
package web

import (
	"html/template"
	"net/http"

	"github.com/henryarin/portfolio-backend-go/internal/api"
)

type BlogIndexData struct {
	Posts []api.Post
}

func BlogIndex(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := api.GetPublishedPosts()
		if err != nil {
			http.Error(w, "failed to load posts", http.StatusInternalServerError)
			return
		}

		data := BlogIndexData{
			Posts: posts,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := t.ExecuteTemplate(w, "layout", data); err != nil {
			http.Error(w, "template render failed", http.StatusInternalServerError)
			return
		}
	}
}

type BlogPostData struct {
	Post api.Post
}

func BlogShow(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.URL.Path[len("/blog/"):]
		if slug == "" {
			http.NotFound(w, r)
			return
		}

		post, err := api.GetPublishedPostBySlug(slug)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		data := BlogPostData{
			Post: post,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := t.ExecuteTemplate(w, "layout", data); err != nil {
			http.Error(w, "template render failed", http.StatusInternalServerError)
			return
		}
	}
}
