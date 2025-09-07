package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	gofeed "github.com/mmcdole/gofeed"
)

var (
	tokenVariableName = "TELEGRAM_BOT_TOKEN"
)

func getRssFeedsLinks() ([]string, error) {
	rssFeeds := GetEnvOrDefault("RSS_FEEDS", "")

	if len(rssFeeds) == 0 {
		return RSSFeeds, nil
	}

	var linksList []string
	err := json.Unmarshal([]byte(rssFeeds), &linksList)

	if err != nil {
		return nil, err
	}

	fmt.Printf("os.Getenv(\"RSS_FEEDS\"): %v\n", os.Getenv("RSS_FEEDS"))
	return linksList, nil
}

func parseFeed(rssLink string, db *Database) (string, error) {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(rssLink)

	if err != nil {
		return "", err
	}

	resultMessage := ""

	for _, item := range feed.Items {
		hasLink, err := db.HasLink(item.Link)
		if err != nil {
			log.Printf("Error while checking if link is already in database: %v", err)
		}

		if hasLink {
			continue
		}

		resultMessage += fmt.Sprintln(item.Title) + fmt.Sprintln(item.Link) + "\n"
		db.AddLink(item.Link)
	}

	if len(resultMessage) > 0 {
		resultMessage = fmt.Sprintln(feed.Title) + "\n" + resultMessage
	}

	return resultMessage, nil
}

func initEnvAndCheckToken() (*string, error) {
	tokenValue, envVarExists := os.LookupEnv(tokenVariableName)
	if !envVarExists {
		if _, err := os.Stat(".env"); err == nil {
			err := godotenv.Load(".env")
			if err != nil {
				return nil, fmt.Errorf(".env file load error: %w", err)
			}

			tokenValue, envVarExists = os.LookupEnv(tokenVariableName)
			if !envVarExists {
				return nil, fmt.Errorf("token variable %s not found in .env file", tokenVariableName)
			}
		} else {
			return nil, fmt.Errorf(".env file not found and token variable %s not set", tokenVariableName)
		}
	}

	return &tokenValue, nil
}

func main() {
	logger := NewLogger()

	token, err := initEnvAndCheckToken()
	if err != nil {
		logger.Error("Error initialize environment variables", "error", err)
		os.Exit(1)
	}

	db, err := NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	logger.Debug("Environment variables successfully initialized")

	bot, err := NewTelebot(*token)

	if err != nil {
		logger.Error("error while creating bot", "error", err)
		os.Exit(1)
	}

	logger.Debug("Telebot successfully created")

	recipients, err := GetRecipients()

	if err != nil {
		logger.Error("Can not initialize messages recipients", "error", err)
		os.Exit(1)
	}

	// bot.Handle("/start", func(c telebot.Context) error {
	// 	fp := gofeed.NewParser()

	// 	feed, err := fp.ParseURL("https://cprss.s3.amazonaws.com/javascriptweekly.com.xml")

	// 	if err != nil {
	// 		log.Println(fmt.Errorf("can not request javascriptweekly.com: %w", err))
	// 	}

	// 	item := feed.Items[0]

	// 	description, err := htmltomarkdown.ConvertString(item.Description)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	fmt.Printf("c.Chat().ID: %v\n", c.Chat().ID)

	// 	return c.Send(description[0:4096])

	// 	return c.Send("bot started")
	// })

	rssFeeds, err := getRssFeedsLinks()

	if err != nil {
		logger.Error("Error while parsing RSS Feeds", "error", err)
		os.Exit(1)
	}

	logger.Debug("RSS Feeds successfully parsed", "rssFeeds", rssFeeds)

	logger.Debug("Start fetching RSS Feeds")
	for _, rssFeed := range rssFeeds {
		logger.Debug("Fetching RSS Feed", "rssFeed", rssFeed)
		message, err := parseFeed(rssFeed, db)

		if err != nil {
			logger.Warn("Failed to fetch RSS Feed", "rssFeed", rssFeed)
		} else {
			logger.Debug("RSS Feed successfully fetched and handled", "rssFeed", rssFeed)

			logger.Debug("Start sending messages to recipients")
			for _, recipient := range recipients {
				bot.Send(recipient, message)
			}
			logger.Debug("Finish sending messages to recipients")
		}

	}

	// bot.Start()
}
