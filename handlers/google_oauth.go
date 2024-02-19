package handlers

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/michaelcosj/pluto-reader/services"
	"github.com/michaelcosj/pluto-reader/views/pages"
)

type GoogleOauthHandler struct {
	service        *services.GoogleOauthService
	sessionManager *scs.SessionManager
}

func GoogleOauth(service *services.GoogleOauthService, sessionManager *scs.SessionManager) *GoogleOauthHandler {
	return &GoogleOauthHandler{service, sessionManager}
}

func (h *GoogleOauthHandler) ShowSignInPage(w http.ResponseWriter, r *http.Request) {
	userId := h.sessionManager.GetInt32(r.Context(), "userID")
	if userId > 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	authPage := pages.Auth()
	authPage.Render(r.Context(), w)
}

func (h *GoogleOauthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	url := h.service.GetAuthUrl()
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *GoogleOauthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != "qwerty" {
		log.Fatalf("invalid state %s", state)
	}

	code := r.URL.Query().Get("code")
	if err := h.service.Authenticate(r.Context(), h.sessionManager, code); err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *GoogleOauthHandler) Logout(w http.ResponseWriter, r *http.Request) {
}
