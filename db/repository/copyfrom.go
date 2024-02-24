// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: copyfrom.go

package repository

import (
	"context"
)

// iteratorForFeedAddItems implements pgx.CopyFromSource.
type iteratorForFeedAddItems struct {
	rows                 []FeedAddItemsParams
	skippedFirstNextCall bool
}

func (r *iteratorForFeedAddItems) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForFeedAddItems) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].EntryID,
		r.rows[0].Title,
		r.rows[0].Summary,
		r.rows[0].Link,
		r.rows[0].Content,
		r.rows[0].ItemUpdated,
		r.rows[0].ItemPublished,
		r.rows[0].FeedID,
	}, nil
}

func (r iteratorForFeedAddItems) Err() error {
	return nil
}

func (q *Queries) FeedAddItems(ctx context.Context, arg []FeedAddItemsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"feed_items"}, []string{"entry_id", "title", "summary", "link", "content", "item_updated", "item_published", "feed_id"}, &iteratorForFeedAddItems{rows: arg})
}

// iteratorForUserAddFeedItems implements pgx.CopyFromSource.
type iteratorForUserAddFeedItems struct {
	rows                 []UserAddFeedItemsParams
	skippedFirstNextCall bool
}

func (r *iteratorForUserAddFeedItems) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForUserAddFeedItems) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].UserID,
		r.rows[0].ItemID,
	}, nil
}

func (r iteratorForUserAddFeedItems) Err() error {
	return nil
}

func (q *Queries) UserAddFeedItems(ctx context.Context, arg []UserAddFeedItemsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"user_feed_items"}, []string{"user_id", "item_id"}, &iteratorForUserAddFeedItems{rows: arg})
}
