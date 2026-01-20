package main

import (
	"log"
	"net/http"
	"time"

	"github.com/henryarin/portfolio-backend-go/internal/api"
	admin "github.com/henryarin/portfolio-backend-go/internal/api/admin"
	"github.com/henryarin/portfolio-backend-go/internal/config"
	"github.com/henryarin/portfolio-backend-go/internal/db"
	"github.com/henryarin/portfolio-backend-go/internal/middleware"
)

func main() {
	cfg := config.Load()

	database := db.Open(cfg.DBPath)
	if err := db.Init(database); err != nil {
		log.Fatal(err)
	}

	api.SetDB(database)

	mux := http.NewServeMux()

	adminLimiter := middleware.NewRateLimiter(5, time.Minute)

	mux.Handle(
		"/api/admin/posts",
		adminLimiter.Middleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					admin.ListPosts(database, cfg.AdminToken)(w, r)
				case http.MethodPost:
					admin.CreatePost(database, cfg.AdminToken)(w, r)
				default:
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				}
			}),
		),
	)

	mux.Handle(
		"/api/admin/posts/",
		adminLimiter.Middleware(
			admin.UpdatePost(database, cfg.AdminToken),
		),
	)

	mux.HandleFunc("/api/posts", api.ListPosts)
	mux.HandleFunc("/api/posts/", api.GetPostBySlug)

	handler := middleware.CORS(cfg.AllowedOrigin, mux)

	log.Println("listening on :" + cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
}
