package parser

import "github.com/michaelcosj/pluto-reader/models"

func Parse(data []byte) (*models.FeedDTO, error) {
	return parseAtom(data)
}
