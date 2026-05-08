package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/phmshk/bootdev-blog-aggr/internal/api"
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
	createdUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
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
	if len(cmd.Args) > 0 {
		return fmt.Errorf("no arguments should be provided for this command, %s", cmd.Name)
	}

	feed, err := api.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("an error occurred: %v", err)
	}

	fmt.Printf("%+v\n", feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	currUserName := s.cfg.CurrentUserName
	currUser, err := s.db.GetUser(context.Background(), currUserName)
	if err != nil {
		return fmt.Errorf("error fetching user: %v", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    currUser.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
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
