package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/phmshk/bootdev-blog-aggr/internal/api"
	"github.com/phmshk/bootdev-blog-aggr/internal/database"
)

func scrapeFeeds(s *state) error {
	nextFeedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error scraping feeds: %v", err)
	}
	_, err = s.db.MarkFeedFetched(context.Background(), nextFeedToFetch.ID)
	if err != nil {
		return fmt.Errorf("error marking feed: %v", err)
	}

	feed, err := api.FetchFeed(context.Background(), nextFeedToFetch.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	var countPosts int
	for _, item := range feed.Channel.Item {
		pubDate := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			pubDate = sql.NullTime{Time: t, Valid: true}
		} else if t, err := time.Parse(time.RFC1123, item.PubDate); err == nil {
			pubDate = sql.NullTime{Time: t, Valid: true}
		}

		description := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: pubDate,
			FeedID:      nextFeedToFetch.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
		countPosts++
	}

	fmt.Println(countPosts, " posts added")
	return nil
}
