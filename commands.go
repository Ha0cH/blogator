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

func handlerAddFeed(s *state, cmd command) error {
	currentUser := s.cfg.CurrentUserName
	if currentUser == "" {
		return fmt.Errorf("No user is currently logged in. Please login first.")
	}

	if len(cmd.args) < 2 {
		return fmt.Errorf("AddFeed command expects two arguments: Feed Name and Feed URL")
	}

	user, err := s.db.GetUser(context.Background(), currentUser)
	if err != nil {
		return fmt.Errorf("Error retrieving user: %v", err)
	}

	userID := user.ID
	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    userID,
	}
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error creating feed: %v", err)
	}

	fmt.Printf("Feed added successfully: %+v\n", feed)

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("Error creating feed follow: %v", err)
	}

	fmt.Printf("User %s is now following feed %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerGetAllFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("feeds command does not expect any arguments")
	}

	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving feeds: %v", err)
	}

	fmt.Println("List of feeds:")
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Error retrieving user for feed: %v", err)
		}
		fmt.Printf("* %s (%s), by %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func handlerFollowFeed(s *state, cmd command) error {
	currentUser := s.cfg.CurrentUserName
	if currentUser == "" {
		return fmt.Errorf("No user is currently logged in. Please login first.")
	}

	if len(cmd.args) != 1 {
		return fmt.Errorf("Follow command expects one argument: Feed URL")
	}

	feedURL := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Error retrieving feed: %v", err)
	}

	user, err := s.db.GetUser(context.Background(), currentUser)
	if err != nil {
		return fmt.Errorf("Error retrieving user: %v", err)
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error creating feed follow: %v", err)
	}

	fmt.Printf("User %s is now following feed %s\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerGetFeedFollowsForUser(s *state, cmd command) error {
	currentUser := s.cfg.CurrentUserName
	if currentUser == "" {
		return fmt.Errorf("No user is currently logged in. Please login first.")
	}

	user, err := s.db.GetUser(context.Background(), currentUser)
	if err != nil {
		return fmt.Errorf("Error retrieving user: %v", err)
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error retrieving feed follows: %v", err)
	}

	fmt.Printf("Feeds followed by user %s:\n", currentUser)
	for _, feedFollow := range feedFollows {
		fmt.Printf("* %s\n", feedFollow.FeedName)
	}
	return nil
}
