package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/michaelcosj/pluto-reader/models"
	"github.com/michaelcosj/pluto-reader/services"
	"github.com/michaelcosj/pluto-reader/utils"
	"github.com/michaelcosj/pluto-reader/views/pages"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const GOOGLE_USER_PROFILE_API = "https://www.googleapis.com/oauth2/v2/userinfo"

type GoogleOauthHandler struct {
	oauthConfig    *oauth2.Config
	userService    *services.UsersService
	sessionManager *scs.SessionManager
}

func GoogleOauth(service *services.UsersService, sessionManager *scs.SessionManager) *GoogleOauthHandler {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/auth/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleOauthHandler{config, service, sessionManager}
}

func (h *GoogleOauthHandler) Index(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("hello?")
	authPage := pages.Auth()
	authPage.Render(r.Context(), w)
}

func (h *GoogleOauthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	state := utils.GenerateRandomString(5)
	h.sessionManager.Put(r.Context(), "state", state)

	url := h.oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *GoogleOauthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	storedState := h.sessionManager.GetString(r.Context(), "state")
	if state != storedState {
		log.Fatalf("invalid state %s", state)
	}

	code := r.URL.Query().Get("code")
	token, err := h.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Fatalf("error exchanging code %w", err)
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", GOOGLE_USER_PROFILE_API, nil)
	if err != nil {
		log.Fatalf("error creating auth request %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error retrieving user data %w", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading user data %w", err)
	}

	var userData models.UserDTO
	if err := json.Unmarshal(data, &userData); err != nil {
		log.Fatalf("error parsing user data %w", err)
	}

	user, err := h.userService.CreateUser(r.Context(), userData)
	if err := h.sessionManager.RenewToken(r.Context()); err != nil {
		log.Fatalf("error renewing session token %w", err)
	}
	h.sessionManager.Put(r.Context(), "userID", user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *GoogleOauthHandler) Logout(w http.ResponseWriter, r *http.Request) {
}
