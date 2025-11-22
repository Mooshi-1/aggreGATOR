package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"mooshi-1/aggregator/internal/database"
	"net/http"
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

func handlerAgg(s *state, cmd command) error {

	url := "https://www.wagslane.dev/index.xml"
	filled, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("agg failed: %w", err)
	}
	fmt.Print(filled)

	return nil
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch feed failure: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("resp failure: %w", err)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io failure: %w", err)
	}

	filledStruct := &RSSFeed{}

	err = xml.Unmarshal(content, filledStruct)
	if err != nil {
		return nil, fmt.Errorf("xml unmarshal error: %w", err)
	}

	filledStruct.Channel.Title = html.UnescapeString(filledStruct.Channel.Title)
	filledStruct.Channel.Description = html.UnescapeString(filledStruct.Channel.Description)
	for _, sli := range filledStruct.Channel.Item {
		sli.Title = html.UnescapeString(sli.Title)
		sli.Description = html.UnescapeString(sli.Description)
	}

	return filledStruct, nil
}

func handleAddFeed(s *state, cmd command) error {
	usr, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
	if err != nil {
		return fmt.Errorf("addFeed: getUser: failure|%w", err)
	}

	if len(cmd.Args) != 2 {
		return fmt.Errorf("addFeed: args incorrect")
	}

	fname := cmd.Args[0]
	url := cmd.Args[1]

	newFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   fname,
		Url:    url,
		UserID: usr.ID,
	})
	if err != nil {
		return fmt.Errorf("addFeed: ceateFeed: failure|%w", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:     uuid.New(),
			UserID: usr.ID,
			FeedID: newFeed.ID,
		})
	if err != nil {
		return fmt.Errorf("addFeed : create FF : falure %w", err)
	}

	fmt.Println(newFeed)
	return nil

}

func handleFeeds(s *state, cmd command) error {

	allFeeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("getallFeeds: falure|%w", err)
	}

	for _, f := range allFeeds {
		fmt.Println(f)
	}

	return nil
}

func handleFollow(s *state, cmd command) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("cmd requires 1 arg")
	}
	uurl := cmd.Args[0]

	usr, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
	if err != nil {
		return fmt.Errorf("follow : getusr err %w", err)
	}

	efeed, err := s.db.GetFeedByURL(context.Background(), uurl)
	if err != nil {
		return fmt.Errorf("follow : getfeed err %w", err)
	}

	ff, err := s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:     uuid.New(),
			UserID: usr.ID,
			FeedID: efeed.ID,
		})
	if err != nil {
		return fmt.Errorf("follow : create ff err %w", err)
	}

	fmt.Println(ff)

	return nil
}

func handleFollowing(s *state, cmd command) error {

	ff, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUser)
	if err != nil {
		return fmt.Errorf("following : get ff : %w", err)
	}

	for _, f := range ff {
		fmt.Println(f)
	}

	return nil

}
