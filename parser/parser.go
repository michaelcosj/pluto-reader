package parser

import "github.com/michaelcosj/pluto-reader/model"

func Parse(data []byte) (*models.FeedDTO, error) {
	return parseAtom(data)
}
