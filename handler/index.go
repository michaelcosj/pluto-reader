package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/michaelcosj/pluto-reader/parser"
	"github.com/michaelcosj/pluto-reader/util"
	"github.com/michaelcosj/pluto-reader/view/component"
	"github.com/michaelcosj/pluto-reader/view/page"
)

type IndexHandler struct { }

func Index() *IndexHandler {
	return &IndexHandler{}
}

func (h *IndexHandler) ShowIndexPage(w http.ResponseWriter, r *http.Request) {
	home := pages.Index()
	home.Render(context.Background(), w)
}

func (h *IndexHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	link := r.FormValue("link")
	if link == "" {
		log.Fatalf("invalid url")
		return
	}

	data, err := util.Fetch(link)
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
