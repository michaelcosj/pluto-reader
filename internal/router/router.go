package router

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/michaelcosj/pluto-reader/internal/handlers"
	"github.com/michaelcosj/pluto-reader/internal/repository"
	"github.com/michaelcosj/pluto-reader/internal/services"
)

func Run() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
    r.Use(middleware.CleanPath)

    ctx := context.Background()
    dbConn, err := pgx.Connect(ctx, os.Getenv("PG_DSN"))
    if err != nil {
        return fmt.Errorf("error connecting to database: %w", err)
    }
    defer dbConn.Close(ctx)

    queries := repository.New(dbConn)

    sessionManager := scs.New()
    sessionManager.Store = memstore.New()

	googleAuthService := services.GoogleOauth(queries)
	googleAuthHandler := handlers.GoogleOauth(googleAuthService, sessionManager)

	r.Group(func(r chi.Router) {
        r.Get("/auth", googleAuthHandler.ShowSignInPage)
		r.Get("/auth/signin", googleAuthHandler.SignIn)
		r.Get("/auth/callback", googleAuthHandler.Callback)
	})

	r.Get("/", handlers.ShowIndexPage)
	r.Post("/getfeed", handlers.GetFeed)

	fs := http.FileServer(http.Dir("./assets/dist/"))
	r.Handle("/dist/*", http.StripPrefix("/dist", fs))

	return http.ListenAndServe(":3000", sessionManager.LoadAndSave(r))
}
