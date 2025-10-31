package main

import (
	"fmt"
	"mooshi-1/aggregator/internal/config"
	"os"
)

type state struct {
	cfg *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	allCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {

	if _, ok := c.allCommands[cmd.Name]; !ok {
		return fmt.Errorf("%v command does not exist", cmd.Name)
	}

	h := c.allCommands[cmd.Name]
	err := h(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {

	c.allCommands[name] = f
}

func main() {
	cfgG := config.ReadConfig()

	currentState := &state{}
	currentState.cfg = cfgG

	commandsMap := commands{
		allCommands: make(map[string]func(*state, command) error),
	}
	commandsMap.register("login", handlerLogin)

	entry := os.Args
	if len(entry) < 2 {
		fmt.Printf("only 1 arg provided, %v\n", entry)
		os.Exit(1)
	}
	var title string
	var titleArgs []string

	for index, single := range entry {
		if index == 1 {
			title = single
		} else if index != 0 {
			titleArgs = append(titleArgs, single)
		}
	}

	perform := command{
		Name: title,
		Args: titleArgs,
	}

	err := commandsMap.run(currentState, perform)
	if err != nil {
		fmt.Print("error performing action\n")
		os.Exit(1)
	}

}

func handlerLogin(s *state, cmd command) error {
	if cmd.Args == nil {
		return fmt.Errorf("handler login requires username")
	}
	name := cmd.Args[0]
	s.cfg.SetUser(name)
	fmt.Printf("username set to %s\n", name)
	return nil
}
