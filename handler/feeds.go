package handlers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/michaelcosj/pluto-reader/parser"
	"github.com/michaelcosj/pluto-reader/utils"
	"github.com/michaelcosj/pluto-reader/views/components"
)

type FeedsHandler struct {
}

func Feeds() *FeedsHandler {
	return &FeedsHandler{}
}

func (h *FeedsHandler) AddFeed(w http.ResponseWriter, r *http.Request) {
	link := r.FormValue("link")
	_, err := url.ParseRequestURI(link)
	if err != nil {
		log.Fatalf("error parsing url: %v\n", err)
	}

	feedName := r.FormValue("name")
	if len(feedName) < 3 {
		log.Fatalf("invalid feed name: %s\n", feedName)
	}

	data, err := utils.Fetch(link)
	if err != nil {
		log.Fatalf("error fetching url %s: %v\n", link, err)
	}

	feed, err := parser.Parse(data)
	if err != nil {
		log.Fatalf("error parsing feed data: %v\n", err)
	}

	cardComponent := components.Card(feed.Title, feed.Description, feed.FeedLink)
	cardComponent.Render(r.Context(), w)

    // TODO: save feed in the database
}
