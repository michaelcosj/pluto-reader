package parser

import "github.com/michaelcosj/pluto-reader/models"

func Parse(data []byte) (*models.Feed, error) {
	return parseAtom(data)
}
