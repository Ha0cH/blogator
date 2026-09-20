package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Ha0cH/blogator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmdLookUp map[string]func(*state, command) error
}

func (cmds *commands) run(s *state, cmd command) error {
	c, ok := cmds.cmdLookUp[cmd.name]
	if !ok {
		return fmt.Errorf("Command not found")
	}

	return c(s, cmd)
}

func (cmds *commands) register(name string, f func(*state, command) error) {
	cmds.cmdLookUp[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Login command expects a single argument: Username")
	}

	userName := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("The user has been set.")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Register command expects an argument: Username")
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	u, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(u.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been created with ID: %s\n", u.Name, u.ID)
	fmt.Printf("%+v\n", u)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("All users have been deleted.")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("List of users:")

	currentUser := s.cfg.CurrentUserName
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %s\n", user.Name)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	f, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("Fetched feed: %+v\n", f)
	return nil
}
