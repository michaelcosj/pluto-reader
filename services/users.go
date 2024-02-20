package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/michaelcosj/pluto-reader/db/repository"
)

type UsersService struct {
	queries *repository.Queries
}

func Users(queries *repository.Queries) *UsersService {
	return &UsersService{
		queries: queries,
	}
}

func (srv *UsersService) CreateUser(context context.Context, email, name, oauthSub string) (*repository.User, error) {
	user, err := srv.queries.GetUserByOauthSub(context, oauthSub)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("error checking if user exists %w", err)
		}

		createUserParams := repository.CreateUserParams{
			Email:    email,
			Name:     name,
			OauthSub: oauthSub,
		}

		user, err = srv.queries.CreateUser(context, createUserParams)
		if err != nil {
			return nil, fmt.Errorf("error creating user %w", err)
		}
	}

	return &user, nil
}
