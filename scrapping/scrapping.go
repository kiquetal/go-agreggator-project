package scrapping

import (
	"context"
	"database/sql"
	"encoding/xml"
	"github.com/kiquetal/go-agreggator-project/internal/database"
	"log"
	"net/http"
	"sync"
	"time"
)

func FetchFeeds(db *database.Queries) {
	// Fetch from the database
	feedToFetch, err := db.GetNexFeedsToFetch(context.Background())
	if err != nil {
		log.Fatalf("error fetching feeds to fetch: %v", err)
	}
	// iterate over the feeds to fetch
	var allFetchedFeeds []Item = make([]Item, 0)
	var wg sync.WaitGroup
	wg.Add(len(feedToFetch))
	for _, feed := range feedToFetch {
		go func(feed database.Feed, db *database.Queries) {
			defer wg.Done()
			items := downloadFeed(feed)
			allFetchedFeeds = append(allFetchedFeeds, items...)
			feed, err := db.MarkedFetched(context.Background(), database.MarkedFetchedParams{
				ID:        feed.ID,
				UpdatedAt: time.Now(),
				LastFetchedAt: sql.NullTime{
					Time: time.Now(),
				},
			})
			if err != nil {
				log.Printf("error marking feed as fetched: %v", err)

			}
		}(feed, db)

	}
	wg.Wait()
	for _, item := range allFetchedFeeds {
		log.Printf("item: %v\n", item.Title)

	}
}

func downloadFeed(feed database.Feed) (allFetchedFeeds []Item) {
	response, err := http.Get(feed.Url)
	if err != nil {
		log.Printf("error downloading feed %s: %v", feed.Url, err)
		return
	}
	defer response.Body.Close()
	var RssFeed Rss
	dataFromServer := xml.NewDecoder(response.Body)
	err = dataFromServer.Decode(&RssFeed)
	if err != nil {
		log.Printf("error parsing feed %s: %v", feed.Url, err)
		return

	}
	return RssFeed.Channel.Items
}

type Rss struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Guid        string `xml:"guid""`
	Description string `xml:"description""`
}
