package handler

import (
	"log"
	"net/http"
	"net/url"

	"github.com/alexedwards/scs/v2"
	"github.com/michaelcosj/pluto-reader/service"
	pages "github.com/michaelcosj/pluto-reader/view/page"
)

type FeedsHandler struct {
	feedService    *service.FeedService
	userService    *service.UserService
	sessionManager *scs.SessionManager
}

func Feeds(feedService *service.FeedService, userService *service.UserService, sessionManager *scs.SessionManager) *FeedsHandler {
	return &FeedsHandler{feedService, userService, sessionManager}
}

func (h *FeedsHandler) AddFeed(w http.ResponseWriter, r *http.Request) {
	userID := h.sessionManager.GetInt32(r.Context(), "userID")
	if userID == 0 {
		log.Fatalf("user not authenticated\n")
	}

	feedUrl := r.FormValue("url")
	_, err := url.ParseRequestURI(feedUrl)
	if err != nil {
		log.Fatalf("error parsing url: %v\n", err)
	}

	feedName := r.FormValue("name")
	if len(feedName) < 3 {
		log.Fatalf("invalid feed name: %s\n", feedName)
	}

	feedID, err := h.feedService.ParseAndCreateFeed(r.Context(), feedUrl)
	if err != nil {
		log.Fatalf("error creating feed: %v\n", err)
	}

	err = h.userService.AddFeedToUser(r.Context(), userID, feedID, feedName)
	if err != nil {
		log.Fatalf("error adding feed to user: %v\n", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *FeedsHandler) GetUserFeedItems(w http.ResponseWriter, r *http.Request)  {
	userID := h.sessionManager.GetInt32(r.Context(), "userID")
	if userID == 0 {
		log.Fatalf("user not authenticated\n")
	}
    
    feedItems, err := h.userService.GetUserFeedItems(r.Context(), userID)
    if err != nil {
        log.Fatalf("error getting user feed items: %v\n", err)
    }

    feedItemList := pages.FeedItemList(feedItems)
    feedItemList.Render(r.Context(), w)
}
