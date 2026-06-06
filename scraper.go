package main

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"rssag/internal/database"
	"sync"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func startScraping(db *database.Queries, concurrency int, timeout time.Duration) {
	log.Printf("Starting scraping on %v goroutines every %s duration", concurrency, timeout)
	ticker := time.NewTicker(timeout)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("error fetching feeds: %s", err)
		}
		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(wg, db, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("error scraping feed %v: %v", feed.ID, err)
	}
	rssFeed, err := FetchFeed(feed.Url)
	if err != nil {
		log.Printf("error fetching feed %v: %v", feed.ID, err)
		return
	}
	for _, item := range rssFeed.Channel.Item {
		log.Printf("Feed collected: %s - %s, Totally: %v", item.Title, item.Link, len(rssFeed.Channel.Item))
	}
}

func FetchFeed(feedUrl string) (RSSFeed, error) {
	httpClient := &http.Client{}
	resp, err := httpClient.Get(feedUrl)
	if err != nil {
		return RSSFeed{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, err
	}
	rssFeed := RSSFeed{}
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}
	return rssFeed, nil
}
