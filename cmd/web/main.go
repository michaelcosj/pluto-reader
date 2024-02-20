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
	"github.com/michaelcosj/pluto-reader/handlers"
	"github.com/michaelcosj/pluto-reader/services"
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
	usersService := services.Users(queries)

	indexHandler := handlers.Index()
	r.Group(func(r chi.Router) {
		r.Get("/", indexHandler.ShowIndexPage)
		r.Post("/getfeed", indexHandler.GetFeed)
	})

	googleOauthHandler := handlers.GoogleOauth(usersService, sessionManager)
	r.Group(func(r chi.Router) {
		r.Get("/auth", googleOauthHandler.ShowSignInPage)
		r.Get("/auth/signin", googleOauthHandler.SignIn)
		r.Get("/auth/callback", googleOauthHandler.Callback)
	})

	r.Handle("/dist/*", assets.Mount())

	http.ListenAndServe(":3000", sessionManager.LoadAndSave(r))
}
