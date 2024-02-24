package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/michaelcosj/pluto-reader/assets"
	"github.com/michaelcosj/pluto-reader/db/repository"
	"github.com/michaelcosj/pluto-reader/handler"
	"github.com/michaelcosj/pluto-reader/service"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)

	sessionManager := scs.New()
	sessionManager.Store = memstore.New()

	ctx := context.Background()
	dbConn, err := pgx.Connect(ctx, os.Getenv("PG_DSN"))
	if err != nil {
		log.Fatalf("error connecting to database: %w", err)
	}
	defer dbConn.Close(ctx)

	queries := repository.New(dbConn)

	usersService := service.User(queries)
	feedService := service.Feed(queries)

	indexHandler := handler.Index()
	r.Get("/", indexHandler.ShowIndexPage)
	r.Post("/getfeed", indexHandler.GetFeed)

	googleOauthHandler := handler.GoogleOauth(usersService, sessionManager)
	r.Route("/auth", func(r chi.Router) {
		r.Get("/", googleOauthHandler.Index)
		r.Get("/signin", googleOauthHandler.SignIn)
		r.Get("/callback", googleOauthHandler.Callback)
	})

	feedsHandler := handler.Feeds(feedService, usersService, sessionManager)
	r.Route("/feed", func(r chi.Router) {
		r.Post("/new", feedsHandler.AddFeed)
	})

	r.Handle("/dist/*", assets.Mount())
	http.ListenAndServe(":3000", sessionManager.LoadAndSave(r))
}
