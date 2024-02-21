package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/michaelcosj/pluto-reader/db/repository"
	"github.com/michaelcosj/pluto-reader/models"
)

type UsersService struct {
	queries *repository.Queries
}

func Users(queries *repository.Queries) *UsersService {
	return &UsersService{
		queries: queries,
	}
}

func (srv *UsersService) CreateUser(context context.Context, userData models.UserDTO) (*repository.User, error) {
	user, err := srv.queries.GetUserByOauthSub(context, userData.OauthSub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("error checking if user exists %w", err)
		}

		createUserParams := repository.CreateUserParams{
			Email:    userData.Email,
			Name:     userData.Name,
			OauthSub: userData.OauthSub,
		}

		if user, err = srv.queries.CreateUser(context, createUserParams); err != nil {
			return nil, fmt.Errorf("error creating user %w", err)
		}
	}

	return &user, nil
}
