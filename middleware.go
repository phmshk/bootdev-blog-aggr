package main

import (
	"context"
	"fmt"

	"github.com/phmshk/bootdev-blog-aggr/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currUserName := s.cfg.CurrentUserName
		currUser, err := s.db.GetUser(context.Background(), currUserName)
		if err != nil {
			return fmt.Errorf("error fetching user: %v", err)
		}

		return handler(s, cmd, currUser)
	}
}
