package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/michaelcosj/pluto-reader/service"
	"github.com/michaelcosj/pluto-reader/util"
	"github.com/michaelcosj/pluto-reader/view/page"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const OauthProfileAPI = "https://www.googleapis.com/oauth2/v2/userinfo"

type OauthProfileData struct {
	Sub   string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type GoogleOauthHandler struct {
	oauthConfig    *oauth2.Config
	userService    *service.UserService
	sessionManager *scs.SessionManager
}

func GoogleOauth(service *service.UserService, sessionManager *scs.SessionManager) *GoogleOauthHandler {
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
	authPage := page.Auth()
	authPage.Render(r.Context(), w)
}

func (h *GoogleOauthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	state := util.GenerateRandomString(5)
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
		log.Fatalf("error exchanging code %v", err)
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", OauthProfileAPI, nil)
	if err != nil {
		log.Fatalf("error creating auth request %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error retrieving user data %v", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading user data %v", err)
	}

	profile := &OauthProfileData{}
	if err := json.Unmarshal(data, profile); err != nil {
		log.Fatalf("error parsing user data %v", err)
	}

	userID, err := h.userService.CreateUser(r.Context(), profile.Sub, profile.Email, profile.Name)
	if err := h.sessionManager.RenewToken(r.Context()); err != nil {
		log.Fatalf("error renewing session token %v", err)
	}
	h.sessionManager.Put(r.Context(), "userID", userID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *GoogleOauthHandler) Logout(w http.ResponseWriter, r *http.Request) {
}
