package parser

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type atomFeedData struct {
	title       string
	description string
	link        string
	feedLink    string
}

type atomItemData struct {
	title   string
	summary string
	link    string
	id      string
	content string
}

func TestParseAtomFeed(t *testing.T) {
	tests := map[string]atomFeedData{
		"atom_test.xml": {
			title:       "Example Feed",
			description: "A subtitle.",
			link:        "http://example.org/",
			feedLink:    "http://example.org/feed/",
		},
	}

	for test, want := range tests {
		fileName := filepath.Join("../testdata", test)
		data, err := os.ReadFile(fileName)
		if err != nil {
			t.Fatalf("reading file %s :%v", fileName, err)
		}

		feed, err := parseAtom(data)
		if err != nil {
			t.Fatalf("parsing data from %s :%v", fileName, err)
		}

		if feed.Title != want.title {
			t.Fatalf("[title] got %s expected %s in %s", feed.Title, want.title, fileName)
		}

		if feed.Description != want.description {
			t.Fatalf("[description] got %s expected %s in %s", feed.Description, want.description, fileName)
		}

		if feed.Link != want.link {
			t.Fatalf("[link] got %s expected %s in %s", feed.Link, want.link, fileName)
		}

		if feed.FeedLink != want.feedLink {
			t.Fatalf("[feedlink] got %s expected %s in %s", feed.FeedLink, want.feedLink, fileName)
		}
	}
}

func TestParseAtomItems(t *testing.T) {
	tests := map[string][]atomItemData{
		"atom_test.xml": {
			{
				title:   "Atom-Powered Robots Run Amok",
				summary: "Some text.",
				link:    "http://example.org/2003/12/13/atom03",
				id:      "urn:uuid:1225c695-cfb8-4ebb-aaaa-80da344efa6a",
				content: "<p>Hello World</p>",
			},
		},
	}

	for test, want := range tests {
		fileName := filepath.Join("../testdata", test)
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			t.Fatalf("reading file %s :%v", fileName, err)
		}

		feed, err := parseAtom(data)
		if err != nil {
			t.Fatalf("parsing data from %s :%v", fileName, err)
		}

		if len(feed.Items) != len(want) {
			t.Fatalf("got %d number of items, expected %d", len(feed.Items), len(want))
		}

		for k, item := range feed.Items {
			if item.Title != want[k].title {
				t.Fatalf("[title] got %s expected %s in %s", item.Title, want[k].title, fileName)
			}

			if item.Summary != want[k].summary {
				t.Fatalf("[summary] got %s expected %s in %s", item.Summary, want[k].summary, fileName)
			}

			if item.Link != want[k].link {
				t.Fatalf("[link] got %s expected %s in %s", item.Link, want[k].link, fileName)
			}

			if item.ItemID != want[k].id {
				t.Fatalf("[id] got %s expected %s in %s", item.ItemID, want[k].id, fileName)
			}

			if strings.TrimSpace(item.Content) != want[k].content {
				t.Fatalf("[content] got %s expected %s in %s", item.Content, want[k].content, fileName)
			}

			if !item.IsDateValid {
				t.Fatalf("[data] item date %q is not valid in %s", item.Date, fileName)
			}
		}

	}
}
