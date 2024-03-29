package parser

import (
	"html/template"
	"time"
)

type FeedEnclosure struct {
	Type   string
	Href   string
	Length int
}

type FeedItem struct {
	ID          int
	IsRead      bool
	EntryID     string
	Title       string
	Summary     string
	Link        string
	Content     template.HTML
	Date        time.Time
	IsDateValid bool
	Enclosures  []FeedEnclosure
}

type Feed struct {
	ID           int
	Title        string
	Description  string
	Link         string
	FeedLink     string
	Refresh      time.Time
	Items        []*FeedItem
	ItemCheckMap map[string]struct{}
}

func Parse(data []byte) (*Feed, error) {
	return parseAtom(data)
}
