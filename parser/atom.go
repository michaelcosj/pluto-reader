package parser

import (
	"encoding/xml"
	"log"

	"github.com/michaelcosj/pluto-reader/util"
)

type atomAuthor struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name"`
	Uri     string   `xml:"uri"`
	Email   string   `xml:"email"`
}

type atomLink struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
	Length  int      `xml:"length,attr"`
}

type atomRawContent struct {
	Raw string `xml:",innerxml"`
}

type atomCategory struct {
	XMLName xml.Name `xml:"category"`
	Label   string   `xml:"label,attr"`
	Term    string   `xml:"term,attr"`
	Scheme  string   `xml:"scheme,attr"`
}

type atomEntry struct {
	XMLName xml.Name       `xml:"entry"`
	ID      string         `xml:"id"`
	Title   string         `xml:"title"`
	Summary string         `xml:"summary"`
	Updated string         `xml:"updated"`
	Links   []atomLink     `xml:"link"`
	Author  atomAuthor     `xml:"author"`
	Content atomRawContent `xml:"content"`
}

type atomFeed struct {
	XMLName  xml.Name    `xml:"feed"`
	Title    string      `xml:"title"`
	Subtitle string      `xml:"subtitle"`
	Updated  string      `xml:"updated"`
	Author   atomAuthor  `xml:"author"`
	Links    []atomLink  `xml:"link"`
	Entries  []atomEntry `xml:"entry"`
}

func parseAtom(data []byte) (*Feed, error) {
	var feed atomFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}

	res := &Feed{
		Title:        feed.Title,
		Description:  feed.Subtitle,
		Items:        make([]*FeedItem, 0),
		ItemCheckMap: make(map[string]struct{}),
	}

	for _, link := range feed.Links {
		if link.Rel == "self" {
			res.FeedLink = link.Href
		} else if link.Rel == "alternate" || link.Rel == "" {
			res.Link = link.Href
		}
	}

	for _, entry := range feed.Entries {
		if _, ok := res.ItemCheckMap[entry.ID]; ok {
			log.Printf("Item %s has duplicate id: %s\n", entry.Title, entry.Updated)
			continue
		}

		item := &FeedItem{
			Title:       entry.Title,
			Summary:     entry.Summary,
			Content:     entry.Content.Raw,
			EntryID:      entry.ID,
			IsRead:      false,
			IsDateValid: true,
		}

		for _, link := range entry.Links {
			if link.Rel == "alternate" || link.Rel == "" {
				item.Link = link.Href
			} else if link.Rel == "enclosure" {
				item.Enclosures = append(item.Enclosures, FeedEnclosure{
					Type:   link.Type,
					Href:   link.Href,
					Length: link.Length,
				})
			}
		}

		if entry.Updated != "" {
			var err error
			item.Date, err = util.ParseTime(entry.Updated)
			if err != nil {
				item.IsDateValid = false
				log.Printf("error parsing date: %s\n", entry.Updated)
			}
		}

		res.Items = append(res.Items, item)
		res.ItemCheckMap[item.EntryID] = struct{}{}
	}

	return res, nil
}
