package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/michaelcosj/pluto-reader/parser"
	"github.com/michaelcosj/pluto-reader/utils"
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
		component := views.Home("Michael")
		component.Render(context.Background(), w)
	})

	r.Post("/getfeed", func(w http.ResponseWriter, r *http.Request) {
		link := r.FormValue("link")
		if link == "" {
			log.Fatal("invalid url")
		}

		data, err := utils.Fetch(link)
		if err != nil {
			log.Fatalf("error fetching url, %v", err)
		}

		feed, err := parser.Parse(data)
		if err != nil {
			log.Fatalf("error parsing data, %v", err)
		}

		fmt.Println(*feed)

		cardList := views.CardList(feed)
		cardList.Render(context.Background(), w)

	})

	fs := http.FileServer(http.Dir("./public/"))
	r.Handle("/*", fs)
	http.ListenAndServe(":3000", r)
}
