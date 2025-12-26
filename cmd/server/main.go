package main

import (
	"log"
	"net/http"

	"github.com/henryarin/portfolio-backend-go/internal/api"
	admin "github.com/henryarin/portfolio-backend-go/internal/api/admin"
	"github.com/henryarin/portfolio-backend-go/internal/config"
	"github.com/henryarin/portfolio-backend-go/internal/db"
	"github.com/henryarin/portfolio-backend-go/internal/middleware"
)

func main() {
	// load config (.env, defaults, etc.)
	cfg := config.Load()

	// open database
	database := db.Open(cfg.DBPath)
	if err := db.Init(database); err != nil {
		log.Fatal(err)
	}

	// set DB for public API
	api.SetDB(database)

	// router
	mux := http.NewServeMux()

	// public blog routes
	mux.HandleFunc("/api/posts", api.ListPosts)
	mux.HandleFunc("/api/posts/", api.GetPostBySlug)

	// admin blog routes
	mux.HandleFunc(
		"/api/admin/posts",
		admin.CreatePost(database, cfg.AdminToken),
	)

	mux.HandleFunc(
		"/api/admin/posts/",
		admin.UpdatePost(database, cfg.AdminToken),
	)

	// middleware
	handler := middleware.CORS(cfg.AllowedOrigin, mux)

	// start server
	log.Println("listening on :" + cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
}
