package parser

import "github.com/michaelcosj/pluto-reader/model"

func Parse(data []byte) (*model.FeedDTO, error) {
	return parseAtom(data)
}
