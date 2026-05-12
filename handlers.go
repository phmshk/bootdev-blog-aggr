package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/phmshk/bootdev-blog-aggr/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]
	createdUser, err := s.db.CreateUser(
		context.Background(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      name,
		},
	)
	if err != nil {
		return fmt.Errorf("an arror occured while creating a user, %v", err)
	}
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("an error occurred while setting a user, %v", err)
	}
	fmt.Println("User created successfully:")
	printUser(createdUser)
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("no arguments should be provided for this command, %s", cmd.Name)
	}

	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("reset failed: %v", err)
	}

	fmt.Println("DB was successfully reset")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("no arguments should be provided for this command, %s", cmd.Name)
	}

	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("an error occurred fetching all users: %v", err)
	}

	printAllUsers(users, s.cfg.CurrentUserName)
	return nil
}

func printAllUsers(users []database.User, currUser string) {
	for _, user := range users {
		var name string
		if user.Name == currUser {
			name = fmt.Sprintf("%s (current)", user.Name)
		} else {
			name = user.Name
		}
		fmt.Printf("* %s\n", name)
	}
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time>", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %s...\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)

	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Printf("Scraping error: %v\n", err)
			continue
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error following feed: %v", err)
	}
	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf(" * ID:\t\t%v\n", feed.ID)
	fmt.Printf(" * Created:\t%v\n", feed.CreatedAt)
	fmt.Printf(" * Updated:\t%v\n", feed.UpdatedAt)
	fmt.Printf(" * Name:\t%v\n", feed.Name)
	fmt.Printf(" * URL:\t\t%v\n", feed.Url)
	fmt.Printf(" * USER_ID:\t%v\n", feed.UserID)
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("no arguments should be provided for this command, %s", cmd.Name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found in the database.")
		return nil
	}

	for _, feed := range feeds {
		fmt.Printf("* Name:        %s\n", feed.Name)
		fmt.Printf("* URL:         %s\n", feed.Url)
		fmt.Printf("* Created by:  %s\n", feed.UserName)
		fmt.Println("--------------------")
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("an error occurred: %v", err)
	}

	createdFeedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("an error occurred: %v", err)
	}

	fmt.Println("Feed Followed Successfully:")
	fmt.Printf("Name: %s\n", createdFeedFollow.FeedName)
	fmt.Printf("CurrentUserName: %s\n", user.Name)
	fmt.Println("==============================================")

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("no arguments should be provided for this command, %s", cmd.Name)
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("an error occurred: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Printf("%s is not following any feeds\n", user.Name)
		return nil
	}

	fmt.Printf("%s is following all the feeds listed below:\n", user.Name)
	for i, feed := range feeds {
		fmt.Printf("%d. %s\n", i+1, feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	feedUrl := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error unfollowing feed %s: %v", feedUrl, err)
	}

	fmt.Println("Feed successfully unfollowed!")
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) > 0 {
		if i, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = i
		} else {
			return fmt.Errorf("invalid limit: %s", cmd.Args[0])
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error getting posts from db: %v", err)
	}
	fmt.Printf("Showing %d posts for user %s:\n", len(posts), user.Name)

	for _, post := range posts {
		pubDate := "Unknown date"
		if post.PublishedAt.Valid {
			pubDate = post.PublishedAt.Time.Format("2006-01-02 15:04")
		}

		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("PubDate: %s\n", pubDate)
		fmt.Printf("Link:   %s\n", post.Url)
		if post.Description.Valid {
			fmt.Printf("Desc:   %s\n", post.Description.String)
		}
		fmt.Println()
	}

	return nil
}
