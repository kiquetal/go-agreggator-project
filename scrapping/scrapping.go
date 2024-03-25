package scrapping

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
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
					Time:  time.Now(),
					Valid: true,
				},
			})
			if err != nil {
				log.Printf("error marking feed as fetched: %v", err)

			}
		}(feed, db)

	}
	wg.Wait()
	for _, item := range allFetchedFeeds {
		// insert into post

		timeToDB, errDate := convertData(item.PubDate)
		isValid := true
		if errDate != nil {
			log.Printf("error parsing date: %v", errDate)
			isValid = false
		}
		_, err := db.InsertPost(context.Background(), database.InsertPostParams{
			ID: uuid.New(),
			Title: sql.NullString{
				String: item.Title,
				Valid:  true,
			},
			Url: sql.NullString{
				String: item.Link,
				Valid:  true,
			},
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: sql.NullTime{
				Time:  timeToDB,
				Valid: isValid,
			},
			FeedID: uuid.NullUUID{
				UUID:  uuid.MustParse(item.FeedId),
				Valid: true,
			},
		})
		if err != nil {
			log.Printf("error inserting post: %v", err)
		}
		log.Println("inserted post: ", item.Title)

	}
}

func convertData(date string) (time.Time, error) {
	var t time.Time
	var err error
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
	}

	for _, format := range formats {
		t, err = time.Parse(format, date)
		if err == nil {
			return t, nil
		}
	}

	return t, fmt.Errorf("unable to parse date: %s", date)

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
	var items []Item
	for _, item := range RssFeed.Channel.Items {
		item.FeedId = feed.ID.String()
		items = append(items, item)
	}
	return items
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
	PubDate     string `xml:"pubDate""`
	FeedId      string
}
