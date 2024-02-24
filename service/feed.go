package service

import (
	"context"
	"fmt"

	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/michaelcosj/pluto-reader/db/repository"
	"github.com/michaelcosj/pluto-reader/parser"
)

type FeedService struct {
	queries *repository.Queries
}

func Feed(queries *repository.Queries) *FeedService {
	return &FeedService{queries}
}

func (s *FeedService) ParseAndCreateFeed(ctx context.Context, feedUrl string) (int32, error) {
	data, err := fetchData(feedUrl)
	if err != nil {
		return 0, fmt.Errorf("error fetching url %s: %w\n", feedUrl, err)
	}

	feedModel, err := parser.Parse(data)
	if err != nil {
		return 0, fmt.Errorf("error parsing feed data: %v\n", err)
	}

	feedID, err := s.queries.FeedCreate(ctx, repository.FeedCreateParams{
		Title:       pgtype.Text{String: feedModel.Title, Valid: true},
		Description: pgtype.Text{String: feedModel.Description, Valid: true},
		SiteLink:    pgtype.Text{String: feedModel.Link, Valid: true},
		FeedLink:    feedModel.FeedLink,
	})

	if err != nil {
		return 0, fmt.Errorf("error inserting feed in database: %w", err)
	}

	var itemsParams []repository.FeedAddItemsParams
	for _, item := range feedModel.Items {
		if item != nil {
			itemsParams = append(itemsParams, repository.FeedAddItemsParams{
				Link:        item.Link,
				FeedID:      feedID,
				EntryID:     pgtype.Text{String: item.EntryID, Valid: true},
				Title:       pgtype.Text{String: item.Title, Valid: true},
				Summary:     pgtype.Text{String: item.Summary, Valid: true},
				Content:     pgtype.Text{String: item.Content, Valid: true},
				ItemUpdated: pgtype.Timestamptz{Time: item.Date, Valid: item.IsDateValid},
			})
		}
	}

	_, err = s.queries.FeedAddItems(ctx, itemsParams)
	if err != nil {
		return 0, fmt.Errorf("error inserting feed items in database: %w", err)
	}

	return feedID, nil
}

func fetchData(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf(res.Status)
	}

	return io.ReadAll(res.Body)
}
