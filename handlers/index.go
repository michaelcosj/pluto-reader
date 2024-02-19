package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/michaelcosj/pluto-reader/parser"
	"github.com/michaelcosj/pluto-reader/utils"
	"github.com/michaelcosj/pluto-reader/views/components"
	"github.com/michaelcosj/pluto-reader/views/pages"
)

func ShowIndexPage(w http.ResponseWriter, r *http.Request) {
	home := pages.Index()
	home.Render(context.Background(), w)
}

func GetFeed(w http.ResponseWriter, r *http.Request) {
	link := r.FormValue("link")
	if link == "" {
		log.Fatalf("invalid url")
		return
	}

	data, err := utils.Fetch(link)
	if err != nil {
		log.Fatalf("error fetching url, %v\n", err)
		return
	}

	feed, err := parser.Parse(data)
	if err != nil {
		log.Fatalf("error parsing data, %v\n", err)
	}

	cardList := components.CardList(feed)
	cardList.Render(r.Context(), w)
}
