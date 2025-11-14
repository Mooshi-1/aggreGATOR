package main

import (
	"context"
	"fmt"
	"mooshi-1/aggregator/internal/database"
	"os"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if cmd.Args == nil {
		return fmt.Errorf("handler login requires username")
	}
	name := cmd.Args[0]

	userFromDB, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		fmt.Printf("user does not found")
		os.Exit(1)
	}
	s.cfg.SetUser(userFromDB.Name)
	fmt.Printf("username set to %s\n", userFromDB.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if cmd.Args == nil {
		return fmt.Errorf("args are empty")
	}

	user := database.CreateUserParams{
		ID:   uuid.New(),
		Name: cmd.Args[0],
	}

	userFromDB, err := s.db.GetUser(context.Background(), user.Name)
	if err != nil {
		fmt.Printf("user does not exist yet, creating")
	}

	if userFromDB.Name == user.Name {
		fmt.Printf("user already exists")
		os.Exit(1)
	}

	newUser, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("db error: %v", err)
	}

	s.cfg.SetUser(newUser.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {

	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("issue resetting db: %v", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {

	usr, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("db error all users: %v", err)
	}

	for _, id := range usr {
		if s.cfg.CurrentUser == id.Name {
			fmt.Printf("%v (current)\n", id.Name)
		} else {
			fmt.Printf("%v\n", id.Name)
		}
	}

	return nil
}
