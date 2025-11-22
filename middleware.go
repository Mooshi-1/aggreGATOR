package main

import "mooshi-1/aggregator/internal/database"

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return nil
}
