package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/michaelcosj/pluto-reader/db/repository"
	"github.com/michaelcosj/pluto-reader/model"
)

type UserService struct {
	queries *repository.Queries
}

func User(queries *repository.Queries) *UserService {
	return &UserService{queries}
}

type AddFeedParams struct {
	UserID   int32
	FeedID   int32
	FeedName string
}

func (s *UserService) CreateUser(ctx context.Context, userData *model.UserDTO) (int32, error) {
	var userID int32 = 0
	user, err := s.queries.UserGetByOauthSub(ctx, userData.OauthSub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("error checking if user exists %w", err)
		}

		createUserParams := repository.UserCreateParams{
			Email:    userData.Email,
			Name:     userData.Name,
			OauthSub: userData.OauthSub,
		}

		if userID, err = s.queries.UserCreate(ctx, createUserParams); err != nil {
			return 0, fmt.Errorf("error creating user %w", err)
		}
	} else {
		userID = user.ID
	}

	return userID, nil
}

func (s *UserService) AddFeedToUser(ctx context.Context, userID, feedID int32, feedName string) error {
	err := s.queries.UserAddFeed(ctx, repository.UserAddFeedParams{
		UserID:         userID,
		FeedID:         feedID,
		FeedName:       feedName,
		UpdateInterval: pgtype.Interval{Microseconds: time.Hour.Microseconds()},
	})

	if err != nil {
		return fmt.Errorf("error adding feed to user feeds: %w", err)
	}

	feedItems, err := s.queries.FeedGetItems(ctx, feedID)
	if err != nil {
		return fmt.Errorf("error getting feed items: %w", err)
	}

	var addFeedItemsParams []repository.UserAddFeedItemsParams
	for _, item := range feedItems {
		addFeedItemsParams = append(addFeedItemsParams, repository.UserAddFeedItemsParams{
			UserID: userID,
			ItemID: item.ID,
		})
	}

	_, err = s.queries.UserAddFeedItems(ctx, addFeedItemsParams)
	if err != nil {
		return fmt.Errorf("error adding feed items to user feed items: %w", err)
	}

	return nil
}
