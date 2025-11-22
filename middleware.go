package main

import (
	"context"
	"fmt"
	"mooshi-1/aggregator/internal/database"
)

func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) func(*state, command) error {

	return func(s *state, cmd command) error {
		usr, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
		if err != nil {
			return fmt.Errorf("middleware: logged in failure|%w", err)
		}

		return handler(s, cmd, usr)
	}
}
