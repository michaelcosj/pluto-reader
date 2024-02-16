package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/michaelcosj/pluto-reader/views"
)

/*

Features
- user accounts
- use google oauth for authentication
- data pesistence in sqlite
- add rss feed from site url
- open feed url in a seperate tab
- create collections to group feeds (an all collection is created by default and stores all feeds)
- bookmark a feed update
- favourite a feed (stored as a favourites collection which is created by default)
- setting to refresh feed on a certain interval
- service workers and cache for offline functionality

timeline 3 weeks

*/

func main() {
	// Create a simple form that accepts a feed url and returns all its articles
	fmt.Println("Pluto!")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		component := views.Hello("Michael")
		component.Render(context.Background(), w)
	})

	http.ListenAndServe(":3000", r)
}
