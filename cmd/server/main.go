package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/henryarin/portfolio-backend-go/internal/api"
	admin "github.com/henryarin/portfolio-backend-go/internal/api/admin"
	"github.com/henryarin/portfolio-backend-go/internal/auth"
	"github.com/henryarin/portfolio-backend-go/internal/config"
	"github.com/henryarin/portfolio-backend-go/internal/db"
	"github.com/henryarin/portfolio-backend-go/internal/handlers"
	"github.com/henryarin/portfolio-backend-go/internal/middleware"
	"github.com/henryarin/portfolio-backend-go/internal/web"
)

func main() {
	cfg := config.Load()

	database := db.Open(cfg.DBPath)
	if err := db.Init(database); err != nil {
		log.Fatal(err)
	}

	api.SetDB(database)

	isProd := os.Getenv("ENV") == "prod"

	sessionManager := auth.NewSessionManager(database, isProd)

	adminAuth := &handlers.AdminAuthHandler{
		Sessions: sessionManager,
	}

	requireAdminSession := middleware.RequireAdminSession(sessionManager)

	mux := http.NewServeMux()

	tmpl := template.Must(template.ParseFiles(
		filepath.Join("internal", "web", "templates", "layout.html"),
		filepath.Join("internal", "web", "templates", "blog_index.html"),
	))

	mux.HandleFunc("/blog", web.BlogIndex(tmpl))

	adminLimiter := middleware.NewRateLimiter(5, time.Minute)

	mux.HandleFunc("/api/admin/login", adminAuth.Login)
	mux.HandleFunc("/api/admin/logout", adminAuth.Logout)

	mux.Handle(
		"/api/admin/posts",
		adminLimiter.Middleware(
			requireAdminSession(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch r.Method {
					case http.MethodGet:
						admin.ListPosts(database)(w, r)
					case http.MethodPost:
						admin.CreatePost(database)(w, r)
					default:
						http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
					}
				}),
			),
		),
	)

	mux.Handle(
		"/api/admin/posts/",
		adminLimiter.Middleware(
			requireAdminSession(
				admin.UpdatePost(database),
			),
		),
	)

	mux.HandleFunc("/api/posts", api.ListPosts)
	mux.HandleFunc("/api/posts/", api.GetPostBySlug)
	mux.HandleFunc("/api/posts/previews", api.ListPostPreviews)
	mux.HandleFunc("/blog/", web.BlogShow(tmpl))

	handler := sessionManager.LoadAndSave(
		middleware.CORS(cfg.AllowedOrigin, mux),
	)

	log.Println("listening on :" + cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
}
