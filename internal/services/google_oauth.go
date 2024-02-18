package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/michaelcosj/pluto-reader/internal/repository"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const GOOGLE_USER_PROFILE_API = "https://www.googleapis.com/oauth2/v2/userinfo"

type googleUserData struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Sub   string `json:"id"`
}

type GoogleOauthService struct {
	config  *oauth2.Config
	queries *repository.Queries
}

func GoogleOauth(queries *repository.Queries) *GoogleOauthService {
	return &GoogleOauthService{
		config: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/auth/callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		queries: queries,
	}
}

func (srv *GoogleOauthService) GetAuthUrl() string {
	// state := utils.GenerateRandomString(5)
	url := srv.config.AuthCodeURL("qwerty")
	return url
}

func (srv *GoogleOauthService) Authenticate(context context.Context, sessionManager *scs.SessionManager, code string) error {
	token, err := srv.config.Exchange(context, code)
	if err != nil {
		return fmt.Errorf("error exchanging code %w", err)
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", GOOGLE_USER_PROFILE_API, nil)
	if err != nil {
		return fmt.Errorf("error creating auth request %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error retrieving user data %w", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading user data %w", err)
	}

	var userData googleUserData
	if err := json.Unmarshal(data, &userData); err != nil {
		return fmt.Errorf("error parsing user data %w", err)
	}

	user, err := srv.queries.GetUserByOauthSub(context, userData.Sub)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("error checking if user exists %w", err)
		}

		createUserParams := repository.CreateUserParams{
			Email:    userData.Email,
			Name:     userData.Name,
			OauthSub: userData.Sub,
		}

		user, err = srv.queries.CreateUser(context, createUserParams)
		if err != nil {
			return fmt.Errorf("error creating user %w", err)
		}
	}

	if err := sessionManager.RenewToken(context); err != nil {
		return fmt.Errorf("error renewing session token %w", err)
	}
	sessionManager.Put(context, "userID", user.ID)

	return nil
}
