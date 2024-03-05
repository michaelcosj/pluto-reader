package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/michaelcosj/pluto-reader/parser"
	"github.com/michaelcosj/pluto-reader/util"
	"github.com/michaelcosj/pluto-reader/view/component"
	"github.com/michaelcosj/pluto-reader/view/page"
)

type IndexHandler struct{
    sessionManager *scs.SessionManager
}

func Index(sessionManager *scs.SessionManager) *IndexHandler {
	return &IndexHandler{sessionManager}
}

func (h *IndexHandler) ShowIndexPage(w http.ResponseWriter, r *http.Request) {
	userID := h.sessionManager.GetInt32(r.Context(), "userID")
	if userID == 0 {
		log.Printf("user not authenticated\n")
		http.Redirect(w, r, "/auth/", http.StatusSeeOther)
	}

	home := page.Index()
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

	_, err = parser.Parse(data)
	if err != nil {
		log.Fatalf("error parsing data, %v\n", err)
	}

	cardList := component.CardList([]string{"hello", "world"})
	cardList.Render(r.Context(), w)
}
