package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/michaelcosj/pluto-reader/internal/router"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
    fmt.Println("pluto reader")
	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}

/*

Features
- user accounts
- use google oauth for authentication
- data pesistence in sqlite
- add rss feed from site url
- open feed url in a seperate tab
- create collections to group feeds (an all collection is created by default and stores all feeds)
- bookmark a feed update
- favourite a feed (stored as a favourites collection which is created by default)
- setting to refresh feed on a certain interval
- service workers and cache for offline functionality

timeline 3 weeks

*/
