package parser

import "github.com/michaelcosj/pluto-reader/internal/models"

func Parse(data []byte) (*models.Feed, error) {
	return parseAtom(data)
}
